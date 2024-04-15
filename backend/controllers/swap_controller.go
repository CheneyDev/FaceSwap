package controllers

import (
	"bytes"
	"encoding/base64"
	"errors"
	"face-swap/config"
	"face-swap/models"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
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

	// TODO: 调用图片替换服务,获取结果图片
	resultImageURL := "https://example.com/result6666666.jpg"

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
	imageURL := fmt.Sprintf("https://%s.r2.cloudflarestorage.com/%s", os.Getenv("R2_ACCOUNT_ID"), imageName)
	return imageURL, nil
}
