package config

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"os"
)

func NewS3Client() *s3.S3 {
	// 从环境变量中读取R2的访问密钥和秘密访问密钥
	accessKey := os.Getenv("R2_ACCESS_KEY")
	secretKey := os.Getenv("R2_SECRET_KEY") // 请确保设置了这个环境变量
	accountID := os.Getenv("R2_ACCOUNT_ID")

	// 检查密钥是否为空
	if accessKey == "" || secretKey == "" {
		fmt.Println("R2_ACCESS_KEY or R2_SECRET_KEY is not set in environment variables")
		return nil
	}

	// 配置S3客户端
	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(accessKey, secretKey, ""),
		Endpoint:         aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountID)),
		Region:           aws.String("auto"),
		S3ForcePathStyle: aws.Bool(true),
	}
	newSession := session.New(s3Config)
	s3Client := s3.New(newSession)

	return s3Client
}
