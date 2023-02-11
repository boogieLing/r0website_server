package admin

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"r0Website-server/models/vo"
	"r0Website-server/service"
	"r0Website-server/utils/msg"
)

type CategoryController struct {
	CategoryService *service.CategoryService `R0Ioc:"true"`
}

// ArchiveArticle 关联文章
func (cc *CategoryController) ArchiveArticle(c *gin.Context) {
	var input vo.ArchiveArticleVo
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, msg.NewMsg().Failed("参数异常"))
		return
	}
	result, err := cc.CategoryService.ArchiveArticle(input.ArticleId, input.CategoryName)
	if err != nil {
		c.JSON(http.StatusBadRequest, msg.NewMsg().Failed(err.Error()))
		return
	}
	c.JSON(http.StatusOK, msg.NewMsg().Success(result))
}
