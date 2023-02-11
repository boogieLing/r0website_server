package service

import (
	"go.mongodb.org/mongo-driver/mongo"
	"r0Website-server/dao"
	"r0Website-server/models/bo"
	"r0Website-server/models/vo"
)

type CategoryService struct {
	CategoryDao *dao.CategoryDao `R0Ioc:"true"`
}

// All 所有分类
func (cs *CategoryService) All() (*vo.CategorySearchResultVo, error) {
	return cs.CategoryDao.AllCategories()
}

// ArchiveArticle 关联文章
func (cs *CategoryService) ArchiveArticle(articleId, categoryName string) (*mongo.UpdateResult, error) {
	if articleId == "" {
		return nil, &bo.NullError{NullField: "articleId"}
	}
	if categoryName == "" {
		return nil, &bo.NullError{NullField: "categoryName"}
	}
	return cs.CategoryDao.ArchiveArticle(articleId, categoryName)
}
