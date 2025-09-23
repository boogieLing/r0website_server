// Package vo
/**
 * @Author: r0
 * @Mail: boogieLing_o@qq.com
 * @Description: 图片相关VO结构体
 * @File:  image_vo
 * @Version: 1.0.0
 * @Date: 2024/09/24
 */
package vo

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"r0Website-server/models/po"
)

// UploadImageVo 上传图片请求参数
type UploadImageVo struct {
	Name string   `form:"name"`        // 图片名称，可选
	Tags []string `form:"tags"`        // 标签数组，可选
}

// ImageDetailVo 图片详情返回数据
type ImageDetailVo struct {
	ID         primitive.ObjectID            `json:"id"`
	Name       string                        `json:"name"`
	CosURL     string                        `json:"cos_url"`
	Width      int                           `json:"width"`
	Height     int                           `json:"height"`
	Size       int64                         `json:"size"`
	Format     string                        `json:"format"`
	Tags       []string                      `json:"tags"`
	Positions  map[string]po.CategoryPosition `json:"positions"`
	UploadedAt string                        `json:"uploaded_at"`
	Exif       map[string]string             `json:"exif"`
}

// ImageListVo 图片列表返回数据
type ImageListVo struct {
	Images []ImageDetailVo `json:"images"`
	Total  int64           `json:"total"`
	Page   int             `json:"page"`
	Size   int             `json:"size"`
}

// UpdateImagePositionVo 更新图片位置请求参数
type UpdateImagePositionVo struct {
	CategoryID string              `json:"categoryId" binding:"required"`
	Position   po.CategoryPosition `json:"position" binding:"required"`
}

// ImageSearchVo 图片搜索参数
type ImageSearchVo struct {
	Keyword  string `form:"keyword"`
	Page     int    `form:"page"`
	PageSize int    `form:"pageSize"`
}

// ImagesByCategoryVo 分类下图片查询参数
type ImagesByCategoryVo struct {
	CategoryID string `uri:"categoryId" binding:"required"`
	Page       int    `form:"page"`
	PageSize   int    `form:"pageSize"`
	Sort       string `form:"sort"`
	Order      string `form:"order"`
}