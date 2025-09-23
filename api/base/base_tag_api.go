package base

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"r0Website-server/dao"
	"r0Website-server/global"
	"r0Website-server/models/po"
	"strconv"
)

type TagController struct {
	TagDao   *dao.TagDao   `R0Ioc:"true"`
	ImageDao *dao.ImageDao `R0Ioc:"true"`
}

// CreateTag 创建新标签
func (c *TagController) CreateTag(ctx *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		DisplayName string `json:"displayName"`
		Description string `json:"description"`
		Color       string `json:"color"`
		Category    string `json:"category"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tag := &po.Tag{
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Description: req.Description,
		Color:       req.Color,
		Category:    req.Category,
	}

	if err := c.TagDao.CreateTag(tag); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "创建标签失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"msg": "标签创建成功", "tagId": tag.ID})
}

// GetTag 获取标签详情
func (c *TagController) GetTag(ctx *gin.Context) {
	tagID := ctx.Param("id")
	if tagID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "标签ID不能为空"})
		return
	}

	tag, err := c.TagDao.GetTagByID(tagID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "标签不存在"})
		return
	}

	ctx.JSON(http.StatusOK, tag)
}

// ListTags 获取所有标签
func (c *TagController) ListTags(ctx *gin.Context) {
	category := ctx.Query("category")

	tags, err := c.TagDao.ListTags(category)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取标签列表失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"tags":  tags,
		"total": len(tags),
	})
}

// GetPopularTags 获取热门标签
func (c *TagController) GetPopularTags(ctx *gin.Context) {
	limit := 20
	if l := ctx.Query("limit"); l != "" {
		if lInt, err := strconv.Atoi(l); err == nil && lInt > 0 {
			limit = lInt
		}
	}

	category := ctx.Query("category")

	tagStats, err := c.TagDao.GetPopularTags(limit, category)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取热门标签失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"tags":  tagStats,
		"total": len(tagStats),
	})
}

// UpdateTag 更新标签信息
func (c *TagController) UpdateTag(ctx *gin.Context) {
	tagID := ctx.Param("id")
	if tagID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "标签ID不能为空"})
		return
	}

	var update bson.M
	if err := ctx.ShouldBindJSON(&update); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.TagDao.UpdateTag(tagID, update); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "更新标签失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"msg": "标签更新成功"})
}

// DeleteTag 删除标签
func (c *TagController) DeleteTag(ctx *gin.Context) {
	tagID := ctx.Param("id")
	if tagID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "标签ID不能为空"})
		return
	}

	if err := c.TagDao.DeleteTag(tagID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"msg": "标签删除成功"})
}

// GetTagImages 获取标签中的图片
func (c *TagController) GetTagImages(ctx *gin.Context) {
	tagID := ctx.Param("id")
	if tagID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "标签ID不能为空"})
		return
	}

	// 获取分页参数
	page := 1
	pageSize := 20
	if p := ctx.Query("page"); p != "" {
		if pInt, err := strconv.Atoi(p); err == nil && pInt > 0 {
			page = pInt
		}
	}
	if ps := ctx.Query("pageSize"); ps != "" {
		if psInt, err := strconv.Atoi(ps); err == nil && psInt > 0 {
			pageSize = psInt
		}
	}

	// 获取标签中的图片引用
	imageRefs, total, err := c.TagDao.GetTagImages(tagID, page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取标签图片失败"})
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

// SearchTags 搜索标签
func (c *TagController) SearchTags(ctx *gin.Context) {
	keyword := ctx.Query("keyword")
	if keyword == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "搜索关键词不能为空"})
		return
	}

	// 搜索标签名称或显示名称
	filter := bson.M{
		"$or": []bson.M{
			{"name": bson.M{"$regex": keyword, "$options": "i"}},
			{"displayName": bson.M{"$regex": keyword, "$options": "i"}},
		},
	}

	cursor, err := c.TagDao.Collection().Find(ctx, filter, options.Find().SetLimit(50).SetSort(bson.M{"imageCount": -1}))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "搜索标签失败"})
		return
	}
	defer cursor.Close(ctx)

	var tags []*po.Tag
	if err = cursor.All(ctx, &tags); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "解析搜索结果失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"tags":  tags,
		"total": len(tags),
	})
}

// BatchCreateTags 批量创建标签
func (c *TagController) BatchCreateTags(ctx *gin.Context) {
	var req struct {
		Tags []struct {
			Name        string `json:"name" binding:"required"`
			DisplayName string `json:"displayName"`
			Description string `json:"description"`
			Color       string `json:"color"`
			Category    string `json:"category"`
		} `json:"tags" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdTags := []string{}
	for _, tagReq := range req.Tags {
		tag := &po.Tag{
			Name:        tagReq.Name,
			DisplayName: tagReq.DisplayName,
			Description: tagReq.Description,
			Color:       tagReq.Color,
			Category:    tagReq.Category,
		}

		if err := c.TagDao.CreateTag(tag); err != nil {
			global.Logger.Errorf("创建标签 %s 失败: %v", tagReq.Name, err)
			continue
		}
		createdTags = append(createdTags, tag.ID)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":         "批量创建标签完成",
		"createdTags": createdTags,
		"count":       len(createdTags),
	})
}