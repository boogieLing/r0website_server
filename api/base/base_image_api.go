// Package base
/**
 * @Author: r0
 * @Mail: boogieLing_o@qq.com
 * @Description: 图片相关操作的API
 * @File:  base_image_api
 * @Version: 1.0.0
 * @Date: 2024/09/24
 */
package base

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"r0Website-server/models/vo"
	"r0Website-server/service"
	"r0Website-server/utils/msg"
)

type PicBedImageController struct {
	ImageService *service.ImageService `R0Ioc:"true"`
}

// UploadImage 上传图片（支持文件上传和数据库记录）
func (c *PicBedImageController) UploadImage(ctx *gin.Context) {
	// 获取上传的文件
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, msg.NewMsg().Failed("请选择要上传的文件"))
		return
	}
	defer file.Close()

	// 绑定其他表单参数
	var params vo.UploadImageVo
	if err := ctx.ShouldBind(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, msg.NewMsg().Failed("参数绑定失败"))
		return
	}

	// 调用服务层上传图片
	result, err := c.ImageService.UploadImage(file, header, params)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, msg.NewMsg().Failed(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, msg.NewMsg().Success(result))
}

// GetImageDetail 获取图片信息
func (c *PicBedImageController) GetImageDetail(ctx *gin.Context) {
	// 解析图片ID
	idStr := ctx.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, msg.NewMsg().Failed("非法图片ID"))
		return
	}

	// 调用服务层获取图片详情
	img, err := c.ImageService.GetImageDetail(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, msg.NewMsg().Failed("图片不存在"))
		return
	}

	ctx.JSON(http.StatusOK, msg.NewMsg().Success(img))
}

// ListImages 获取所有图片列表
func (c *PicBedImageController) ListImages(ctx *gin.Context) {
	// 调用服务层获取图片列表
	imgs, err := c.ImageService.ListImages()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, msg.NewMsg().Failed("获取失败"))
		return
	}

	ctx.JSON(http.StatusOK, msg.NewMsg().Success(imgs))
}

// FindImagesByTag 通过标签查询图片
func (c *PicBedImageController) FindImagesByTag(ctx *gin.Context) {
	tag := ctx.Param("tag")

	// 调用服务层查询图片
	imgs, err := c.ImageService.FindImagesByTag(tag)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, msg.NewMsg().Failed("查询失败"))
		return
	}

	ctx.JSON(http.StatusOK, msg.NewMsg().Success(imgs))
}

// SearchImageByName 通过关键词模糊查找图片
func (c *PicBedImageController) SearchImageByName(ctx *gin.Context) {
	kw := ctx.Param("kw")

	// 调用服务层搜索图片
	imgs, err := c.ImageService.SearchImagesByName(kw)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, msg.NewMsg().Failed("搜索失败"))
		return
	}

	ctx.JSON(http.StatusOK, msg.NewMsg().Success(imgs))
}

// GetImageAlbums 获取当前图片被引用的所有图集
func (c *PicBedImageController) GetImageAlbums(ctx *gin.Context) {
	// 解析图片ID
	imageID, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, msg.NewMsg().Failed("非法图片ID"))
		return
	}

	// 调用服务层获取图集
	albums, err := c.ImageService.GetImageAlbums(imageID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, msg.NewMsg().Failed("查询失败"))
		return
	}

	ctx.JSON(http.StatusOK, msg.NewMsg().Success(albums))
}

// DeleteImage 删除图片记录
func (c *PicBedImageController) DeleteImage(ctx *gin.Context) {
	// 解析图片ID
	imageID, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, msg.NewMsg().Failed("非法图片ID"))
		return
	}

	// 调用服务层删除图片
	err = c.ImageService.DeleteImage(imageID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, msg.NewMsg().Failed("删除失败"))
		return
	}

	ctx.JSON(http.StatusOK, msg.NewMsg().Success("图片已删除"))
}

// UpdateImagePosition 更新图片在分类中的位置
func (c *PicBedImageController) UpdateImagePosition(ctx *gin.Context) {
	// 解析图片ID
	imageID, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, msg.NewMsg().Failed("非法图片ID"))
		return
	}

	// 绑定请求参数
	var params vo.UpdateImagePositionVo
	if err := ctx.ShouldBindJSON(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, msg.NewMsg().Failed("参数绑定失败"))
		return
	}

	// 调用服务层更新位置
	err = c.ImageService.UpdateImagePosition(imageID, params.CategoryID, &params.Position)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, msg.NewMsg().Failed("更新位置失败"))
		return
	}

	ctx.JSON(http.StatusOK, msg.NewMsg().Success("位置更新成功"))
}

// RemoveImageFromCategory 从分类中移除图片
func (c *PicBedImageController) RemoveImageFromCategory(ctx *gin.Context) {
	// 解析图片ID
	imageID, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, msg.NewMsg().Failed("非法图片ID"))
		return
	}

	// 获取分类ID
	categoryID := ctx.Query("categoryId")
	if categoryID == "" {
		ctx.JSON(http.StatusBadRequest, msg.NewMsg().Failed("分类ID不能为空"))
		return
	}

	// 调用服务层移除图片
	err = c.ImageService.RemoveImageFromCategory(imageID, categoryID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, msg.NewMsg().Failed("移除失败"))
		return
	}

	ctx.JSON(http.StatusOK, msg.NewMsg().Success("已从分类中移除"))
}

// GetImagesByCategory 获取分类下的图片（使用倒排索引优化）
func (c *PicBedImageController) GetImagesByCategory(ctx *gin.Context) {
	// 绑定请求参数
	var params vo.ImagesByCategoryVo
	if err := ctx.ShouldBindUri(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, msg.NewMsg().Failed("参数绑定失败"))
		return
	}
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, msg.NewMsg().Failed("查询参数绑定失败"))
		return
	}

	// 设置默认值
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 20
	}
	if params.Sort == "" {
		params.Sort = "position"
	}
	if params.Order == "" {
		params.Order = "asc"
	}

	// 调用服务层获取分类图片
	result, err := c.ImageService.GetImagesByCategory(params.CategoryID, params.Page, params.PageSize, params.Sort, params.Order)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, msg.NewMsg().Failed("获取分类图片失败"))
		return
	}

	ctx.JSON(http.StatusOK, msg.NewMsg().Success(result))
}