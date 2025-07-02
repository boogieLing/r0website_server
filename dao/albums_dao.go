package dao

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"r0Website-server/global"
	"r0Website-server/models/po"
	"time"
)

type AlbumDao struct {
	*BasicDaoMongo `R0Ioc:"true"`
}

func (*AlbumDao) CollectionName() string {
	return "albums"
}

func (ad *AlbumDao) Collection() *mongo.Collection {
	return ad.Mdb.Collection(ad.CollectionName())
}

// CreateAlbum 创建图集
func (ad *AlbumDao) CreateAlbum(album *po.Album) (*mongo.InsertOneResult, error) {
	album.CreatedAt = time.Now()
	album.UpdatedAt = time.Now()
	res, err := ad.Collection().InsertOne(context.TODO(), album)
	if err != nil {
		global.Logger.Errorf("❌ 创建图集失败: %v", err)
	}
	return res, err
}

// GetAlbumByID 根据 ID 查询图集
func (ad *AlbumDao) GetAlbumByID(id primitive.ObjectID) (*po.Album, error) {
	var album po.Album
	err := ad.Collection().FindOne(context.TODO(), bson.M{"_id": id}).Decode(&album)
	if err != nil {
		global.Logger.Errorf("❌ 查询图集失败 id=%s: %v", id.Hex(), err)
		return nil, err
	}
	return &album, nil
}

// UpdateAlbum 更新图集元信息
func (ad *AlbumDao) UpdateAlbum(album *po.Album) error {
	album.UpdatedAt = time.Now()
	_, err := ad.Collection().UpdateByID(context.TODO(), album.ID, bson.M{
		"$set": bson.M{
			"title":       album.Title,
			"description": album.Description,
			"tags":        album.Tags,
			"visibility":  album.Visibility,
			"cover_image": album.CoverImage,
			"updated_at":  album.UpdatedAt,
		},
	})
	if err != nil {
		global.Logger.Errorf("❌ 更新图集失败: %v", err)
	}
	return err
}

// DeleteAlbum 删除图集
func (ad *AlbumDao) DeleteAlbum(id primitive.ObjectID) error {
	_, err := ad.Collection().DeleteOne(context.TODO(), bson.M{"_id": id})
	if err != nil {
		global.Logger.Errorf("❌ 删除图集失败: %v", err)
	}
	return err
}

// AddOrUpdateImageRef 添加或更新图集中某张图的布局/描述
func (ad *AlbumDao) AddOrUpdateImageRef(albumID, imageID primitive.ObjectID, ref *po.AlbumImageRef) error {
	ctx := context.TODO()
	filter := bson.M{"_id": albumID, "image_refs.image_id": imageID}

	update := bson.M{
		"$set": bson.M{
			"image_refs.$.position":    ref.Position,
			"image_refs.$.caption":     ref.Caption,
			"image_refs.$.description": ref.Description,
			"updated_at":               time.Now(),
		},
	}
	result, err := ad.Collection().UpdateOne(ctx, filter, update)
	if err != nil {
		global.Logger.Errorf("❌ 更新图像引用失败: %v", err)
		return err
	}
	if result.MatchedCount == 0 {
		push := bson.M{
			"$push": bson.M{"image_refs": ref},
			"$set":  bson.M{"updated_at": time.Now()},
		}
		_, err = ad.Collection().UpdateByID(ctx, albumID, push)
		if err != nil {
			global.Logger.Errorf("❌ 添加图像引用失败: %v", err)
			return err
		}
	}
	return nil
}

// RemoveImageFromAlbum 从图集中移除图像
func (ad *AlbumDao) RemoveImageFromAlbum(albumID, imageID primitive.ObjectID) error {
	update := bson.M{
		"$pull": bson.M{"image_refs": bson.M{"image_id": imageID}},
		"$set":  bson.M{"updated_at": time.Now()},
	}
	_, err := ad.Collection().UpdateByID(context.TODO(), albumID, update)
	if err != nil {
		global.Logger.Errorf("❌ 移除图像失败: %v", err)
	}
	return err
}

// MoveImageToAnotherAlbum 图像在图集间移动
func (ad *AlbumDao) MoveImageToAnotherAlbum(fromAlbumID, toAlbumID, imageID primitive.ObjectID) error {
	album, err := ad.GetAlbumByID(fromAlbumID)
	if err != nil {
		return err
	}
	var ref *po.AlbumImageRef
	for _, r := range album.ImageRefs {
		if r.ImageID == imageID {
			ref = r
			break
		}
	}
	if ref == nil {
		return fmt.Errorf("image not found in source album")
	}
	if err := ad.RemoveImageFromAlbum(fromAlbumID, imageID); err != nil {
		return err
	}
	return ad.AddOrUpdateImageRef(toAlbumID, imageID, ref)
}

// ListAlbums 获取所有图集
func (ad *AlbumDao) ListAlbums() ([]*po.Album, error) {
	cursor, err := ad.Collection().Find(context.TODO(), bson.M{})
	if err != nil {
		global.Logger.Errorf("❌ 获取图集列表失败: %v", err)
		return nil, err
	}
	var albums []*po.Album
	if err = cursor.All(context.TODO(), &albums); err != nil {
		global.Logger.Errorf("❌ 解析图集列表失败: %v", err)
	}
	return albums, err
}

// FindAlbumsByAuthor 根据作者查图集
func (ad *AlbumDao) FindAlbumsByAuthor(author string) ([]*po.Album, error) {
	cursor, err := ad.Collection().Find(context.TODO(), bson.M{"author": author})
	if err != nil {
		global.Logger.Errorf("❌ 查询作者图集失败: %v", err)
		return nil, err
	}
	var albums []*po.Album
	if err = cursor.All(context.TODO(), &albums); err != nil {
		global.Logger.Errorf("❌ 解析作者图集失败: %v", err)
	}
	return albums, err
}

// FindAlbumsByTag 根据标签查图集
func (ad *AlbumDao) FindAlbumsByTag(tag string) ([]*po.Album, error) {
	cursor, err := ad.Collection().Find(context.TODO(), bson.M{"tags": tag})
	if err != nil {
		global.Logger.Errorf("❌ 查询标签图集失败: %v", err)
		return nil, err
	}
	var albums []*po.Album
	if err = cursor.All(context.TODO(), &albums); err != nil {
		global.Logger.Errorf("❌ 解析标签图集失败: %v", err)
	}
	return albums, err
}

// SearchAlbumsByKeyword 关键字搜索图集
func (ad *AlbumDao) SearchAlbumsByKeyword(keyword string) ([]*po.Album, error) {
	filter := bson.M{"$text": bson.M{"$search": keyword}}
	cursor, err := ad.Collection().Find(context.TODO(), filter)
	if err != nil {
		global.Logger.Errorf("❌ 搜索图集失败: %v", err)
		return nil, err
	}
	var albums []*po.Album
	if err = cursor.All(context.TODO(), &albums); err != nil {
		global.Logger.Errorf("❌ 解析搜索结果失败: %v", err)
	}
	return albums, err
}

// UpdateImageLayoutInAlbum 更新图集中某张图片的布局信息（位置、大小、描述等）
func (ad *AlbumDao) UpdateImageLayoutInAlbum(
	albumID, imageID primitive.ObjectID, layout *po.AlbumPosition, caption, description string,
) error {
	ctx := context.TODO()

	filter := bson.M{
		"_id":                 albumID,
		"image_refs.image_id": imageID,
	}

	update := bson.M{
		"$set": bson.M{
			"image_refs.$.position":    layout,
			"image_refs.$.caption":     caption,
			"image_refs.$.description": description,
			"updated_at":               time.Now(),
		},
	}

	result, err := ad.Collection().UpdateOne(ctx, filter, update)
	if err != nil {
		global.Logger.Errorf("❌ 更新图集中图片布局失败 album_id=%s image_id=%s: %v", albumID.Hex(), imageID.Hex(), err)
		return err
	}
	if result.MatchedCount == 0 {
		global.Logger.Warnf("⚠️ 图集中未找到 image_id=%s 无法更新布局", imageID.Hex())
		return mongo.ErrNoDocuments
	}

	global.Logger.Infof("✅ 更新图像布局成功 album_id=%s image_id=%s", albumID.Hex(), imageID.Hex())
	return nil
}
