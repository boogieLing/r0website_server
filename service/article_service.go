// Package service
/**
 * @Author: r0
 * @Mail: boogieLing_o@qq.com
 * @Description: 文章功能
 * @File:  article
 * @Version: 1.0.0
 * @Date: 2022/7/4 20:20
 */
package service

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io/ioutil"
	"r0Website-server/dao"
	"r0Website-server/global"
	"r0Website-server/models/bo"
	"r0Website-server/models/po"
	"r0Website-server/models/vo"
	"r0Website-server/utils"
	"time"
)

type ArticleService struct {
	ArticleDao *dao.ArticleDao `R0Ioc:"true"`
}

const ArticleColl = "articles"

// ArticleADDFile 通过上传文件增加文章
func (article *ArticleService) ArticleADDFile(
	params vo.AdminArticleAddFileVo, id string,
) (ans *vo.AdminArticleAddFileResultVo, err error) {
	var result vo.AdminArticleAddFileResultVo
	var input po.Article
	// 获取文章内容
	open, err := params.File.Open()
	if err != nil {
		global.Logger.Error(err)
		return nil, err
	}
	content, err := ioutil.ReadAll(open)
	if err != nil {
		global.Logger.Error(err)
		return nil, errors.New("ArticleADDFile: 获取文件内容失败")
	}
	input.Markdown = string(content)
	updateArticleMetaByParams(&input, params, id)
	insertResult, err := article.ArticleDao.CreateArticle(&input)
	if err != nil {
		global.Logger.Error(err)
	} else {
		result = vo.AdminArticleAddFileResultVo{Title: input.Title, Id: insertResult.InsertedID.(primitive.ObjectID)}
	}
	return &result, err
}

// ArticleADDForm 通过提交表单增加文章
func (article *ArticleService) ArticleADDForm(
	params vo.AdminArticleAddFormVo, id string,
) (ans *vo.AdminArticleAddFormResultVo, err error) {
	var result vo.AdminArticleAddFormResultVo
	var input po.Article
	input.Markdown = params.Markdown
	updateArticleMetaByParams(&input, params, id)
	insertResult, err := article.ArticleDao.CreateArticle(&input)
	if err != nil {
		global.Logger.Error(err)
	} else {
		result = vo.AdminArticleAddFormResultVo{Title: input.Title, Id: insertResult.InsertedID.(primitive.ObjectID)}
	}
	return &result, err
}

// ArticleBaseSearch 基础权限的文章搜索功能
// 可选作者，可选模糊内容，可选两种时间排序，可选id
// **如果使用id检索，其他检索全部失效**
// 分页 利用opt构造的skip和limit
func (article *ArticleService) ArticleBaseSearch(
	params vo.BaseArticleSearchVo, id string,
) (ans *vo.BaseArticleSearchResultVo, err error) {
	return article.ArticleDao.ArticleBaseSearch(params, id)
}

// updateArticleMetaByParams 用输入的参数更新文章元信息
func updateArticleMetaByParams(input *po.Article, params interface{}, uuid string) {
	switch params.(type) {
	case vo.AdminArticleAddFormVo, vo.AdminArticleAddFileVo:
		var meta vo.AdminArticleAddMetaVo
		if value, ok := params.(vo.AdminArticleAddFormVo); ok {
			meta = value.AdminArticleAddMetaVo
		} else if value, ok := params.(vo.AdminArticleAddFileVo); ok {
			meta = value.AdminArticleAddMetaVo
		}
		// input.Uuid = uuid
		var err error
		uuid = utils.String2HexString24(uuid)
		if input.Id, err = primitive.ObjectIDFromHex(uuid); err != nil {
			global.Logger.Error(err)
		}
		input.Title = meta.Title
		input.Author = meta.Author
		input.Synopsis = meta.Synopsis
		input.DeleteFlag = false
		input.DraftFlag = meta.DraftFlag
		// input.Detail = ""
		// input.Markdown = meta.Markdown
		input.Overhead = meta.Overhead
		input.ArtLength = 0
		input.ReadsNumber = 0
		input.CommentsNumber = 0
		input.PraiseNumber = 0
		input.Tags = meta.Tags
		input.Categories = meta.Categories
		var curTime = time.Now()
		input.UpdateTime = curTime
		input.CreateTime = curTime
		// 接下来“修补”模型的值
		// input.Detail = utils.Markdown2Html(input.Markdown)
		// input.ArtLength = len(input.Markdown)
		wordCounter := utils.WordCounter{}
		wordCounter.Stat(input.Markdown)
		input.ArtLength = int64(wordCounter.Total)
		input.MdWords = utils.WordSplitForSearching(input.Markdown)
		input.TitleWords = utils.WordSplitForSearching(input.Title)
	}

}

// checkAndPatchArticleUuid 检查并修补文章uuid的值
func (article *ArticleService) checkAndPatchArticleUuid(uuid string) (string, error) {
	var assignedUuid = true
	if uuid == "" {
		// 空id 将自动生成uuid
		uuid = utils.GenSonyflake()
		assignedUuid = false
	}
	if articleUuidCount := article.ArticleDao.ArticleCountByUUid(uuid); articleUuidCount > 0 {
		// 表示id重复 重复生成一次 两次Sonyflake碰撞的概率非常低
		// 但对于指定了idq的POST请求的碰撞 应该直接给UniqueError
		if assignedUuid {
			return "", &bo.UniqueError{UniqueField: "article-uuid", Msg: uuid, Count: articleUuidCount}
		} else {
			uuid = utils.GenSonyflake()
		}
	}
	return uuid, nil
}
