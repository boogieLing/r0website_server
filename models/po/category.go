package po

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// ImageCategory 图片分类集合，维护图片与分类的倒排关系
// 使用字符串ID作为主键，支持如 "nexus", "stillness" 等语义化ID
type ImageCategory struct {
	ID          string             `bson:"_id"`                    // 分类ID，如 "nexus", "stillness" 等
	Name        string             `bson:"name"`                   // 分类名称
	Description string             `bson:"description,omitempty"`  // 分类描述
	ImageCount  int                `bson:"imageCount"`             // 分类中的图片数量
	CoverImage  primitive.ObjectID `bson:"coverImage,omitempty"`   // 封面图片ID
	Settings    CategorySettings   `bson:"settings"`               // 分类设置
	Images      []CategoryImageRef `bson:"images,omitempty"`       // 分类中的图片引用列表（倒排索引）
	CreatedAt   time.Time          `bson:"createdAt"`              // 创建时间
	UpdatedAt   time.Time          `bson:"updatedAt"`              // 更新时间
}

// CategorySettings 分类设置
type CategorySettings struct {
	LayoutMode    string  `bson:"layoutMode"`    // 布局模式: freeform, flex, grid
	GridSize      int     `bson:"gridSize"`      // 网格大小
	DefaultWidth  float64 `bson:"defaultWidth"`  // 默认宽度
	DefaultHeight float64 `bson:"defaultHeight"` // 默认高度
	AutoArrange   bool    `bson:"autoArrange"`   // 是否自动排列
}

// CategoryImageRef 分类中的图片引用（用于倒排索引）
type CategoryImageRef struct {
	ImageID   primitive.ObjectID `bson:"imageId"`   // 图片ID
	SortOrder int                `bson:"sortOrder"` // 在分类中的排序序号
	AddedAt   time.Time          `bson:"addedAt"`   // 添加到分类的时间
}

// CategoryWithImages 分类及其图片列表（聚合查询结果）
type CategoryWithImages struct {
	Category  ImageCategory      `bson:"category"`
	Images    []CategoryImageRef `bson:"images"`
	Total     int64              `bson:"total"`
}