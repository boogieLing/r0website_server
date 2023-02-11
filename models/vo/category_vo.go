package vo

import "r0Website-server/models/po"

// CategorySearchResultVo 查找分类的结果
type CategorySearchResultVo struct {
	Categories []po.Category `json:"categories"`
	TotalCount int64         `json:"total_count"` // 总数
}

// ArchiveArticleVo 关联文章模型
type ArchiveArticleVo struct {
	ArticleId    string `json:"article_id"`
	CategoryName string `json:"category_name"`
}
