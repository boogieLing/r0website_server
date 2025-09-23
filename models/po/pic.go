package po

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// Image 表示一张原始图片的元信息，存储于 images 集合中
type Image struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`  // MongoDB 主键
	Name       string             `bson:"name"`           // 图片名称或标题
	CosURL     string             `bson:"cos_url"`        // 腾讯云 COS 图片地址
	Width      int                `bson:"width"`          // 原始图片宽度（像素）
	Height     int                `bson:"height"`         // 原始图片高度（像素）
	UploadedAt time.Time          `bson:"uploaded_at"`    // 上传时间
	Tags       []string           `bson:"tags,omitempty"` // 可选：标签列表
	EXIF       map[string]string  `bson:"exif,omitempty"` // 可选：EXIF 数据（相机型号、光圈等）

	// 分类和位置信息 - 支持一个图片在多个分类中有不同的位置
	Positions map[string]CategoryPosition `bson:"positions" json:"positions"` // key: categoryID, value: 分类中的位置信息
}

// AlbumPosition 表示一张图在图集页面上的显示布局（坐标、尺寸、样式）
type AlbumPosition struct {
	X                 float64 `bson:"x"`                   // 横坐标（相对于画布或容器，建议归一化，如0.1表示10%）
	Y                 float64 `bson:"y"`                   // 纵坐标（同上）
	Width             float64 `bson:"width"`               // 显示宽度
	Height            float64 `bson:"height"`              // 显示高度
	Unit              string  `bson:"unit"`                // 尺寸单位："px"（绝对像素）、"%"（相对父容器）等
	Rotate            int     `bson:"rotate"`              // 旋转角度（单位：度）
	ZIndex            int     `bson:"z_index"`             // 层级顺序（前后遮挡）
	AspectRatioLocked bool    `bson:"aspect_ratio_locked"` // 是否保持原图宽高比
	Opacity           float64 `bson:"opacity"`             // 可选：透明度，0.0-1.0
	BorderRadius      float64 `bson:"border_radius"`       // 可选：圆角半径（px或相对单位）
	Shadow            string  `bson:"shadow,omitempty"`    // 可选：阴影参数，如 CSS box-shadow 格式
}

// AlbumImageRef 表示图集中一张图的引用及其布局信息
type AlbumImageRef struct {
	ImageID     primitive.ObjectID `bson:"image_id"`    // 引用 images 表中的图片ID
	Position    *AlbumPosition     `bson:"position"`    // 图在页面中的位置与尺寸布局
	Caption     string             `bson:"caption"`     // 可选：该图在图集中的标题或说明
	Description string             `bson:"description"` // 每张图的详细描述
}

// Album 表示一个图集，包含多个图片引用及其布局
type Album struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`    // MongoDB 主键
	Title       string             `bson:"title"`            // 图集标题
	Description string             `bson:"description"`      // 图集说明或备注
	CoverImage  primitive.ObjectID `bson:"cover_image"`      // 图集封面图（引用某个 image_id）
	CreatedAt   time.Time          `bson:"created_at"`       // 图集创建时间
	UpdatedAt   time.Time          `bson:"updated_at"`       // 图集最后修改时间
	ImageRefs   []*AlbumImageRef   `bson:"image_refs"`       // 图集中所有图片的引用与布局信息
	Tags        []string           `bson:"tags,omitempty"`   // 可选：图集标签
	Author      string             `bson:"author,omitempty"` // 可选：图集创建者
	Visibility  string             `bson:"visibility"`       // 可见性："public" | "private" | "unlisted"
}

// CategoryPosition 分类中的位置信息（freeform模式核心）
type CategoryPosition struct {
	// 核心位置信息
	X      float64 `bson:"x" json:"x"`          // X坐标（像素）
	Y      float64 `bson:"y" json:"y"`          // Y坐标（像素）
	Width  float64 `bson:"width" json:"width"`  // 宽度（像素）
	Height float64 `bson:"height" json:"height"` // 高度（像素）

	// 网格对齐信息
	GridX    int `bson:"gridX" json:"gridX"`     // 网格X位置
	GridY    int `bson:"gridY" json:"gridY"`     // 网格Y位置
	GridSize int `bson:"gridSize" json:"gridSize"` // 网格大小（默认10px）

	// 排序信息
	Row int `bson:"row" json:"row"` // 行号（用于排序）
	Col int `bson:"col" json:"col"` // 列号（用于排序）

	// 布局模式
	LayoutMode string `bson:"layoutMode" json:"layoutMode"` // 布局模式：freeform/flex

	// 层级和可见性
	ZIndex    int  `bson:"zIndex" json:"zIndex"`    // 层级顺序
	IsVisible bool `bson:"isVisible" json:"isVisible"` // 是否可见

	// 版本控制
	Version   int       `bson:"version" json:"version"`   // 数据版本
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"` // 位置更新时间

	// 分类信息
	CategoryName string    `bson:"categoryName" json:"categoryName"` // 分类名称
	SortOrder    int       `bson:"sortOrder" json:"sortOrder"`       // 在分类中的排序序号
	AddedAt      time.Time `bson:"addedAt" json:"addedAt"`          // 添加到分类的时间
}
