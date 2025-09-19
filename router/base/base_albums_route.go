// Package base
/**
 * @Author: r0
 * @Mail: boogieLing_o@qq.com
 * @Description: 图床模块的路由定义
 * @File:  picbed_route.go
 * @Version: 1.0.0
 * @Date: 2025/07/05
 */
package base

import (
	"github.com/gin-gonic/gin"
	"r0Website-server/r0Ioc"
)

func InitPicBedRouter(Router *gin.RouterGroup) {
	album := r0Ioc.R0Route.PicBedAlbumController
	image := r0Ioc.R0Route.PicBedImageController

	group := Router.Group("picbed")
	{
		// 图集 Album 操作
		group.POST("album", album.CreateAlbum)                // 创建图集
		group.GET("album/:id", album.GetAlbumDetail)          // 获取图集详情
		group.GET("album", album.ListAlbums)                  // 图集列表
		group.GET("album/tag/:tag", album.FindAlbumsByTag)    // 按标签查询图集
		group.GET("album/author/:author", album.FindByAuthor) // 按作者查图集
		group.GET("album/search/:kw", album.SearchByKeyword)  // 模糊搜索图集
		group.PUT("album/:id", album.UpdateAlbum)             // 更新图集信息
		group.DELETE("album/:id", album.DeleteAlbum)          // 删除图集

		// 图集中图片引用与布局
		group.PUT("album/:albumId/image", album.AddOrUpdateImageRef)               // 添加/更新图片引用
		group.PUT("album/:albumId/image/:imageId/layout", album.UpdateImageLayout) // 更新布局
		group.DELETE("album/:albumId/image/:imageId", album.RemoveImageFromAlbum)  // 移除引用
		group.PUT("image/move", album.MoveImageToAnotherAlbum)                     // 移动图片到另一个图集

		// 图片 Image 操作
		group.POST("image", image.UploadImage)                 // 上传图片
		group.GET("image/:id", image.GetImageDetail)           // 获取图片详情
		group.GET("image", image.ListImages)                   // 获取所有图片
		group.GET("image/tag/:tag", image.FindImagesByTag)     // 按标签查图
		group.GET("image/search/:kw", image.SearchImageByName) // 模糊查图
		group.GET("image/:id/albums", image.GetImageAlbums)    // 查询在哪些图集中
		group.DELETE("image/:id", image.DeleteImage)           // 删除图片
	}
}
