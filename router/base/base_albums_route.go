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
	"r0Website-server/r0Ioc"

	"github.com/gin-gonic/gin"
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
		group.PUT("album/:id/image", album.AddOrUpdateImageRef)               // 添加/更新图片引用
		group.PUT("album/:id/image/:imageId/layout", album.UpdateImageLayout) // 更新布局
		group.DELETE("album/:id/image/:imageId", album.RemoveImageFromAlbum)  // 移除引用
		group.PUT("image/move", album.MoveImageToAnotherAlbum)                // 移动图片到另一个图集

		// 图片 Image 操作
		group.POST("image", image.UploadImage)                             // 上传图片
		group.GET("image/:id", image.GetImageDetail)                       // 获取图片详情
		group.GET("image", image.ListImages)                               // 获取所有图片
		group.GET("image/tag/:tag", image.FindImagesByTag)                 // 按标签查图
		group.GET("image/search/:kw", image.SearchImageByName)             // 模糊查图
		group.GET("image/:id/albums", image.GetImageAlbums)                // 查询在哪些图集中
		group.DELETE("image/:id", image.DeleteImage)                       // 删除图片
		group.PUT("image/:id/position", image.UpdateImagePosition)         // 更新图片在分类中的位置
		group.DELETE("image/:id/category", image.RemoveImageFromCategory)  // 从分类中移除图片
		group.GET("image/category/:categoryId", image.GetImagesByCategory) // 获取分类下的图片

		// 图片分类管理
		category := r0Ioc.R0Route.ImageCategoryController
		group.POST("category", category.CreateCategory)                       // 创建分类
		group.GET("category", category.ListCategories)                        // 获取所有分类
		group.GET("category/:id", category.GetCategory)                       // 获取分类详情
		group.PUT("category/:id", category.UpdateCategory)                    // 更新分类
		group.PUT("category/:id/layout", category.UpdateCategoryLayoutMode)   // 调整分类布局方式
		group.DELETE("category/:id", category.DeleteCategory)                 // 删除分类
		group.GET("category/:id/images", category.GetCategoryImages)          // 获取分类中的图片
		group.POST("category/:id/images", category.AddImageToCategory)        // 添加图片到分类
		group.DELETE("category/:id/images", category.RemoveImageFromCategory) // 从分类移除图片
		group.PUT("category/:id/images/sort", category.UpdateImageSortOrder)  // 更新图片排序
		group.PUT("category/:id/cover", category.SetCategoryCover)            // 设置分类封面

		// 标签管理
		tag := r0Ioc.R0Route.TagController
		group.POST("tag", tag.CreateTag)              // 创建标签
		group.GET("tag", tag.ListTags)                // 获取所有标签
		group.GET("tag/popular", tag.GetPopularTags)  // 获取热门标签
		group.GET("tag/search", tag.SearchTags)       // 搜索标签
		group.POST("tag/batch", tag.BatchCreateTags)  // 批量创建标签
		group.GET("tag/:id", tag.GetTag)              // 获取标签详情
		group.PUT("tag/:id", tag.UpdateTag)           // 更新标签
		group.DELETE("tag/:id", tag.DeleteTag)        // 删除标签
		group.GET("tag/:id/images", tag.GetTagImages) // 获取标签中的图片
	}
}
