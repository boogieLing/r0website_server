package utils

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"strings"

	"github.com/disintegration/imaging"
)

// ThumbnailConfig 缩略图配置
type ThumbnailConfig struct {
	MaxWidth  int
	MaxHeight int
	Quality   int
}

// DefaultThumbnailConfig 默认缩略图配置
var DefaultThumbnailConfig = ThumbnailConfig{
	MaxWidth:  300,
	MaxHeight: 300,
	Quality:   85,
}

// GenerateThumbnail 生成缩略图
func GenerateThumbnail(file multipart.File, header *multipart.FileHeader, config ThumbnailConfig) ([]byte, int, int, error) {
	// 读取文件内容
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("读取文件失败: %v", err)
	}

	// 解码图片
	img, _, err := image.Decode(bytes.NewReader(fileBytes))
	if err != nil {
		return nil, 0, 0, fmt.Errorf("解码图片失败: %v", err)
	}

	// 获取原始尺寸
	originalWidth := img.Bounds().Dx()
	originalHeight := img.Bounds().Dy()

	// 计算缩略图尺寸（保持宽高比）
	thumbnailWidth, thumbnailHeight := calculateThumbnailSize(originalWidth, originalHeight, config.MaxWidth, config.MaxHeight)

	// 生成缩略图
	thumbnail := imaging.Resize(img, thumbnailWidth, thumbnailHeight, imaging.Lanczos)

	// 编码缩略图
	var buf bytes.Buffer
	contentType := header.Header.Get("Content-Type")

	switch strings.ToLower(contentType) {
	case "image/jpeg", "image/jpg":
		err = jpeg.Encode(&buf, thumbnail, &jpeg.Options{Quality: config.Quality})
	case "image/png":
		err = png.Encode(&buf, thumbnail)
	default:
		// 默认使用JPEG格式
		err = jpeg.Encode(&buf, thumbnail, &jpeg.Options{Quality: config.Quality})
	}

	if err != nil {
		return nil, 0, 0, fmt.Errorf("编码缩略图失败: %v", err)
	}

	return buf.Bytes(), thumbnailWidth, thumbnailHeight, nil
}

// calculateThumbnailSize 计算缩略图尺寸（保持宽高比）
func calculateThumbnailSize(originalWidth, originalHeight, maxWidth, maxHeight int) (int, int) {
	// 如果原图已经小于最大尺寸，返回原图尺寸
	if originalWidth <= maxWidth && originalHeight <= maxHeight {
		return originalWidth, originalHeight
	}

	// 计算缩放比例
	widthRatio := float64(maxWidth) / float64(originalWidth)
	heightRatio := float64(maxHeight) / float64(originalHeight)

	// 使用较小的比例，确保图片能完全显示在缩略图尺寸内
	scale := widthRatio
	if heightRatio < scale {
		scale = heightRatio
	}

	// 计算缩略图尺寸
	thumbnailWidth := int(float64(originalWidth) * scale)
	thumbnailHeight := int(float64(originalHeight) * scale)

	return thumbnailWidth, thumbnailHeight
}