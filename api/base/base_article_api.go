// Package base
/**
 * @Author: r0
 * @Mail: boogieLing_o@qq.com
 * @Description: base的文章api
 * @File:  base_article_api
 * @Version: 1.0.0
 * @Date: 2022/7/5 16:33
 */
package base

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"r0Website-server/global"
	"r0Website-server/models/vo"
	"r0Website-server/service"
	"r0Website-server/utils/msg"
)

type ArticleController struct {
	ArticleService *service.ArticleService `R0Ioc:"true"`
}

// AddPraise 增加一次赞
func (article *ArticleController) AddPraise(c *gin.Context) {
	articleID := c.Param("id")
	if res, err := article.ArticleService.AddPraise(articleID); err != nil {
		global.Logger.Error(err)
		c.JSON(http.StatusBadRequest, msg.NewMsg().Failed(err.Error()))
		return
	} else {
		c.JSON(http.StatusOK, msg.NewMsg().Success(res))
	}
}

// AddPV 增加一次pv
func (article *ArticleController) AddPV(c *gin.Context) {
	articleID := c.Param("id")
	if res, err := article.ArticleService.AddPV(articleID); err != nil {
		global.Logger.Error(err)
		c.JSON(http.StatusBadRequest, msg.NewMsg().Failed(err.Error()))
		return
	} else {
		c.JSON(http.StatusOK, msg.NewMsg().Success(res))
	}
}

// ArticleSearch 模糊搜索文章内容，依赖分词冗杂，允许带空格
// 如果带id那就是精准查找
// ShouldBindJSON > ShouldBind
// Follow: https://ofstack.com/Golang/29196/gin-golang-web-development-model-binding-implementation-process-analysis.html
func (article *ArticleController) ArticleSearch(c *gin.Context) {
	var params vo.BaseArticleSearchVo
	articleID := c.Param("id")
	if err := c.ShouldBind(&params); err != nil {
		global.Logger.Error(err)
		c.JSON(http.StatusBadRequest, msg.NewMsg().Failed("查询参数异常"))
		return
	}
	result, err := article.ArticleService.ArticleBaseSearch(params, articleID)
	if err != nil {
		global.Logger.Error(err)
		c.JSON(http.StatusBadRequest, msg.NewMsg().Failed(err.Error()))
		return
	}
	c.JSON(http.StatusOK, msg.NewMsg().Success(result))
}

// ArticleInCategory 获取分类下的文章
func (article *ArticleController) ArticleInCategory(c *gin.Context) {
	var params vo.ArticleSearchByCategoryVo
	name := c.Param("name")
	if err := c.ShouldBind(&params); err != nil {
		global.Logger.Error(err)
		c.JSON(http.StatusBadRequest, msg.NewMsg().Failed("查询参数异常"))
		return
	}
	params.CategoryName = name
	result, err := article.ArticleService.ArticleInCategory(params)
	if err != nil {
		global.Logger.Error(err)
		c.JSON(http.StatusBadRequest, msg.NewMsg().Failed(err.Error()))
		return
	}
	c.JSON(http.StatusOK, msg.NewMsg().Success(result))
}
