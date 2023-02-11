package po

import "go.mongodb.org/mongo-driver/bson/primitive"

// Category 分类的实体模型
type Category struct {
	Id         primitive.ObjectID `bson:"_id,omitempty"` // Mongo 主键 _id
	Name       string             `json:"name" bson:"name"`
	Count      int64              `json:"count" bson:"count"`
	ArticleIds []string           `json:"article_ids" bson:"article_ids"`
}
