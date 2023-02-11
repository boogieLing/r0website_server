package base

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"r0Website-server/service"
	"r0Website-server/utils/msg"
)

type CategoryController struct {
	CategoryService *service.CategoryService `R0Ioc:"true"`
}

func (cc *CategoryController) All(c *gin.Context) {
	result, err := cc.CategoryService.All()
	if err != nil {
		c.JSON(http.StatusBadRequest, msg.NewMsg().Failed(err.Error()))
		return
	}
	c.JSON(http.StatusOK, msg.NewMsg().Success(result))
	return
}
