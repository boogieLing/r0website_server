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
	"strings"
	"time"
)

type TagDao struct {
	*BasicDaoMongo `R0Ioc:"true"`
}

func (*TagDao) CollectionName() string {
	return "tags"
}

func (td *TagDao) Collection() *mongo.Collection {
	return td.Mdb.Collection(td.CollectionName())
}

// CreateTag 创建标签
func (td *TagDao) CreateTag(tag *po.Tag) error {
	tag.CreatedAt = time.Now()
	tag.UpdatedAt = time.Now()
	tag.ImageCount = 0

	// 标准化标签ID
	tag.ID = strings.ToLower(strings.TrimSpace(tag.Name))
	if tag.DisplayName == "" {
		tag.DisplayName = tag.Name
	}

	_, err := td.Collection().InsertOne(context.TODO(), tag)
	if err != nil {
		global.Logger.Errorf("❌ 创建标签失败: %v", err)
		return err
	}
	return nil
}

// GetTagByID 根据ID获取标签
func (td *TagDao) GetTagByID(tagID string) (*po.Tag, error) {
	var tag po.Tag
	err := td.Collection().FindOne(context.TODO(), bson.M{"_id": tagID}).Decode(&tag)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("标签不存在")
		}
		global.Logger.Errorf("❌ 获取标签失败: %v", err)
		return nil, err
	}
	return &tag, nil
}

// GetTagByName 根据名称获取标签
func (td *TagDao) GetTagByName(name string) (*po.Tag, error) {
	tagID := strings.ToLower(strings.TrimSpace(name))
	return td.GetTagByID(tagID)
}

// UpdateTag 更新标签信息
func (td *TagDao) UpdateTag(tagID string, update bson.M) error {
	if _, ok := update["$set"]; !ok {
		update["$set"] = bson.M{}
	}
	update["$set"].(bson.M)["updatedAt"] = time.Now()

	_, err := td.Collection().UpdateOne(context.TODO(), bson.M{"_id": tagID}, update)
	if err != nil {
		global.Logger.Errorf("❌ 更新标签失败: %v", err)
		return err
	}
	return nil
}

// DeleteTag 删除标签
func (td *TagDao) DeleteTag(tagID string) error {
	// 检查标签中是否有图片
	count, err := td.GetTagImageCount(tagID)
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("标签中还存在图片，无法删除")
	}

	_, err = td.Collection().DeleteOne(context.TODO(), bson.M{"_id": tagID})
	if err != nil {
		global.Logger.Errorf("❌ 删除标签失败: %v", err)
		return err
	}
	return nil
}

// ListTags 获取所有标签
func (td *TagDao) ListTags(category string) ([]*po.Tag, error) {
	filter := bson.M{}
	if category != "" {
		filter["category"] = category
	}

	cursor, err := td.Collection().Find(context.TODO(), filter, options.Find().SetSort(bson.M{"name": 1}))
	if err != nil {
		global.Logger.Errorf("❌ 获取标签列表失败: %v", err)
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var tags []*po.Tag
	if err = cursor.All(context.TODO(), &tags); err != nil {
		global.Logger.Errorf("❌ 解析标签列表失败: %v", err)
		return nil, err
	}
	return tags, nil
}

// AddImageToTag 添加图片到标签（维护倒排索引）
func (td *TagDao) AddImageToTag(tagID string, imageID primitive.ObjectID, imageName string) error {
	// 检查图片是否已存在
	exists, err := td.IsImageInTag(tagID, imageID)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("图片已存在于该标签中")
	}

	// 添加到标签的图片列表中
	imageRef := po.TagImageRef{
		ImageID:   imageID,
		ImageName: imageName,
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

	_, err = td.Collection().UpdateOne(context.TODO(), bson.M{"_id": tagID}, update)
	if err != nil {
		global.Logger.Errorf("❌ 添加图片到标签失败: %v", err)
		return err
	}
	return nil
}

// RemoveImageFromTag 从标签中移除图片
func (td *TagDao) RemoveImageFromTag(tagID string, imageID primitive.ObjectID) error {
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

	result, err := td.Collection().UpdateOne(context.TODO(), bson.M{"_id": tagID}, update)
	if err != nil {
		global.Logger.Errorf("❌ 从标签中移除图片失败: %v", err)
		return err
	}

	if result.ModifiedCount == 0 {
		return fmt.Errorf("图片不在该标签中")
	}

	return nil
}

// GetTagImages 获取标签中的图片列表
func (td *TagDao) GetTagImages(tagID string, page, pageSize int) ([]po.TagImageRef, int64, error) {
	// 获取标签信息
	tag, err := td.GetTagByID(tagID)
	if err != nil {
		return nil, 0, err
	}
	total := int64(tag.ImageCount)

	// 如果分页参数无效，返回所有图片引用
	if page <= 0 || pageSize <= 0 {
		return tag.Images, total, nil
	}

	// 分页查询
	skip := (page - 1) * pageSize
	pipeline := []bson.M{
		{"$match": bson.M{"_id": tagID}},
		{"$unwind": "$images"},
		{"$skip": skip},
		{"$limit": pageSize},
		{"$replaceRoot": bson.M{"newRoot": "$images"}},
	}

	cursor, err := td.Collection().Aggregate(context.TODO(), pipeline)
	if err != nil {
		global.Logger.Errorf("❌ 获取标签图片失败: %v", err)
		return nil, 0, err
	}
	defer cursor.Close(context.TODO())

	var imageRefs []po.TagImageRef
	if err = cursor.All(context.TODO(), &imageRefs); err != nil {
		global.Logger.Errorf("❌ 解析标签图片失败: %v", err)
		return nil, 0, err
	}

	return imageRefs, total, nil
}

// GetTagImageCount 获取标签中的图片数量
func (td *TagDao) GetTagImageCount(tagID string) (int, error) {
	tag, err := td.GetTagByID(tagID)
	if err != nil {
		return 0, err
	}
	return tag.ImageCount, nil
}

// IsImageInTag 检查图片是否在标签中
func (td *TagDao) IsImageInTag(tagID string, imageID primitive.ObjectID) (bool, error) {
	count, err := td.Collection().CountDocuments(context.TODO(),
		bson.M{
			"_id": tagID,
			"images.imageId": imageID,
		})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetPopularTags 获取热门标签
func (td *TagDao) GetPopularTags(limit int, category string) ([]*po.TagStats, error) {
	if limit <= 0 {
		limit = 20
	}

	pipeline := []bson.M{
		{"$match": bson.M{}},
	}

	if category != "" {
		pipeline[0] = bson.M{"$match": bson.M{"category": category}}
	}

	pipeline = append(pipeline, []bson.M{
		{"$project": bson.M{
			"tagId":       "$_id",
			"name":        "$name",
			"displayName": "$displayName",
			"imageCount":  "$imageCount",
		}},
		{"$sort": bson.M{"imageCount": -1}},
		{"$limit": limit},
	}...)

	cursor, err := td.Collection().Aggregate(context.TODO(), pipeline)
	if err != nil {
		global.Logger.Errorf("❌ 获取热门标签失败: %v", err)
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var stats []*po.TagStats
	if err = cursor.All(context.TODO(), &stats); err != nil {
		global.Logger.Errorf("❌ 解析热门标签失败: %v", err)
		return nil, err
	}
	return stats, nil
}

// GetOrCreateTag 获取或创建标签
func (td *TagDao) GetOrCreateTag(name string, displayName string, category string) (*po.Tag, error) {
	tagID := strings.ToLower(strings.TrimSpace(name))

	// 尝试获取现有标签
	tag, err := td.GetTagByID(tagID)
	if err == nil {
		return tag, nil
	}

	// 标签不存在，创建新标签
	newTag := &po.Tag{
		ID:          tagID,
		Name:        name,
		DisplayName: displayName,
		Category:    category,
		ImageCount:  0,
		Images:      []po.TagImageRef{},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := td.CreateTag(newTag); err != nil {
		return nil, err
	}

	return newTag, nil
}

// SyncImageTags 同步图片的标签（批量更新）
func (td *TagDao) SyncImageTags(imageID primitive.ObjectID, imageName string, oldTags []string, newTags []string) error {
	// 找出需要移除的标签
	tagsToRemove := []string{}
	for _, oldTag := range oldTags {
		found := false
		for _, newTag := range newTags {
			if strings.EqualFold(oldTag, newTag) {
				found = true
				break
			}
		}
		if !found {
			tagsToRemove = append(tagsToRemove, oldTag)
		}
	}

	// 找出需要添加的标签
	tagsToAdd := []string{}
	for _, newTag := range newTags {
		found := false
		for _, oldTag := range oldTags {
			if strings.EqualFold(oldTag, newTag) {
				found = true
				break
			}
		}
		if !found {
			tagsToAdd = append(tagsToAdd, newTag)
		}
	}

	// 移除旧标签
	for _, tagName := range tagsToRemove {
		if tag, err := td.GetTagByName(tagName); err == nil {
			if err := td.RemoveImageFromTag(tag.ID, imageID); err != nil {
				global.Logger.Errorf("从标签 %s 移除图片失败: %v", tagName, err)
			}
		}
	}

	// 添加新标签
	for _, tagName := range tagsToAdd {
		tag, err := td.GetOrCreateTag(tagName, tagName, "")
		if err != nil {
			global.Logger.Errorf("获取或创建标签 %s 失败: %v", tagName, err)
			continue
		}
		if err := td.AddImageToTag(tag.ID, imageID, imageName); err != nil {
			global.Logger.Errorf("添加图片到标签 %s 失败: %v", tagName, err)
		}
	}

	return nil
}