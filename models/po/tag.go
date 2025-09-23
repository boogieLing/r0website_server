package po

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// Tag 标签实体，维护标签与图片的倒排关系
type Tag struct {
	ID          string             `bson:"_id"`                    // 标签ID，使用标签名称的小写形式
	Name        string             `bson:"name"`                   // 标签名称（原始大小写）
	DisplayName string             `bson:"displayName"`            // 显示名称
	Description string             `bson:"description,omitempty"`  // 标签描述
	ImageCount  int                `bson:"imageCount"`             // 使用该标签的图片数量
	Color       string             `bson:"color,omitempty"`        // 标签颜色（可选）
	Category    string             `bson:"category,omitempty"`     // 标签分类（如：风格、主题、技术等）
	Images      []TagImageRef      `bson:"images,omitempty"`       // 使用该标签的图片列表（倒排索引）
	CreatedAt   time.Time          `bson:"createdAt"`              // 创建时间
	UpdatedAt   time.Time          `bson:"updatedAt"`              // 更新时间
}

// TagImageRef 标签中的图片引用（用于倒排索引）
type TagImageRef struct {
	ImageID   primitive.ObjectID `bson:"imageId"`   // 图片ID
	ImageName string             `bson:"imageName"` // 图片名称
	AddedAt   time.Time          `bson:"addedAt"`   // 添加时间
}

// TagStats 标签统计信息
type TagStats struct {
	TagID       string `bson:"tagId"`
	Name        string `bson:"name"`
	ImageCount  int    `bson:"imageCount"`
	DisplayName string `bson:"displayName"`
}