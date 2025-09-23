package dao

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"r0Website-server/global"
	"r0Website-server/models/po"
	"time"
)

type ImageCategoryDao struct {
	*BasicDaoMongo `R0Ioc:"true"`
}

func (*ImageCategoryDao) CollectionName() string {
	return "image_categories"
}

func (cd *ImageCategoryDao) Collection() *mongo.Collection {
	return cd.Mdb.Collection(cd.CollectionName())
}

// CreateCategory 创建分类
func (cd *ImageCategoryDao) CreateCategory(category *po.ImageCategory) error {
	category.CreatedAt = time.Now()
	category.UpdatedAt = time.Now()
	category.ImageCount = 0

	// 设置默认设置
	if category.Settings.LayoutMode == "" {
		category.Settings.LayoutMode = "freeform"
	}
	if category.Settings.GridSize == 0 {
		category.Settings.GridSize = 10
	}
	if category.Settings.DefaultWidth == 0 {
		category.Settings.DefaultWidth = 200
	}
	if category.Settings.DefaultHeight == 0 {
		category.Settings.DefaultHeight = 280
	}

	_, err := cd.Collection().InsertOne(context.TODO(), category)
	if err != nil {
		global.Logger.Errorf("❌ 创建图片分类失败: %v", err)
		return err
	}
	return nil
}

// GetCategoryByID 根据ID获取分类
func (cd *ImageCategoryDao) GetCategoryByID(categoryID string) (*po.ImageCategory, error) {
	var category po.ImageCategory
	err := cd.Collection().FindOne(context.TODO(), bson.M{"_id": categoryID}).Decode(&category)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("图片分类不存在")
		}
		global.Logger.Errorf("❌ 获取图片分类失败: %v", err)
		return nil, err
	}
	return &category, nil
}

// UpdateCategory 更新分类信息
func (cd *ImageCategoryDao) UpdateCategory(categoryID string, update bson.M) error {
	if _, ok := update["$set"]; !ok {
		update["$set"] = bson.M{}
	}
	update["$set"].(bson.M)["updatedAt"] = time.Now()

	_, err := cd.Collection().UpdateOne(context.TODO(), bson.M{"_id": categoryID}, update)
	if err != nil {
		global.Logger.Errorf("❌ 更新图片分类失败: %v", err)
		return err
	}
	return nil
}

// DeleteCategory 删除分类
func (cd *ImageCategoryDao) DeleteCategory(categoryID string) error {
	// 检查分类中是否有图片
	count, err := cd.GetCategoryImageCount(categoryID)
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("分类中还存在图片，无法删除")
	}

	_, err = cd.Collection().DeleteOne(context.TODO(), bson.M{"_id": categoryID})
	if err != nil {
		global.Logger.Errorf("❌ 删除图片分类失败: %v", err)
		return err
	}
	return nil
}

// ListCategories 获取所有分类
func (cd *ImageCategoryDao) ListCategories() ([]*po.ImageCategory, error) {
	cursor, err := cd.Collection().Find(context.TODO(), bson.M{}, options.Find().SetSort(bson.M{"name": 1}))
	if err != nil {
		global.Logger.Errorf("❌ 获取图片分类列表失败: %v", err)
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var categories []*po.ImageCategory
	if err = cursor.All(context.TODO(), &categories); err != nil {
		global.Logger.Errorf("❌ 解析图片分类列表失败: %v", err)
		return nil, err
	}
	return categories, nil
}

// AddImageToCategory 添加图片到分类（维护倒排索引）
func (cd *ImageCategoryDao) AddImageToCategory(categoryID string, imageID primitive.ObjectID, sortOrder int) error {
	// 获取当前最大排序序号
	maxSortOrder, err := cd.getMaxSortOrder(categoryID)
	if err != nil {
		return err
	}
	if sortOrder == 0 {
		sortOrder = maxSortOrder + 1
	}

	// 检查图片是否已存在
	exists, err := cd.IsImageInCategory(categoryID, imageID)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("图片已存在于该分类中")
	}

	// 添加到分类的图片列表中
	imageRef := po.CategoryImageRef{
		ImageID:   imageID,
		SortOrder: sortOrder,
		AddedAt:   time.Now(),
	}

	update := bson.M{
		"$push": bson.M{
			"images": imageRef,
		},
		"$inc": bson.M{
			"imageCount": 1,
		},
		"$set": bson.M{
			"updatedAt": time.Now(),
		},
	}

	_, err = cd.Collection().UpdateOne(context.TODO(), bson.M{"_id": categoryID}, update)
	if err != nil {
		global.Logger.Errorf("❌ 添加图片到分类失败: %v", err)
		return err
	}
	return nil
}

// RemoveImageFromCategory 从分类中移除图片
func (cd *ImageCategoryDao) RemoveImageFromCategory(categoryID string, imageID primitive.ObjectID) error {
	update := bson.M{
		"$pull": bson.M{
			"images": bson.M{"imageId": imageID},
		},
		"$inc": bson.M{
			"imageCount": -1,
		},
		"$set": bson.M{
			"updatedAt": time.Now(),
		},
	}

	result, err := cd.Collection().UpdateOne(context.TODO(), bson.M{"_id": categoryID}, update)
	if err != nil {
		global.Logger.Errorf("❌ 从分类中移除图片失败: %v", err)
		return err
	}

	if result.ModifiedCount == 0 {
		return fmt.Errorf("图片不在该分类中")
	}

	return nil
}

// GetCategoryImages 获取分类中的图片列表
func (cd *ImageCategoryDao) GetCategoryImages(categoryID string, page, pageSize int) ([]po.CategoryImageRef, int64, error) {
	// 获取分类信息
	category, err := cd.GetCategoryByID(categoryID)
	if err != nil {
		return nil, 0, err
	}
	total := int64(category.ImageCount)

	// 如果分页参数无效，返回所有图片
	if page <= 0 || pageSize <= 0 {
		return category.Images, total, nil
	}

	// 分页查询
	skip := (page - 1) * pageSize
	pipeline := []bson.M{
		{"$match": bson.M{"_id": categoryID}},
		{"$unwind": "$images"},
		{"$sort": bson.M{"images.sortOrder": 1}},
		{"$skip": skip},
		{"$limit": pageSize},
		{"$replaceRoot": bson.M{"newRoot": "$images"}},
	}

	cursor, err := cd.Collection().Aggregate(context.TODO(), pipeline)
	if err != nil {
		global.Logger.Errorf("❌ 获取分类图片失败: %v", err)
		return nil, 0, err
	}
	defer cursor.Close(context.TODO())

	var imageRefs []po.CategoryImageRef
	if err = cursor.All(context.TODO(), &imageRefs); err != nil {
		global.Logger.Errorf("❌ 解析分类图片失败: %v", err)
		return nil, 0, err
	}

	return imageRefs, total, nil
}

// GetCategoryImageCount 获取分类中的图片数量
func (cd *ImageCategoryDao) GetCategoryImageCount(categoryID string) (int, error) {
	category, err := cd.GetCategoryByID(categoryID)
	if err != nil {
		return 0, err
	}
	return category.ImageCount, nil
}

// UpdateImageSortOrder 更新图片在分类中的排序
func (cd *ImageCategoryDao) UpdateImageSortOrder(categoryID string, imageID primitive.ObjectID, newSortOrder int) error {
	update := bson.M{
		"$set": bson.M{
			"images.$.sortOrder": newSortOrder,
			"updatedAt":          time.Now(),
		},
	}

	result, err := cd.Collection().UpdateOne(
		context.TODO(),
		bson.M{"_id": categoryID, "images.imageId": imageID},
		update,
	)
	if err != nil {
		global.Logger.Errorf("❌ 更新图片排序失败: %v", err)
		return err
	}

	if result.ModifiedCount == 0 {
		return fmt.Errorf("图片不在该分类中")
	}

	return nil
}

// IsImageInCategory 检查图片是否在分类中
func (cd *ImageCategoryDao) IsImageInCategory(categoryID string, imageID primitive.ObjectID) (bool, error) {
	count, err := cd.Collection().CountDocuments(context.TODO(),
		bson.M{
			"_id": categoryID,
			"images.imageId": imageID,
		})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// getMaxSortOrder 获取分类中的最大排序序号
func (cd *ImageCategoryDao) getMaxSortOrder(categoryID string) (int, error) {
	pipeline := []bson.M{
		{"$match": bson.M{"_id": categoryID}},
		{"$unwind": "$images"},
		{"$group": bson.M{
			"_id":       nil,
			"maxOrder": bson.M{"$max": "$images.sortOrder"},
		}},
	}

	cursor, err := cd.Collection().Aggregate(context.TODO(), pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(context.TODO())

	var result struct {
		MaxOrder int `bson:"maxOrder"`
	}
	if cursor.Next(context.TODO()) {
		if err := cursor.Decode(&result); err != nil {
			return 0, err
		}
		return result.MaxOrder, nil
	}

	return 0, nil // 如果没有图片，返回0
}

// EnsureDefaultCategories 确保默认分类存在
func (cd *ImageCategoryDao) EnsureDefaultCategories() error {
	defaultCategories := []po.ImageCategory{
		{
			ID:          "nexus",
			Name:        "Nexus",
			Description: "默认分类，所有图片初始所属分类",
			ImageCount:  0,
			Settings: po.CategorySettings{
				LayoutMode:    "freeform",
				GridSize:      10,
				DefaultWidth:  200,
				DefaultHeight: 280,
				AutoArrange:   false,
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	for _, category := range defaultCategories {
		_, err := cd.GetCategoryByID(category.ID)
		if err != nil {
			// 分类不存在，创建它
			if err := cd.CreateCategory(&category); err != nil {
				global.Logger.Errorf("❌ 创建默认图片分类失败: %v", err)
				return err
			}
			global.Logger.Infof("✅ 创建默认图片分类: %s", category.Name)
		}
	}

	return nil
}