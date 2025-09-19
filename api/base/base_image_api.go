package base

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"r0Website-server/models/po"
	"r0Website-server/service"
)

type PicBedImageController struct {
	ImageService *service.AlbumService `R0Ioc:"true"`
}

// UploadImage 上传图片（数据库记录）
func (c *PicBedImageController) UploadImage(ctx *gin.Context) {
	var img po.Image
	if err := ctx.ShouldBindJSON(&img); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := c.ImageService.ImageDao.UploadImage(&img)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "上传失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"id": res.InsertedID})
}

// GetImageDetail 获取图片信息
func (c *PicBedImageController) GetImageDetail(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "非法图片ID"})
		return
	}
	img, err := c.ImageService.ImageDao.GetImageByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "图片不存在"})
		return
	}
	ctx.JSON(http.StatusOK, img)
}

// ListImages 获取所有图片列表
func (c *PicBedImageController) ListImages(ctx *gin.Context) {
	imgs, err := c.ImageService.ImageDao.ListImages()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取失败"})
		return
	}
	ctx.JSON(http.StatusOK, imgs)
}

// FindImagesByTag 通过标签查询图片
func (c *PicBedImageController) FindImagesByTag(ctx *gin.Context) {
	tag := ctx.Param("tag")
	imgs, err := c.ImageService.ImageDao.FindImagesByTag(tag)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}
	ctx.JSON(http.StatusOK, imgs)
}

// SearchImageByName 通过关键词模糊查找图片
func (c *PicBedImageController) SearchImageByName(ctx *gin.Context) {
	kw := ctx.Param("kw")
	imgs, err := c.ImageService.ImageDao.SearchImagesByName(kw)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "搜索失败"})
		return
	}
	ctx.JSON(http.StatusOK, imgs)
}

// GetImageAlbums 获取当前图片被引用的所有图集
func (c *PicBedImageController) GetImageAlbums(ctx *gin.Context) {
	imageID, _ := primitive.ObjectIDFromHex(ctx.Param("id"))
	albums, err := c.ImageService.ImageDao.GetAllAlbumsOfImage(imageID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}
	ctx.JSON(http.StatusOK, albums)
}

// DeleteImage 删除图片记录
func (c *PicBedImageController) DeleteImage(ctx *gin.Context) {
	imageID, _ := primitive.ObjectIDFromHex(ctx.Param("id"))
	err := c.ImageService.ImageDao.DeleteImageByID(imageID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "图片已删除"})
}
