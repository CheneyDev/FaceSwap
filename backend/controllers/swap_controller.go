package controllers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"face-swap/config"
	"face-swap/models"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

type SwapInput struct {
	OriginalImage string `json:"original_image"`
	ImageToSwap   string `json:"image_to_swap"`
}

func SwapImage(c *gin.Context) {
	var input SwapInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 上传原始图片到S3
	originalImageURL, err := uploadToS3(input.OriginalImage)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// 上传待替换图片到S3
	imageToSwapURL, err := uploadToS3(input.ImageToSwap)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// 调用Replicate API
	replicateInput := ReplicateRequest{
		Version: "8d0b076a2aff3904dfcec3253c778e0310a68f78483c4699c7fd800f3051d2b3",
		Input: ReplicateInput{
			Image:         originalImageURL,
			ImageToBecome: imageToSwapURL,
		},
	}
	replicateToken := "Bearer " + os.Getenv("REPLICATE_API_TOKEN")
	reqBody, _ := json.Marshal(replicateInput)
	req, _ := http.NewRequest("POST", "https://api.replicate.com/v1/predictions", bytes.NewBuffer(reqBody))
	req.Header.Set("Authorization", replicateToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
		}
	}(resp.Body)

	var prediction ReplicateResponse
	if err := json.NewDecoder(resp.Body).Decode(&prediction); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// 轮询预测结果
	for {
		req, _ := http.NewRequest("GET", prediction.URLs.Get, nil)
		req.Header.Set("Authorization", replicateToken)

		resp, err := client.Do(req)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
			}
		}(resp.Body)

		body, _ := ioutil.ReadAll(resp.Body)
		if err := json.Unmarshal(body, &prediction); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		if prediction.Status == "succeeded" {
			break
		} else if prediction.Status == "failed" {
			c.JSON(500, gin.H{"error": prediction.Error})
			return
		}

		time.Sleep(1 * time.Second)
	}

	resultImageURL := prediction.Output[0]
	// 上传结果图片到S3
	resultImageURL, err = uploadS3ByURL(resultImageURL)

	// 创建 ImageSwapRecord
	record := models.ImageSwapRecord{
		UserID:        1, // 这里简化了,假设用户ID为1
		OriginalImage: originalImageURL,
		ImageToSwap:   imageToSwapURL,
		ResultImage:   resultImageURL,
	}
	config.DB.Create(&record)

	c.JSON(200, record)
}

func uploadToS3(base64Data string) (string, error) {
	// 解析Base64格式并提取MIME类型和纯净的Base64数据
	splitData := strings.SplitN(base64Data, ",", 2)
	if len(splitData) != 2 {
		return "", errors.New("invalid base64 data")
	}
	mimeType := strings.TrimPrefix(splitData[0], "data:")
	mimeType = strings.SplitN(mimeType, ";", 2)[0]
	cleanBase64Data := splitData[1]

	// 确定文件扩展名
	var fileExt string
	switch mimeType {
	case "image/jpeg":
		fileExt = "jpeg"
	case "image/jpg":
		fileExt = "jpg"
	case "image/png":
		fileExt = "png"
	case "image/webp":
		fileExt = "webp"
	// 添加更多的MIME类型和文件扩展名匹配
	default:
		return "", fmt.Errorf("unsupported MIME type: %s", mimeType)
	}

	// 解码Base64图片数据
	dec, err := base64.StdEncoding.DecodeString(cleanBase64Data)
	if err != nil {
		return "", err
	}

	// 生成唯一的图片名称
	imageName := fmt.Sprintf("%d.%s", time.Now().UnixNano(), fileExt)

	// 获取S3客户端
	s3Client := config.NewS3Client()

	// 上传图片到S3
	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(os.Getenv("R2_BUCKET_NAME")),
		Key:           aws.String(imageName),
		Body:          bytes.NewReader(dec),
		ContentType:   aws.String(mimeType),
		ContentLength: aws.Int64(int64(len(dec))),
	})
	if err != nil {
		return "", err
	}

	// 生成并返回图片的公共URL
	imageURL := fmt.Sprintf("%s/%s", os.Getenv("R2_PUBLIC_URL"), imageName)
	return imageURL, nil
}

func uploadS3ByURL(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Error closing response body:", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	imageName := fmt.Sprintf("%d.jpg", time.Now().UnixNano())
	s3Client := config.NewS3Client()

	var contentLength int64
	if resp.ContentLength > 0 {
		contentLength = resp.ContentLength
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	reader := bytes.NewReader(body)
	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(os.Getenv("R2_BUCKET_NAME")),
		Key:           aws.String(imageName),
		Body:          reader,
		ContentType:   aws.String("image/jpeg"),
		ContentLength: aws.Int64(contentLength),
	})

	imageURL := fmt.Sprintf("%s/%s", os.Getenv("R2_PUBLIC_URL"), imageName)
	return imageURL, nil
}
