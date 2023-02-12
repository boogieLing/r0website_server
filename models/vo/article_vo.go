// Package vo
/**
 * @Author: r0
 * @Mail: boogieLing_o@qq.com
 * @Description: 文章视图模型
 * @File:  article
 * @Version: 1.0.0
 * @Date: 2022/7/3 20:30
 */
package vo

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"mime/multipart"
	"r0Website-server/models/bo"
	"time"
)

// AdminArticleAddFileVo 通过上传文件增加文章
type AdminArticleAddFileVo struct {
	AdminArticleAddMetaVo
	File *multipart.FileHeader `form:"file"` // 上传的文件
}

// AdminArticleAddFormVo 文章增加(form表单在网页进行编辑)Vo     "tui-editor": "^1.4.10",
type AdminArticleAddFormVo struct {
	AdminArticleAddMetaVo
	Markdown string `form:"markdown"` // 文章的md数据
}

// AdminArticleAddMetaVo admin权限下article增加功能的元信息模型
type AdminArticleAddMetaVo struct {
	Title      string   `form:"title"`      // 文章标题
	Author     string   `form:"author"`     // 主人
	Synopsis   string   `form:"synopsis"`   // 备注
	Tags       []string `form:"tags"`       // 标签
	Categories []string `form:"categories"` // 分类
	DraftFlag  bool     `form:"draft_flag"` // 是否为草稿
	Overhead   bool     `form:"overhead"`   // 是否顶置
	PicUrl     string   `form:"pic_url"`    // 图片的链接
}

// BaseArticleSetPVResultVo 设置pv之后返回的模型
type BaseArticleSetPVResultVo struct {
	MatchedCount  int64  `json:"matched_count" bson:"matched_count"`   // The number of documents matched by the filter.
	ModifiedCount int64  `json:"modified_count" bson:"modified_count"` // The number of documents modified by the operation.
	UpsertedCount int64  `json:"upserted_count" bson:"upserted_count"` // The number of documents upserted by the operation.
	Id            string `json:"id" bson:"id"`
}

type baseParams struct {
	Lazy           bool        `json:"lazy" form:"lazy"`                         // 如果采用懒惰加载，则不返回实体内容
	UpdateTimeSort bo.TimeSort `json:"update_time_sort" form:"update_time_sort"` // 更新时间排序的方向
	CreateTimeSort bo.TimeSort `json:"create_time_sort" form:"create_time_sort"` // 创建时间排序的方向
	PageNumber     int64       `json:"page_number" form:"page_number"`           // 分页使用，页码，页码从1开始
	PageSize       int64       `json:"page_size" form:"page_size"`               // 分页使用，页大小
}

// AdminArticleAddFileResultVo 通过AdminArticleAddFileVo提交之后的返回模型
type AdminArticleAddFileResultVo struct {
	Title string             `json:"title" bson:"title"`       // 文章标题
	Id    primitive.ObjectID `json:"_id" bson:"_id,omitempty"` // Mongo 主键 _id
}

// AdminArticleAddFormResultVo 通过AdminArticleAddFormVo提交之后的返回模型
type AdminArticleAddFormResultVo struct {
	Title string             `json:"title" bson:"title"`       // 文章标题
	Id    primitive.ObjectID `json:"_id" bson:"_id,omitempty"` // Mongo 主键 _id
}

// BaseArticleSearchVo 搜索模型
type BaseArticleSearchVo struct {
	SearchText string `json:"search_text" form:"search_text"` // 模糊搜素的内容 允许空格
	Author     string `json:"author" form:"author"`           // 作者名
	baseParams
}

// BaseArticleSearchResultVo 模糊搜索返回的结果
type BaseArticleSearchResultVo struct {
	Articles   []SingleBaseArticleSearchResultVo `json:"articles"`    // 文章列表
	PageNumber int64                             `json:"page_number"` // 页码
	PageSize   int64                             `json:"page_size"`   // 页面大小
	AnsCount   int64                             `json:"ans_count"`   // 结果数量
	TotalCount int64                             `json:"total_count"` // 总数
	Msg        string                            `json:"msg"`         // 提示信息
}

// SingleBaseArticleSearchResultVo 单个模糊搜索返回的结果
type SingleBaseArticleSearchResultVo struct {
	Id             primitive.ObjectID `json:"_id" bson:"_id,omitempty"`               // Mongo 主键 _id
	Title          string             `json:"title" bson:"title"`                     // 文章标题
	Author         string             `json:"author" bson:"author"`                   // 作者
	PicUrl         string             `json:"pic_url" bson:"pic_url"`                 // 图片的链接
	Markdown       string             `json:"markdown" bson:"markdown"`               // md内容
	ArtLength      int64              `json:"art_length" bson:"art_length"`           // 文章长度
	ReadsNumber    int64              `json:"reads_number" bson:"reads_number"`       // 阅读数
	CommentsNumber int64              `json:"comments_number" bson:"comments_number"` // 评论数
	PraiseNumber   int64              `json:"praise_number" bson:"praise_number"`     // 点赞数
	Tags           []string           `json:"tags" bson:"tags"`                       // 标签
	Categories     []string           `json:"categories" bson:"categories"`           // 分类
	CreateTime     time.Time          `json:"create_time" bson:"create_time"`         // 创建时间
	UpdateTime     time.Time          `json:"update_time" bson:"update_time"`         // 更新时间
	Score          float64            `json:"score" bson:"score"`                     // mongo全文检索评分
}

// ArticleSearchByCategoryVo 通过分类搜索文章的模型
type ArticleSearchByCategoryVo struct {
	baseParams
	CategoryName string `json:"category_name"`
}

// ArticleDeleteParams 删除
type ArticleDeleteParams struct {
	Id string `json:"id" bson:"id"`
}

// ArticleDeleteRes 删除
type ArticleDeleteRes struct {
	Count int64 `json:"count"`
}
