package utils

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/tencentyun/cos-go-sdk-v5"
	"r0Website-server/config"
)

type COSClient struct {
	client *cos.Client
	bucket string
}

// NewCOSClient 创建COS客户端
func NewCOSClient(cfg *config.SystemConfig) (*COSClient, error) {
	// 腾讯云COS配置
	secretID := cfg.TencentCloud.SecretID
	secretKey := cfg.TencentCloud.SecretKey
	region := cfg.TencentCloud.Region
	bucket := cfg.TencentCloud.Bucket

	// 创建COS客户端
	u, _ := cos.NewBucketURL(bucket, region, true)
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  secretID,
			SecretKey: secretKey,
		},
	})

	return &COSClient{
		client: client,
		bucket: bucket,
	}, nil
}

// UploadFile 上传文件到COS
func (c *COSClient) UploadFile(file multipart.File, fileHeader *multipart.FileHeader, objectKey string) (string, error) {
	// 获取文件大小
	fileSize := fileHeader.Size
	if fileSize == 0 {
		return "", fmt.Errorf("文件大小不能为0")
	}

	// 读取文件内容
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		log.Printf("读取文件失败: %v", err)
		return "", err
	}

	// 上传文件
	_, err = c.client.Object.Put(context.Background(), objectKey, strings.NewReader(string(fileBytes)), nil)
	if err != nil {
		log.Printf("上传文件到COS失败: %v", err)
		return "", err
	}

	// 构建访问URL：直接复用 SDK 生成的 BucketURL，避免重复拼接 bucket 和域名
	base := strings.TrimRight(c.client.BaseURL.BucketURL.String(), "/")
	fileURL := fmt.Sprintf("%s/%s", base, objectKey)
	log.Printf("文件上传成功: %s", fileURL)
	return fileURL, nil
}

// UploadThumbnail 上传缩略图到COS
func (c *COSClient) UploadThumbnail(thumbnailBytes []byte, originalFilename string, objectKey string) (string, error) {
	if len(thumbnailBytes) == 0 {
		return "", fmt.Errorf("缩略图数据不能为空")
	}

	// 上传缩略图
	_, err := c.client.Object.Put(context.Background(), objectKey, bytes.NewReader(thumbnailBytes), nil)
	if err != nil {
		log.Printf("上传缩略图到COS失败: %v", err)
		return "", err
	}

	// 构建访问URL：与原图同样逻辑
	base := strings.TrimRight(c.client.BaseURL.BucketURL.String(), "/")
	thumbnailURL := fmt.Sprintf("%s/%s", base, objectKey)
	log.Printf("缩略图上传成功: %s", thumbnailURL)
	return thumbnailURL, nil
}

// GenerateObjectKey 生成COS对象存储路径
func (c *COSClient) GenerateObjectKey(originalFilename string) string {
	// 获取文件扩展名
	ext := strings.ToLower(filepath.Ext(originalFilename))
	if ext == "" {
		ext = ".jpg" // 默认扩展名
	}

	// 生成唯一文件名：时间戳 + 随机字符串
	timestamp := time.Now().UnixNano()
	uniqueName := fmt.Sprintf("somnium/primitive/%d%s", timestamp, ext)

	return uniqueName
}

// GenerateThumbnailObjectKey 生成缩略图COS对象存储路径
func (c *COSClient) GenerateThumbnailObjectKey(originalFilename string) string {
	// 获取文件扩展名
	ext := strings.ToLower(filepath.Ext(originalFilename))
	if ext == "" {
		ext = ".jpg" // 默认扩展名
	}

	// 生成唯一文件名：时间戳 + 随机字符串
	timestamp := time.Now().UnixNano()
	uniqueName := fmt.Sprintf("somnium/compressed/%d_thumb%s", timestamp, ext)

	return uniqueName
}

// DeleteFile 从COS删除文件
func (c *COSClient) DeleteFile(objectKey string) error {
	_, err := c.client.Object.Delete(context.Background(), objectKey)
	if err != nil {
		log.Printf("从COS删除文件失败: %v", err)
		return err
	}
	log.Printf("文件删除成功: %s", objectKey)
	return nil
}

// GetFileInfo 获取文件信息
func (c *COSClient) GetFileInfo(objectKey string) (*cos.Response, error) {
	resp, err := c.client.Object.Get(context.Background(), objectKey, nil)
	if err != nil {
		log.Printf("获取文件信息失败: %v", err)
		return nil, err
	}
	return resp, nil
}
