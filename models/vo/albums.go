package vo

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"r0Website-server/models/po"
)

// AlbumDetailVo 图集详情返回结构
type AlbumDetailVo struct {
	ID          primitive.ObjectID  `json:"id"`
	Title       string              `json:"title"`
	Description string              `json:"description"`
	CoverImage  primitive.ObjectID  `json:"cover_image"`
	CreatedAt   string              `json:"created_at"`
	UpdatedAt   string              `json:"updated_at"`
	Tags        []string            `json:"tags"`
	Visibility  string              `json:"visibility"`
	ImageRefs   []*po.AlbumImageRef `json:"image_refs"` // 原样返回布局与描述
}

// AlbumListVo 图集列表返回结构
type AlbumListVo struct {
	Total  int64            `json:"total"`
	Albums []*AlbumDetailVo `json:"albums"`
}
