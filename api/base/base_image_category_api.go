package base

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"r0Website-server/dao"
	"r0Website-server/models/po"
	"strconv"
	"time"
)

type ImageCategoryController struct {
	ImageCategoryDao *dao.ImageCategoryDao `R0Ioc:"true"`
	ImageDao         *dao.ImageDao         `R0Ioc:"true"`
}

// CreateCategory 创建新的图片分类
func (c *ImageCategoryController) CreateCategory(ctx *gin.Context) {
	var req struct {
		ID          string `json:"id" binding:"required"`   // 分类ID，如 "nexus", "stillness"
		Name        string `json:"name" binding:"required"` // 分类名称
		Description string `json:"description"`             // 分类描述
		LayoutMode  string `json:"layoutMode"`              // 布局模式
		GridSize    int    `json:"gridSize"`                // 网格大小
		AutoArrange bool   `json:"autoArrange"`             // 是否自动排列
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category := &po.ImageCategory{
		ID:          req.ID,
		Name:        req.Name,
		Description: req.Description,
		Settings: po.CategorySettings{
			LayoutMode:    req.LayoutMode,
			GridSize:      req.GridSize,
			DefaultWidth:  200,
			DefaultHeight: 280,
			AutoArrange:   req.AutoArrange,
		},
	}

	if err := c.ImageCategoryDao.CreateCategory(category); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "创建分类失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"msg": "分类创建成功"})
}

// GetCategory 获取分类详情
func (c *ImageCategoryController) GetCategory(ctx *gin.Context) {
	categoryID := ctx.Param("id")
	if categoryID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "分类ID不能为空"})
		return
	}

	category, err := c.ImageCategoryDao.GetCategoryByID(categoryID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "分类不存在"})
		return
	}

	ctx.JSON(http.StatusOK, category)
}

// ListCategories 获取所有分类
func (c *ImageCategoryController) ListCategories(ctx *gin.Context) {
	categories, err := c.ImageCategoryDao.ListCategories()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取分类列表失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"categories": categories,
		"total":      len(categories),
	})
}

// UpdateCategory 更新分类信息
func (c *ImageCategoryController) UpdateCategory(ctx *gin.Context) {
	categoryID := ctx.Param("id")
	if categoryID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "分类ID不能为空"})
		return
	}

	var update bson.M
	if err := ctx.ShouldBindJSON(&update); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.ImageCategoryDao.UpdateCategory(categoryID, update); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "更新分类失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"msg": "分类更新成功"})
}

// UpdateCategoryLayoutMode 更新分类的布局方式（layoutMode）
func (c *ImageCategoryController) UpdateCategoryLayoutMode(ctx *gin.Context) {
	categoryID := ctx.Param("id")
	if categoryID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "分类ID不能为空"})
		return
	}

	var req struct {
		LayoutMode string `json:"layoutMode" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	update := bson.M{
		"$set": bson.M{
			"settings.layoutMode": req.LayoutMode,
		},
	}

	if err := c.ImageCategoryDao.UpdateCategory(categoryID, update); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "更新分类布局方式失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"msg": "分类布局方式已更新"})
}

// DeleteCategory 删除分类
func (c *ImageCategoryController) DeleteCategory(ctx *gin.Context) {
	categoryID := ctx.Param("id")
	if categoryID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "分类ID不能为空"})
		return
	}

	if err := c.ImageCategoryDao.DeleteCategory(categoryID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"msg": "分类删除成功"})
}

// GetCategoryImages 获取分类中的图片
func (c *ImageCategoryController) GetCategoryImages(ctx *gin.Context) {
	categoryID := ctx.Param("id")
	if categoryID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "分类ID不能为空"})
		return
	}

	// 获取分页参数
	page := 1
	pageSize := 20
	if p := ctx.Query("page"); p != "" {
		page, _ = strconv.Atoi(p)
		if page < 1 {
			page = 1
		}
	}
	if ps := ctx.Query("pageSize"); ps != "" {
		pageSize, _ = strconv.Atoi(ps)
		if pageSize < 1 {
			pageSize = 20
		}
	}

	// 获取分类中的图片引用
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

// AddImageToCategory 添加图片到分类
func (c *ImageCategoryController) AddImageToCategory(ctx *gin.Context) {
	categoryID := ctx.Param("id")
	if categoryID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "分类ID不能为空"})
		return
	}

	var req struct {
		ImageID   string `json:"imageId" binding:"required"`
		SortOrder int    `json:"sortOrder"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	imageID, err := primitive.ObjectIDFromHex(req.ImageID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的图片ID"})
		return
	}

	if err := c.ImageCategoryDao.AddImageToCategory(categoryID, imageID, req.SortOrder); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"msg": "图片已添加到分类"})
}

// RemoveImageFromCategory 从分类中移除图片
func (c *ImageCategoryController) RemoveImageFromCategory(ctx *gin.Context) {
	categoryID := ctx.Param("id")
	imageID := ctx.Query("imageId")

	if categoryID == "" || imageID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "分类ID和图片ID不能为空"})
		return
	}

	objID, err := primitive.ObjectIDFromHex(imageID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的图片ID"})
		return
	}

	if err := c.ImageCategoryDao.RemoveImageFromCategory(categoryID, objID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"msg": "图片已从分类中移除"})
}

// UpdateImageSortOrder 更新图片在分类中的排序
func (c *ImageCategoryController) UpdateImageSortOrder(ctx *gin.Context) {
	categoryID := ctx.Param("id")
	if categoryID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "分类ID不能为空"})
		return
	}

	var req struct {
		ImageID   string `json:"imageId" binding:"required"`
		SortOrder int    `json:"sortOrder" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	imageID, err := primitive.ObjectIDFromHex(req.ImageID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的图片ID"})
		return
	}

	if err := c.ImageCategoryDao.UpdateImageSortOrder(categoryID, imageID, req.SortOrder); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"msg": "图片排序已更新"})
}

// SetCategoryCover 设置分类封面
func (c *ImageCategoryController) SetCategoryCover(ctx *gin.Context) {
	categoryID := ctx.Param("id")
	imageID := ctx.Query("imageId")

	if categoryID == "" || imageID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "分类ID和图片ID不能为空"})
		return
	}

	objID, err := primitive.ObjectIDFromHex(imageID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的图片ID"})
		return
	}

	// 检查图片是否在该分类中
	exists, err := c.ImageCategoryDao.IsImageInCategory(categoryID, objID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "检查图片失败"})
		return
	}
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "图片不在该分类中"})
		return
	}

	update := bson.M{
		"$set": bson.M{
			"coverImage": objID,
			"updatedAt":  time.Now(),
		},
	}

	if err := c.ImageCategoryDao.UpdateCategory(categoryID, update); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "设置封面失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"msg": "封面设置成功"})
}
