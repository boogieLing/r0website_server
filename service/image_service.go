// Package service
/**
 * @Author: r0
 * @Mail: boogieLing_o@qq.com
 * @Description: 图片相关服务
 * @File:  image_service
 * @Version: 1.0.0
 * @Date: 2024/09/24
 */
package service

import (
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"mime/multipart"
	"r0Website-server/dao"
	"r0Website-server/global"
	"r0Website-server/models/po"
	"r0Website-server/models/vo"
	"r0Website-server/utils"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ImageService struct {
	ImageDao         *dao.ImageDao         `R0Ioc:"true"`
	ImageCategoryDao *dao.ImageCategoryDao `R0Ioc:"true"`
	TagDao           *dao.TagDao           `R0Ioc:"true"`
	COSClient        *utils.COSClient      `R0Ioc:"true"`
}

const (
	MaxFileSize = 100 * 1024 * 1024 // 100MB
)

var allowedImageTypes = []string{"image/jpeg", "image/jpg", "image/png", "image/gif", "image/webp"}

// UploadImage 上传图片
func (s *ImageService) UploadImage(file multipart.File, header *multipart.FileHeader, params vo.UploadImageVo) (*vo.ImageDetailVo, error) {
	// 文件大小验证
	if header.Size > MaxFileSize {
		return nil, errors.New("文件大小不能超过100MB")
	}

	// 文件类型验证
	contentType := header.Header.Get("Content-Type")
	if !isValidImageType(contentType) {
		return nil, errors.New("不支持的文件类型，只允许: JPG, PNG, GIF, WebP")
	}

	// 获取图片信息（宽度、高度）
	imgConfig, format, err := image.DecodeConfig(file)
	if err != nil {
		return nil, errors.New("无法解析图片文件")
	}

	// 重置文件读取位置
	file.Seek(0, 0)

	// 生成缩略图
	thumbnailBytes, thumbWidth, thumbHeight, err := utils.GenerateThumbnail(file, header, utils.DefaultThumbnailConfig)
	if err != nil {
		global.Logger.Errorf("生成缩略图失败: %v", err)
		// 如果缩略图生成失败，继续上传原图，但不影响主流程
		thumbnailBytes = nil
		thumbWidth = 0
		thumbHeight = 0
	}

	// 重置文件读取位置（用于原图上传）
	file.Seek(0, 0)

	// 生成COS对象键
	objectKey := s.COSClient.GenerateObjectKey(header.Filename)

	// 上传到腾讯云COS
	cosURL, err := s.COSClient.UploadFile(file, header, objectKey)
	if err != nil {
		return nil, errors.New("上传文件失败")
	}

	// 上传缩略图（如果生成成功）
	var thumbURL string
	if thumbnailBytes != nil {
		thumbObjectKey := s.COSClient.GenerateThumbnailObjectKey(header.Filename)
		thumbURL, err = s.COSClient.UploadThumbnail(thumbnailBytes, header.Filename, thumbObjectKey)
		if err != nil {
			global.Logger.Errorf("上传缩略图失败: %v", err)
			// 如果缩略图上传失败，继续主流程，但不影响原图上传
			thumbURL = ""
		}
	}

	// 创建图片名称
	imageName := params.Name
	if imageName == "" {
		imageName = header.Filename
	}

	// 创建图片记录
	image := &po.Image{
		Name:        imageName,
		CosURL:      cosURL,
		ThumbURL:    thumbURL,
		Width:       imgConfig.Width,
		Height:      imgConfig.Height,
		ThumbWidth:  thumbWidth,
		ThumbHeight: thumbHeight,
		UploadedAt:  time.Now(),
		Tags:        params.Tags,
		EXIF:        make(map[string]string),
		Positions:   make(map[string]po.CategoryPosition),
	}

	// 添加到nexus分类（默认分类）
	nexusPosition := po.CategoryPosition{
		X:            0,
		Y:            0,
		Width:        float64(imgConfig.Width),
		Height:       float64(imgConfig.Height),
		GridX:        0,
		GridY:        0,
		GridSize:     10,
		Row:          0,
		Col:          0,
		LayoutMode:   "freeform",
		ZIndex:       0,
		IsVisible:    true,
		Version:      1,
		UpdatedAt:    time.Now(),
		CategoryName: "Nexus",
		SortOrder:    0,
		AddedAt:      time.Now(),
	}
	image.Positions["nexus"] = nexusPosition

	// 保存到数据库
	res, err := s.ImageDao.UploadImage(image)
	if err != nil {
		// 如果数据库保存失败，删除COS文件
		if deleteErr := s.COSClient.DeleteFile(objectKey); deleteErr != nil {
			global.Logger.Errorf("删除COS文件失败: %v", deleteErr)
		}
		return nil, errors.New("保存图片信息失败")
	}

	// 获取插入的图片ID
	imageID := res.InsertedID.(primitive.ObjectID)

	// 维护倒排索引：将图片添加到nexus分类
	if err := s.ImageCategoryDao.AddImageToCategory("nexus", imageID, 0); err != nil {
		global.Logger.Errorf("添加图片到nexus分类失败: %v", err)
		// 不中断主流程，只记录错误
	}

	// 维护标签倒排索引：同步更新标签
	if len(params.Tags) > 0 {
		for _, tagName := range params.Tags {
			if tag, err := s.TagDao.GetOrCreateTag(tagName, tagName, ""); err == nil {
				if err := s.TagDao.AddImageToTag(tag.ID, imageID, imageName); err != nil {
					global.Logger.Errorf("添加图片到标签 %s 失败: %v", tagName, err)
				}
			} else {
				global.Logger.Errorf("获取或创建标签 %s 失败: %v", tagName, err)
			}
		}
	}

	// 返回VO数据
	return &vo.ImageDetailVo{
		ID:          imageID,
		Name:        imageName,
		CosURL:      cosURL,
		ThumbURL:    thumbURL,
		Width:       imgConfig.Width,
		Height:      imgConfig.Height,
		ThumbWidth:  thumbWidth,
		ThumbHeight: thumbHeight,
		Size:        header.Size,
		Format:      format,
		Tags:        params.Tags,
		Positions:   image.Positions,
		UploadedAt:  image.UploadedAt.Format(time.RFC3339),
		Exif:        image.EXIF,
	}, nil
}

// GetImageDetail 获取图片详情
func (s *ImageService) GetImageDetail(imageID primitive.ObjectID) (*po.Image, error) {
	return s.ImageDao.GetImageByID(imageID)
}

// ListImages 获取所有图片列表
func (s *ImageService) ListImages() ([]*po.Image, error) {
	return s.ImageDao.ListImages()
}

// FindImagesByTag 通过标签查询图片
func (s *ImageService) FindImagesByTag(tag string) ([]*po.Image, error) {
	return s.ImageDao.FindImagesByTag(tag)
}

// SearchImagesByName 通过关键词模糊查找图片
func (s *ImageService) SearchImagesByName(keyword string) ([]*po.Image, error) {
	return s.ImageDao.SearchImagesByName(keyword)
}

// GetImageAlbums 获取当前图片被引用的所有图集
func (s *ImageService) GetImageAlbums(imageID primitive.ObjectID) ([]primitive.ObjectID, error) {
	return s.ImageDao.GetAllAlbumsOfImage(imageID)
}

// DeleteImage 删除图片记录
func (s *ImageService) DeleteImage(imageID primitive.ObjectID) error {
	return s.ImageDao.DeleteImageByID(imageID)
}

// UpdateImagePosition 更新图片在分类中的位置
func (s *ImageService) UpdateImagePosition(imageID primitive.ObjectID, categoryID string, position *po.CategoryPosition) error {
	// 设置更新时间
	position.UpdatedAt = time.Now()
	return s.ImageDao.UpdateImageCategoryPosition(imageID, categoryID, position)
}

// RemoveImageFromCategory 从分类中移除图片
func (s *ImageService) RemoveImageFromCategory(imageID primitive.ObjectID, categoryID string) error {
	return s.ImageDao.RemoveImageFromCategory(imageID, categoryID)
}

// GetImagesByCategory 获取分类下的图片（使用倒排索引优化）
func (s *ImageService) GetImagesByCategory(categoryID string, page, pageSize int, sort, order string) (*vo.ImageListVo, error) {
	// 使用倒排索引获取分类中的图片引用
	imageRefs, total, err := s.ImageCategoryDao.GetCategoryImages(categoryID, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("获取分类图片失败: %v", err)
	}

	// 获取图片详细信息
	var images []vo.ImageDetailVo
	for _, ref := range imageRefs {
		img, err := s.ImageDao.GetImageByID(ref.ImageID)
		if err == nil {
			images = append(images, vo.ImageDetailVo{
				ID:          img.ID,
				Name:        img.Name,
				CosURL:      img.CosURL,
				ThumbURL:    img.ThumbURL,
				Width:       img.Width,
				Height:      img.Height,
				ThumbWidth:  img.ThumbWidth,
				ThumbHeight: img.ThumbHeight,
				Size:        0,  // Size信息在Image结构体中不存在
				Format:      "", // Format信息在Image结构体中不存在
				Tags:        img.Tags,
				Positions:   img.Positions,
				UploadedAt:  img.UploadedAt.Format(time.RFC3339),
				Exif:        img.EXIF,
			})
		}
	}

	return &vo.ImageListVo{
		Images: images,
		Total:  total,
		Page:   page,
		Size:   pageSize,
	}, nil
}

// isValidImageType 检查是否为有效的图片类型
func isValidImageType(contentType string) bool {
	for _, t := range allowedImageTypes {
		if contentType == t {
			return true
		}
	}
	return false
}

// parseTags 解析标签字符串
func parseTags(tagsStr string) []string {
	if tagsStr == "" {
		return []string{}
	}
	return strings.Split(tagsStr, ",")
}
