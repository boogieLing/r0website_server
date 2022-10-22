// Package admin
/**
 * @Author: r0
 * @Mail: boogieLing_o@qq.com
 * @Description: 管理员下的文章api
 * @File:  article
 * @Version: 1.0.0
 * @Date: 2022/7/4 20:22
 */
package admin

import (
	"github.com/gin-gonic/gin"
	"r0Website-server/r0Ioc"
)

func InitArticleFileRouter(r *gin.RouterGroup) {
	article := r0Ioc.R0Route.AdminArticleController
	group := r.Group("article")
	{
		// 不使用 /*id 的匹配是因为不想处理前后的"/"
		group.POST("", article.ArticleFormWay)            // 通过编辑的方式增加文章 无id自动生成
		group.POST(":id", article.ArticleFormWay)         // 通过编辑的方式增加文章 id是必选的
		group.POST("/upload", article.ArticleFileWay)     // 通过上传文件的方式增加文章 无id自动生成
		group.POST("/upload/:id", article.ArticleFileWay) // 通过上传文件的方式增加文章 id是必选的
	}
}
