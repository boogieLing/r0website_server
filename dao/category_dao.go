package dao

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"r0Website-server/global"
	"r0Website-server/models/bo"
	"r0Website-server/models/po"
	"r0Website-server/models/vo"
)

type CategoryDao struct {
	*BasicDaoMongo `R0Ioc:"true"`
}

func (*CategoryDao) CollectionName() string {
	return "categories"
}
func (cd *CategoryDao) Collection() *mongo.Collection {
	return cd.Mdb.Collection(cd.CollectionName())
}

// CategorySearch 查询某一分类
func (cd *CategoryDao) CategorySearch(name string) (*po.Category, error) {
	result := &vo.CategorySearchResultVo{}
	filter := bson.D{{"name", name}}
	cursor, err := cd.Collection().Find(context.TODO(), filter, nil)
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context.TODO(), &result.Categories); err != nil {
		global.Logger.Error(err)
		return nil, err
	}
	// defer 关闭游标
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			global.Logger.Error(err)
		}
	}(cursor, context.TODO())
	if len(result.Categories) >= 1 {
		return &result.Categories[0], nil
	} else {
		return &po.Category{}, nil
	}
}

// AddCategory 增加一种分类
func (cd *CategoryDao) AddCategory(name string) (*mongo.InsertOneResult, error) {
	var insertResult *mongo.InsertOneResult
	input := &po.Category{
		Name:       name,
		Count:      0,
		ArticleIds: []string{},
	}
	insertResult, err := cd.Collection().InsertOne(context.TODO(), input)
	if err != nil {
		global.Logger.Error(err)
		return nil, &bo.UniqueError{UniqueField: "category->_id", Msg: input.Id.Hex(), Count: 1}
	}
	return insertResult, nil
}

// AllCategories 所有的分类
func (cd *CategoryDao) AllCategories() (*vo.CategorySearchResultVo, error) {
	res := &vo.CategorySearchResultVo{}
	cursor, err := cd.Collection().Find(context.TODO(), bson.D{})
	if err != nil {
		global.Logger.Error(err)
		return nil, err
	}
	if err = cursor.All(context.TODO(), &res.Categories); err != nil {
		global.Logger.Error(err)
		return nil, err
	}
	res.TotalCount = int64(len(res.Categories))
	return res, nil
}

// ArchiveArticle 归档一篇文章
func (cd *CategoryDao) ArchiveArticle(articleId string, name string) (*mongo.UpdateResult, error) {
	if exists, err := cd.ExistsCategory(name); err != nil {
		return nil, err
	} else if !exists {
		if addResult, addErr := cd.AddCategory(name); addErr != nil {
			return nil, err
		} else {
			global.Logger.Infof("new Category:%s-%v", name, addResult.InsertedID)
		}
	}
	filter := bson.D{{"name", name}}
	update := bson.D{
		{"$addToSet", bson.D{{"article_ids", articleId}}},
		{"$inc", bson.D{{"count", 1}}},
	}
	result, err := cd.Collection().UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// ExistsCategory 该分类是否存在
func (cd *CategoryDao) ExistsCategory(name string) (bool, error) {
	filter := bson.D{{"name", name}}
	count, err := cd.Collection().CountDocuments(context.TODO(), filter)
	if err != nil {
		return false, err
	}
	if count >= 1 {
		return true, nil
	} else {
		return false, nil
	}
}
