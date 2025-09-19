package base

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"r0Website-server/models/po"
	"r0Website-server/service"
)

type PicBedAlbumController struct {
	AlbumService *service.AlbumService `R0Ioc:"true"`
}

// CreateAlbum 创建图集
func (c *PicBedAlbumController) CreateAlbum(ctx *gin.Context) {
	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Author      string `json:"author"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := c.AlbumService.CreateNewAlbum(req.Title, req.Description, req.Author)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"id": id})
}

// GetAlbumDetail 获取图集详情
func (c *PicBedAlbumController) GetAlbumDetail(ctx *gin.Context) {
	id, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "非法图集ID"})
		return
	}
	res, err := c.AlbumService.GetAlbumDetail(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "图集不存在"})
		return
	}
	ctx.JSON(http.StatusOK, res)
}

// ListAlbums 获取所有图集
func (c *PicBedAlbumController) ListAlbums(ctx *gin.Context) {
	res, err := c.AlbumService.AlbumDao.ListAlbums()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取失败"})
		return
	}
	ctx.JSON(http.StatusOK, res)
}

// FindAlbumsByTag 按标签获取图集
func (c *PicBedAlbumController) FindAlbumsByTag(ctx *gin.Context) {
	tag := ctx.Param("tag")
	res, err := c.AlbumService.AlbumDao.FindAlbumsByTag(tag)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}
	ctx.JSON(http.StatusOK, res)
}

// FindByAuthor 按作者查询图集
func (c *PicBedAlbumController) FindByAuthor(ctx *gin.Context) {
	author := ctx.Param("author")
	res, err := c.AlbumService.AlbumDao.FindAlbumsByAuthor(author)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}
	ctx.JSON(http.StatusOK, res)
}

// SearchByKeyword 关键词搜索图集
func (c *PicBedAlbumController) SearchByKeyword(ctx *gin.Context) {
	kw := ctx.Param("kw")
	res, err := c.AlbumService.AlbumDao.SearchAlbumsByKeyword(kw)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "搜索失败"})
		return
	}
	ctx.JSON(http.StatusOK, res)
}

// UpdateAlbum 更新图集信息
func (c *PicBedAlbumController) UpdateAlbum(ctx *gin.Context) {
	id, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "非法ID"})
		return
	}
	var album po.Album
	if err := ctx.ShouldBindJSON(&album); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	album.ID = id
	err = c.AlbumService.AlbumDao.UpdateAlbum(&album)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "更新成功"})
}

// DeleteAlbum 删除图集
func (c *PicBedAlbumController) DeleteAlbum(ctx *gin.Context) {
	id, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "非法ID"})
		return
	}
	err = c.AlbumService.AlbumDao.DeleteAlbum(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "图集已删除"})
}

// AddOrUpdateImageRef 添加或更新图集中的图片引用
func (c *PicBedAlbumController) AddOrUpdateImageRef(ctx *gin.Context) {
	albumID, _ := primitive.ObjectIDFromHex(ctx.Param("albumId"))
	var ref po.AlbumImageRef
	if err := ctx.ShouldBindJSON(&ref); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := c.AlbumService.AlbumDao.AddOrUpdateImageRef(albumID, ref.ImageID, &ref)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "保存失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "操作成功"})
}

// UpdateImageLayout 更新图片在图集中的布局
func (c *PicBedAlbumController) UpdateImageLayout(ctx *gin.Context) {
	albumID, _ := primitive.ObjectIDFromHex(ctx.Param("albumId"))
	imageID, _ := primitive.ObjectIDFromHex(ctx.Param("imageId"))
	var layout struct {
		Position    *po.AlbumPosition `json:"position"`
		Caption     string            `json:"caption"`
		Description string            `json:"description"`
	}
	if err := ctx.ShouldBindJSON(&layout); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := c.AlbumService.UpdateImageLayout(albumID, imageID, layout.Position, layout.Caption, layout.Description)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "布局已更新"})
}

// RemoveImageFromAlbum 从图集中移除图像
func (c *PicBedAlbumController) RemoveImageFromAlbum(ctx *gin.Context) {
	albumID, _ := primitive.ObjectIDFromHex(ctx.Param("albumId"))
	imageID, _ := primitive.ObjectIDFromHex(ctx.Param("imageId"))
	err := c.AlbumService.AlbumDao.RemoveImageFromAlbum(albumID, imageID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "移除失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "已移除"})
}

// MoveImageToAnotherAlbum 移动图片到其他图集
func (c *PicBedAlbumController) MoveImageToAnotherAlbum(ctx *gin.Context) {
	var req struct {
		FromAlbumID string `json:"from_album_id"`
		ToAlbumID   string `json:"to_album_id"`
		ImageID     string `json:"image_id"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fromID, _ := primitive.ObjectIDFromHex(req.FromAlbumID)
	toID, _ := primitive.ObjectIDFromHex(req.ToAlbumID)
	imageID, _ := primitive.ObjectIDFromHex(req.ImageID)
	err := c.AlbumService.MoveImageBetweenAlbums(fromID, toID, imageID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "移动失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "移动成功"})
}
