// Package po
/**
 * @Author: r0
 * @Mail: boogieLing_o@qq.com
 * @Description: 单篇文章的模型
 * @File:  article
 * @Version: 1.0.0
 * @Date: 2022/7/4 18:50
 */
package po

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Article struct {
	Id primitive.ObjectID `bson:"_id,omitempty"` // Mongo 主键 _id
	// Uuid     string             `bson:"uuid"`          // 本系统需要的额外主键
	Title    string `bson:"title"`    // 文章标题
	Author   string `bson:"author"`   // 作者
	Synopsis string `bson:"synopsis"` // 备注
	PicUrl   string `json:"pic_url"`  // 图片的链接
	// Detail         string             `bson:"detail"`          // htm内容
	Markdown       string    `bson:"markdown"`        // md内容
	MdWords        string    `bson:"md_words"`        // md分词内容
	TitleWords     string    `bson:"title_words"`     // 标题分词内容
	DeleteFlag     bool      `bson:"delete_flag"`     // 是否已删除
	DraftFlag      bool      `bson:"draft_flag"`      // 是否为草稿
	Overhead       bool      `bson:"overhead"`        // 是否置顶
	ArtLength      int64     `bson:"art_length"`      // 文章长度
	ReadsNumber    int64     `bson:"reads_number"`    // 阅读数
	CommentsNumber int64     `bson:"comments_number"` // 评论数
	PraiseNumber   int64     `bson:"praise_number"`   // 点赞数
	Tags           []string  `bson:"tags"`            // 标签
	Categories     []string  `bson:"categories"`      // 分类
	CreateTime     time.Time `bson:"create_time"`     // 创建时间
	UpdateTime     time.Time `bson:"update_time"`     // 更新时间
}
