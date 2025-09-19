package admin

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"r0Website-server/utils/msg"
)

type UserController struct {
}

// Login 用户登录
func (u *UserController) Login(c *gin.Context) {
	c.JSON(http.StatusOK, msg.NewMsg().Success("Hello"))
}
