// Package router
/**
 * @Author: r0
 * @Mail: boogieLing_o@qq.com
 * @Description: 初始化路由
 * @File:  route_init
 * @Version: 1.0.0
 * @Date: 2022/7/3 18:39
 */
package router

import (
	"github.com/gin-gonic/gin"
	"r0Website-server/middleware"
	"r0Website-server/r0Ioc"
	"r0Website-server/router/admin"
	"r0Website-server/router/base"
)

// Routers 配置路由，依赖gin
// 日志、跨域、JWT
// public的路由不需要鉴权 admin/login不需要鉴权（处于登录态才颁发JWT）
func Routers() *gin.Engine {
	engine := gin.Default()
	engine.Use(middleware.Logger())
	engine.Use(middleware.Cors())

	root := engine.Group("api")
	{
		base.InitBaseRouter(root)
		adminGroup := root.Group("admin")
		userController := r0Ioc.R0Route.AdminUserController
		adminGroup.POST("login", userController.Login)
		adminGroup.Use(middleware.Jwt())
		{
			// TODO 需要鉴权的admin接口
			admin.InitArticleFileRouter(adminGroup)
		}
	}
	return engine
}
