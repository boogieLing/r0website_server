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
	"r0Website-server/global"
	"r0Website-server/middleware"
	"r0Website-server/r0Ioc"
	"r0Website-server/router/admin"
	"r0Website-server/router/base"

	"github.com/gin-gonic/gin"
)

// Routers 配置路由，依赖gin
// 日志、跨域、JWT
// public的路由不需要鉴权 admin/login不需要鉴权（处于登录态才颁发JWT）
func Routers() *gin.Engine {
	// 初始化默认图片分类
	albumService := r0Ioc.R0Route.PicBedAlbumController.AlbumService
	if albumService != nil {
		if err := albumService.InitDefaultCategories(); err != nil {
			global.Logger.Errorf("初始化默认图片分类失败: %v", err)
		} else {
			global.Logger.Infoln("✅ 默认图片分类初始化完成")
		}
	} else {
		global.Logger.Error("AlbumService 未初始化，跳过默认图片分类初始化")
	}

	engine := gin.Default()
	engine.MaxMultipartMemory = 64 << 20 // 允许更大的 multipart 表单
	engine.Use(middleware.Logger())
	engine.Use(middleware.Cors())

	root := engine.Group("api")
	{
		baseGroup := root.Group("base")
		{
			base.InitBaseRouter(baseGroup)
			base.InitCategoryFileRouter(baseGroup)
		}
		adminGroup := root.Group("admin")
		userController := r0Ioc.R0Route.AdminUserController
		adminGroup.POST("login", userController.Login)
		adminGroup.Use(middleware.Jwt())
		{
			// TODO 需要鉴权的admin接口
			admin.InitArticleFileRouter(adminGroup)
			admin.InitCategoryFileRouter(adminGroup)
		}
	}
	return engine
}
