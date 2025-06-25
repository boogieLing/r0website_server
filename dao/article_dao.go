// Package dao
/**
 * @Author: r0
 * @Mail: boogieLing_o@qq.com
 * @Description: 文章相关的DAO
 * @File:  article_dao
 * @Version: 1.0.0
 * @Date: 2022/7/30 03:03
 */
package dao

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"r0Website-server/global"
	"r0Website-server/models/bo"
	"r0Website-server/models/po"
	"r0Website-server/models/vo"
	"r0Website-server/utils"
)

type ArticleDao struct {
	*BasicDaoMongo `R0Ioc:"true"`
	CategoryDao    *CategoryDao `R0Ioc:"true"`
}

func (*ArticleDao) CollectionName() string {
	return "articles"
}
func (ad *ArticleDao) Collection() *mongo.Collection {
	return ad.Mdb.Collection(ad.CollectionName())
}

// AddPraise 增加一次赞赏
func (ad *ArticleDao) AddPraise(id string) (*mongo.UpdateResult, error) {
	bsonId, err := primitive.ObjectIDFromHex(utils.String2HexString24(id))
	if err != nil {
		return nil, err
	}
	filter := bson.D{{"_id", bsonId}}
	update := bson.D{{"$inc", bson.D{{"praise_number", 1}}}}
	if res, err := ad.Collection().UpdateOne(context.TODO(), filter, update); err != nil {
		global.Logger.Error(err)
		return nil, err
	} else {
		return res, nil
	}
}

// AddPV 增加一次PV
func (ad *ArticleDao) AddPV(id string) (*mongo.UpdateResult, error) {
	bsonId, err := primitive.ObjectIDFromHex(utils.String2HexString24(id))
	if err != nil {
		return nil, err
	}
	filter := bson.D{{"_id", bsonId}}
	update := bson.D{{"$inc", bson.D{{"reads_number", 1}}}}
	if res, err := ad.Collection().UpdateOne(context.TODO(), filter, update); err != nil {
		global.Logger.Error(err)
		return nil, err
	} else {
		return res, nil
	}
}

// CreateArticle  增加文章
func (ad *ArticleDao) CreateArticle(input *po.Article) (ans *mongo.InsertOneResult, err error) {
	var insertResult *mongo.InsertOneResult
	insertResult, err = ad.Collection().InsertOne(context.TODO(), input)
	if err != nil {
		global.Logger.Error(err)
		return nil, &bo.UniqueError{UniqueField: "article->_id", Msg: input.Id.Hex(), Count: 1}
	}
	articleId := insertResult.InsertedID.(primitive.ObjectID)
	// 建立分类的倒排
	for _, category := range input.Categories {
		_, err := ad.CategoryDao.ArchiveArticle(articleId.Hex(), category)
		if err != nil {
			return nil, err
		}
	}
	return insertResult, nil
}

// ArticleBaseSearch 基础权限的文章搜索功能
// FOLLOWS:
// - https://www.mongodb.com/docs/drivers/go/current/fundamentals/crud/read-operations/text/
// - https://www.mongodb.com/docs/drivers/go/current/fundamentals/crud/read-operations/limit/
// db.articles.find({$text:{$search:"xxx"}, author:"xxx"},{score:{$meta : "textScore"}})
//	.sort({"update_time":-1, score:{$meta : "textScore"}})
func (ad *ArticleDao) ArticleBaseSearch(
	params vo.BaseArticleSearchVo, id string,
) (ans *vo.BaseArticleSearchResultVo, err error) {
	var result vo.BaseArticleSearchResultVo
	result.Articles = []vo.SingleBaseArticleSearchResultVo{}
	pageNumber := params.PageNumber
	pageSize := params.PageSize
	filter := ad.getArticleBaseSearchFilter(params, id)
	opts, err := ad.getArticleBaseSearchOption(params, id)
	if err != nil {
		return nil, err
	}
	// 防止全量搜索并构造分页, 页码从1开始，需要同时指定才能生效
	opts = ad.patchPageOption(&pageNumber, &pageSize, opts)
	global.Logger.Infof("ArticleBaseSearch -> Mongo: \n\t[ %+v | %+v ]", filter, opts)
	cursor, err := ad.Collection().Find(context.TODO(), filter, opts)
	if err != nil {
		global.Logger.Error(err)
		return nil, err
	}
	if err = cursor.All(context.TODO(), &result.Articles); err != nil {
		global.Logger.Error(err)
		return nil, err
	}
	for index, val := range result.Articles {
		result.Articles[index].UpdateTime = val.UpdateTime.Local()
		result.Articles[index].CreateTime = val.CreateTime.Local()
		if params.Lazy {
			result.Articles[index].Markdown = ""
		}
	}
	result.PageNumber = pageNumber
	result.PageSize = pageSize
	result.AnsCount = int64(len(result.Articles))
	result.TotalCount = ad.CountDocuments(filter)
	// defer 关闭游标
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			global.Logger.Error(err)
		}
	}(cursor, context.TODO())
	return &result, nil
}

// ArticleInCategory 某一分类下的文章
func (ad *ArticleDao) ArticleInCategory(
	params vo.ArticleSearchByCategoryVo,
) (*vo.BaseArticleSearchResultVo, error) {
	category, err := ad.CategoryDao.CategorySearch(params.CategoryName)
	if err != nil {
		return nil, err
	}
	matchIds := make([]primitive.ObjectID, len(category.ArticleIds))
	for index, id := range category.ArticleIds {
		matchIds[index], _ = primitive.ObjectIDFromHex(id)
	}
	var result vo.BaseArticleSearchResultVo
	result.Articles = []vo.SingleBaseArticleSearchResultVo{}
	pageNumber := params.PageNumber
	pageSize := params.PageSize
	filter := bson.D{{"_id", bson.D{{"$in", matchIds}}}}
	// 防止全量搜索并构造分页, 页码从1开始，需要同时指定才能生效
	opts, err := ad.getArticleBaseSearchOption(vo.BaseArticleSearchVo{
		SearchText: "",
		Author:     "",
		BaseParams: vo.BaseParams{
			Lazy:           params.Lazy,
			UpdateTimeSort: params.UpdateTimeSort,
			CreateTimeSort: params.CreateTimeSort,
			PageNumber:     params.PageNumber,
			PageSize:       params.PageSize,
		},
	}, "")
	if err != nil {
		return nil, err
	}
	opts = ad.patchPageOption(&pageNumber, &pageSize, opts)

	cursor, err := ad.Collection().Find(context.TODO(), filter, opts)
	if err != nil {
		global.Logger.Error(err)
		return nil, err
	}
	if err = cursor.All(context.TODO(), &result.Articles); err != nil {
		global.Logger.Error(err)
		return nil, err
	}
	for index, val := range result.Articles {
		result.Articles[index].UpdateTime = val.UpdateTime.Local()
		result.Articles[index].CreateTime = val.CreateTime.Local()
		if params.Lazy {
			result.Articles[index].Markdown = ""
		}
	}
	result.PageNumber = pageNumber
	result.PageSize = pageSize
	result.AnsCount = int64(len(result.Articles))
	result.TotalCount = ad.CountDocuments(filter)
	// defer 关闭游标
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			global.Logger.Error(err)
		}
	}(cursor, context.TODO())
	return &result, nil
}

// CountDocuments 统计文档总数
func (ad *ArticleDao) CountDocuments(filter interface{}) int64 {
	if ans, err := ad.Collection().CountDocuments(context.TODO(), filter); err != nil {
		global.Logger.Error(err)
		return -1
	} else {
		return ans
	}
}

// ArticleCountByUUid 统计有多少文章有此uuid
func (ad *ArticleDao) ArticleCountByUUid(uuid string) int64 {
	filter := bson.M{"uuid": uuid}
	count, err := ad.Collection().CountDocuments(context.TODO(), filter)
	if err != nil {
		global.Logger.Error(err)
	}
	return count
}

// patchPageParams 修正分页Option，并获取pageSkip
func (ad *ArticleDao) patchPageOption(
	pageNumber, pageSize *int64, opts *options.FindOptions,
) *options.FindOptions {
	if *pageNumber <= 0 {
		*pageNumber = 1
	}
	if *pageSize <= 0 {
		*pageSize = 10
	}
	if *pageSize >= 200 {
		*pageSize = 200
	}
	if *pageNumber != 0 && *pageSize != 0 {
		pageSkip := (*pageNumber - 1) * *pageSize
		opts = opts.SetLimit(*pageSize).SetSkip(pageSkip)
	}
	return opts
}

// getArticleBaseSearchOption 构造ArticleBaseSearch的选项
func (ad *ArticleDao) getArticleBaseSearchOption(
	params vo.BaseArticleSearchVo, id string,
) (*options.FindOptions, error) {
	searchText := params.SearchText
	sort := bson.D{}
	projection := bson.D{}
	opts := options.Find()

	// 构造排序bson
	if params.UpdateTimeSort.SortFlag == true && params.CreateTimeSort.SortFlag == true {
		return nil, errors.New("ArticleBaseSearch: " + "不能同时指定UpdateTime和CreateTime的排序")
	}
	if params.UpdateTimeSort.SortFlag == true {
		sort = bson.D{{"update_time", params.UpdateTimeSort.SortDirection}}
		opts = opts.SetSort(sort)
	}
	if params.CreateTimeSort.SortFlag == true {
		sort = bson.D{{"create_time", params.CreateTimeSort.SortDirection}}
		opts = opts.SetSort(sort)
	}
	// 如果使用id检索，其他检索全部失效，并不应该在排序段增加score
	if searchText != "" && id == "" {
		// 如果包含模糊搜素，那么需要在sort段和投影段增加条件
		sort = append(sort, bson.E{Key: "score", Value: bson.D{{"$meta", "textScore"}}})
		// 构造投影bson
		projection = append(projection, bson.E{Key: "score", Value: bson.D{{"$meta", "textScore"}}})
		// 这里会引起opts的二次SetSort，但目前还没发现问题
		opts = opts.SetSort(sort).SetProjection(projection)
	}
	return opts, nil
}

// getArticleBaseSearchFilter 构造ArticleBaseSearch的过滤
func (ad *ArticleDao) getArticleBaseSearchFilter(
	params vo.BaseArticleSearchVo, id string,
) bson.D {
	searchText := params.SearchText
	author := params.Author
	// 是一个补丁，防止出现WriteNull错误，丑陋的解决方法
	filter := bson.D{}
	// 构造搜索/过滤bson
	if author != "" {
		filter = append(filter, bson.E{Key: "author", Value: author})
	}
	// 如果使用id检索，其他检索全部失效
	if id != "" {
		if bsonId, err := primitive.ObjectIDFromHex(utils.String2HexString24(id)); err != nil {
			global.Logger.Error(err)
		} else {
			filter = bson.D{{"_id", bsonId}}
			global.Logger.Info("使用id检索，其他检索条件全部失效")
		}
	}
	// 如果使用id检索，其他检索全部失效，并不应该在排序段增加score
	if searchText != "" && id == "" {
		filter = append(filter, bson.E{Key: "$text", Value: bson.D{{"$search", searchText}}})
	}
	return filter
}

// DeleteArticle 删除文章
func (ad *ArticleDao) DeleteArticle(id string) (int64, error) {
	if bsonId, err := primitive.ObjectIDFromHex(utils.String2HexString24(id)); err != nil {
		global.Logger.Error(err)
		return 0, err
	} else {
		filter := bson.D{{"_id", bsonId}}
		deleteRes, err := ad.Collection().DeleteOne(context.TODO(), filter)
		if err != nil {
			global.Logger.Error(err)
			return 0, err
		} else {
			return deleteRes.DeletedCount, nil
		}
	}
}
