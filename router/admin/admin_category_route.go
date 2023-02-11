package admin

import (
	"github.com/gin-gonic/gin"
	"r0Website-server/r0Ioc"
)

func InitCategoryFileRouter(r *gin.RouterGroup) {
	article := r0Ioc.R0Route.AdminCategoryController
	group := r.Group("category")
	{
		// 不使用 /*id 的匹配是因为不想处理前后的"/"
		group.POST("/archive", article.ArchiveArticle)
	}
}
