// Package base
/**
 * @Author: r0
 * @Mail: boogieLing_o@qq.com
 * @Description: 基础接口
 * @File:  base
 * @Version: 1.0.0
 * @Date: 2022/7/3 18:44
 */
package base

import (
	"github.com/gin-gonic/gin"
	"r0Website-server/r0Ioc"
)

func InitBaseRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	userController := r0Ioc.R0Route.BaseUserController
	{
		Router.POST("login", userController.Login)
		Router.POST("register", userController.Register)
		InitBaseArticleRouter(Router)
		InitPicBedRouter(Router) // 添加图床路由
	}
	return Router
}
