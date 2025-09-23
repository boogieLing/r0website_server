package base

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"r0Website-server/dao"
	"r0Website-server/global"
	"r0Website-server/models/po"
	"r0Website-server/service"
	"r0Website-server/utils"
	"time"
)

type PicBedImageController struct {
	ImageService      *service.AlbumService      `R0Ioc:"true"`
	COSClient         *utils.COSClient           `R0Ioc:"true"`
	ImageCategoryDao  *dao.ImageCategoryDao      `R0Ioc:"true"`
	ImageDao          *dao.ImageDao              `R0Ioc:"true"`
	TagDao            *dao.TagDao                `R0Ioc:"true"`
}

// UploadImage 上传图片（支持文件上传和数据库记录）
func (c *PicBedImageController) UploadImage(ctx *gin.Context) {
	// 获取上传的文件
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "请选择要上传的文件"})
		return
	}
	defer file.Close()

	// 文件大小验证（最大10MB）
	const maxFileSize = 10 * 1024 * 1024 // 10MB
	if header.Size > maxFileSize {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "文件大小不能超过10MB"})
		return
	}

	// 文件类型验证
	contentType := header.Header.Get("Content-Type")
	allowedTypes := []string{"image/jpeg", "image/jpg", "image/png", "image/gif", "image/webp"}
	isValidType := false
	for _, t := range allowedTypes {
		if contentType == t {
			isValidType = true
			break
		}
	}
	if !isValidType {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "不支持的文件类型，只允许: JPG, PNG, GIF, WebP"})
		return
	}

	// 获取图片信息（宽度、高度）
	imgConfig, format, err := image.DecodeConfig(file)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无法解析图片文件"})
		return
	}

	// 重置文件读取位置
	file.Seek(0, 0)

	// 生成COS对象键
	objectKey := c.COSClient.GenerateObjectKey(header.Filename)

	// 上传到腾讯云COS
	cosURL, err := c.COSClient.UploadFile(file, header, objectKey)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "上传文件失败"})
		return
	}

	// 获取其他表单参数
	name := ctx.PostForm("name")
	if name == "" {
		name = header.Filename
	}
	tags := ctx.PostFormArray("tags")

	// 创建图片记录
	image := &po.Image{
		Name:       name,
		CosURL:     cosURL,
		Width:      imgConfig.Width,
		Height:     imgConfig.Height,
		UploadedAt: time.Now(),
		Tags:       tags,
		EXIF:       make(map[string]string),
		Positions:  make(map[string]po.CategoryPosition),
	}

	// 添加到nexus分类（默认分类）
	nexusPosition := po.CategoryPosition{
		X:            0,
		Y:            0,
		Width:        float64(imgConfig.Width),
		Height:       float64(imgConfig.Height),
		GridX:        0,
		GridY:        0,
		GridSize:     10,
		Row:          0,
		Col:          0,
		LayoutMode:   "freeform",
		ZIndex:       0,
		IsVisible:    true,
		Version:      1,
		UpdatedAt:    time.Now(),
		CategoryName: "Nexus",
		SortOrder:    0,
		AddedAt:      time.Now(),
	}
	image.Positions["nexus"] = nexusPosition

	// 保存到数据库
	res, err := c.ImageService.ImageDao.UploadImage(image)
	if err != nil {
		// 如果数据库保存失败，删除COS文件
		if deleteErr := c.COSClient.DeleteFile(objectKey); deleteErr != nil {
			// 记录删除失败日志，但不影响主要错误返回
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "保存图片信息失败"})
		return
	}

	// 获取插入的图片ID
	imageID := res.InsertedID.(primitive.ObjectID)

	// 维护倒排索引：将图片添加到nexus分类
	if err := c.ImageCategoryDao.AddImageToCategory("nexus", imageID, 0); err != nil {
		global.Logger.Errorf("添加图片到nexus分类失败: %v", err)
		// 不中断主流程，只记录错误
	}

	// 维护标签倒排索引：同步更新标签
	if len(tags) > 0 {
		for _, tagName := range tags {
			if tag, err := c.TagDao.GetOrCreateTag(tagName, tagName, ""); err == nil {
				if err := c.TagDao.AddImageToTag(tag.ID, imageID, name); err != nil {
					global.Logger.Errorf("添加图片到标签 %s 失败: %v", tagName, err)
				}
			} else {
				global.Logger.Errorf("获取或创建标签 %s 失败: %v", tagName, err)
			}
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":     imageID,
		"name":   name,
		"url":    cosURL,
		"width":  imgConfig.Width,
		"height": imgConfig.Height,
		"format": format,
	})
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

// UpdateImagePosition 更新图片在分类中的位置
func (c *PicBedImageController) UpdateImagePosition(ctx *gin.Context) {
	imageID, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "非法图片ID"})
		return
	}

	var req struct {
		CategoryID string                 `json:"categoryId" binding:"required"`
		Position   po.CategoryPosition    `json:"position" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 设置更新时间
	req.Position.UpdatedAt = time.Now()

	// 更新位置信息
	err = c.ImageService.ImageDao.UpdateImageCategoryPosition(imageID, req.CategoryID, &req.Position)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "更新位置失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"msg": "位置更新成功"})
}

// RemoveImageFromCategory 从分类中移除图片
func (c *PicBedImageController) RemoveImageFromCategory(ctx *gin.Context) {
	imageID, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "非法图片ID"})
		return
	}

	categoryID := ctx.Query("categoryId")
	if categoryID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "分类ID不能为空"})
		return
	}

	err = c.ImageService.ImageDao.RemoveImageFromCategory(imageID, categoryID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "移除失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"msg": "已从分类中移除"})
}

// GetImagesByCategory 获取分类下的图片（使用倒排索引优化）
func (c *PicBedImageController) GetImagesByCategory(ctx *gin.Context) {
	categoryID := ctx.Param("categoryId")
	if categoryID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "分类ID不能为空"})
		return
	}

	// 获取分页参数
	page := 1
	pageSize := 20
	if p := ctx.Query("page"); p != "" {
		fmt.Sscanf(p, "%d", &page)
	}
	if ps := ctx.Query("pageSize"); ps != "" {
		fmt.Sscanf(ps, "%d", &pageSize)
	}

	// 使用倒排索引获取分类中的图片引用
	imageRefs, total, err := c.ImageCategoryDao.GetCategoryImages(categoryID, page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取分类图片失败"})
		return
	}

	// 获取图片详细信息
	var images []*po.Image
	for _, ref := range imageRefs {
		img, err := c.ImageDao.GetImageByID(ref.ImageID)
		if err == nil {
			images = append(images, img)
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"images":   images,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}
