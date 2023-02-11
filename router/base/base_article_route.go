// Package base
/**
 * @Author: r0
 * @Mail: boogieLing_o@qq.com
 * @Description:
 * @File:  article_route
 * @Version: 1.0.0
 * @Date: 2022/7/5 16:36
 */
package base

import (
	"github.com/gin-gonic/gin"
	"r0Website-server/r0Ioc"
)

func InitBaseArticleRouter(Router *gin.RouterGroup) {
	article := r0Ioc.R0Route.BaseArticleController
	group := Router.Group("article")
	{
		group.GET("", article.ArticleSearch)                   // 模糊搜素
		group.GET(":id", article.ArticleSearch)                // id精确搜索
		group.PUT(":id/pv", article.AddPV)                     // 设置PV
		group.PUT(":id/praise", article.AddPraise)             // 增加一次praise
		group.GET("category/:name", article.ArticleInCategory) // 某一分类下的文章
	}
}
