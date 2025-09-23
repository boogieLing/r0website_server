package dao

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"r0Website-server/global"
	"r0Website-server/models/po"
	"time"
)

type ImageDao struct {
	*BasicDaoMongo `R0Ioc:"true"`
}

func (*ImageDao) CollectionName() string {
	return "images"
}

func (id *ImageDao) Collection() *mongo.Collection {
	return id.Mdb.Collection(id.CollectionName())
}

// UploadImage 上传图片
func (id *ImageDao) UploadImage(img *po.Image) (*mongo.InsertOneResult, error) {
	img.UploadedAt = time.Now()
	res, err := id.Collection().InsertOne(context.TODO(), img)
	if err != nil {
		global.Logger.Errorf("❌ 上传图片失败: %v", err)
	}
	return res, err
}

// GetImageByID 获取图片
func (id *ImageDao) GetImageByID(imageID primitive.ObjectID) (*po.Image, error) {
	var img po.Image
	err := id.Collection().FindOne(context.TODO(), bson.M{"_id": imageID}).Decode(&img)
	if err != nil {
		global.Logger.Errorf("❌ 获取图片失败: %v", err)
		return nil, err
	}
	return &img, nil
}

// DeleteImageByID 删除图片
func (id *ImageDao) DeleteImageByID(imageID primitive.ObjectID) error {
	_, err := id.Collection().DeleteOne(context.TODO(), bson.M{"_id": imageID})
	if err != nil {
		global.Logger.Errorf("❌ 删除图片失败: %v", err)
	}
	return err
}

// GetAllAlbumsOfImage 查询图片在哪些图集中被引用
func (id *ImageDao) GetAllAlbumsOfImage(imageID primitive.ObjectID) ([]primitive.ObjectID, error) {
	cursor, err := id.Mdb.Collection("albums").Find(context.TODO(), bson.M{"image_refs.image_id": imageID})
	if err != nil {
		global.Logger.Errorf("❌ 查询图像所在图集失败: %v", err)
		return nil, err
	}
	defer cursor.Close(context.TODO())
	var result []primitive.ObjectID
	for cursor.Next(context.TODO()) {
		var album po.Album
		if err := cursor.Decode(&album); err == nil {
			result = append(result, album.ID)
		}
	}
	return result, nil
}

// ListImages 获取所有图片
func (id *ImageDao) ListImages() ([]*po.Image, error) {
	cursor, err := id.Collection().Find(context.TODO(), bson.M{})
	if err != nil {
		global.Logger.Errorf("❌ 获取图片列表失败: %v", err)
		return nil, err
	}
	var imgs []*po.Image
	if err = cursor.All(context.TODO(), &imgs); err != nil {
		global.Logger.Errorf("❌ 解析图片列表失败: %v", err)
	}
	return imgs, err
}

// FindImagesByTag 按标签获取图片
func (id *ImageDao) FindImagesByTag(tag string) ([]*po.Image, error) {
	cursor, err := id.Collection().Find(context.TODO(), bson.M{"tags": tag})
	if err != nil {
		global.Logger.Errorf("❌ 获取标签图片失败: %v", err)
		return nil, err
	}
	var imgs []*po.Image
	if err = cursor.All(context.TODO(), &imgs); err != nil {
		global.Logger.Errorf("❌ 解析标签图片失败: %v", err)
	}
	return imgs, err
}

// SearchImagesByName 模糊搜索图片名
func (id *ImageDao) SearchImagesByName(keyword string) ([]*po.Image, error) {
	cursor, err := id.Collection().Find(context.TODO(), bson.M{
		"name": bson.M{"$regex": keyword, "$options": "i"},
	})
	if err != nil {
		global.Logger.Errorf("❌ 图片搜索失败: %v", err)
		return nil, err
	}
	var imgs []*po.Image
	if err = cursor.All(context.TODO(), &imgs); err != nil {
		global.Logger.Errorf("❌ 解析搜索结果失败: %v", err)
	}
	return imgs, err
}

// UpdateImageCategoryPosition 更新图片在分类中的位置信息
func (id *ImageDao) UpdateImageCategoryPosition(imageID primitive.ObjectID, categoryID string, position *po.CategoryPosition) error {
	update := bson.M{
		"$set": bson.M{
			"positions." + categoryID: position,
		},
	}
	_, err := id.Collection().UpdateOne(context.TODO(), bson.M{"_id": imageID}, update)
	if err != nil {
		global.Logger.Errorf("❌ 更新图片分类位置失败: %v", err)
	}
	return err
}

// RemoveImageFromCategory 从分类中移除图片
func (id *ImageDao) RemoveImageFromCategory(imageID primitive.ObjectID, categoryID string) error {
	update := bson.M{
		"$unset": bson.M{
			"positions." + categoryID: "",
		},
	}
	_, err := id.Collection().UpdateOne(context.TODO(), bson.M{"_id": imageID}, update)
	if err != nil {
		global.Logger.Errorf("❌ 从分类中移除图片失败: %v", err)
	}
	return err
}

// GetImagesByCategory 获取分类下的所有图片
func (id *ImageDao) GetImagesByCategory(categoryID string) ([]*po.Image, error) {
	// 使用MongoDB的聚合查询，筛选出包含指定分类的图片
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"positions." + categoryID: bson.M{"$exists": true},
			},
		},
		{
			"$addFields": bson.M{
				"sortOrder": "$positions." + categoryID + ".sortOrder",
			},
		},
		{
			"$sort": bson.M{"sortOrder": 1},
		},
	}

	cursor, err := id.Collection().Aggregate(context.TODO(), pipeline)
	if err != nil {
		global.Logger.Errorf("❌ 获取分类图片失败: %v", err)
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var imgs []*po.Image
	if err = cursor.All(context.TODO(), &imgs); err != nil {
		global.Logger.Errorf("❌ 解析分类图片失败: %v", err)
		return nil, err
	}
	return imgs, nil
}

// GetImagesByCategoryWithPagination 分页获取分类下的图片
func (id *ImageDao) GetImagesByCategoryWithPagination(categoryID string, page, pageSize int) ([]*po.Image, int64, error) {
	skip := (page - 1) * pageSize

	// 计算总数
	count, err := id.Collection().CountDocuments(context.TODO(), bson.M{
		"positions." + categoryID: bson.M{"$exists": true},
	})
	if err != nil {
		global.Logger.Errorf("❌ 计算分类图片数量失败: %v", err)
		return nil, 0, err
	}

	// 分页查询
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"positions." + categoryID: bson.M{"$exists": true},
			},
		},
		{
			"$addFields": bson.M{
				"sortOrder": "$positions." + categoryID + ".sortOrder",
			},
		},
		{
			"$sort": bson.M{"sortOrder": 1},
		},
		{
			"$skip": skip,
		},
		{
			"$limit": pageSize,
		},
	}

	cursor, err := id.Collection().Aggregate(context.TODO(), pipeline)
	if err != nil {
		global.Logger.Errorf("❌ 分页获取分类图片失败: %v", err)
		return nil, 0, err
	}
	defer cursor.Close(context.TODO())

	var imgs []*po.Image
	if err = cursor.All(context.TODO(), &imgs); err != nil {
		global.Logger.Errorf("❌ 解析分页分类图片失败: %v", err)
		return nil, 0, err
	}
	return imgs, count, nil
}
