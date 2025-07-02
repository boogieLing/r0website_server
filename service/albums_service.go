package service

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"r0Website-server/dao"
	"r0Website-server/models/po"
	"r0Website-server/models/vo"
	"time"
)

// AlbumService 画廊服务
type AlbumService struct {
	AlbumDao *dao.AlbumDao `R0Ioc:"true"`
	ImageDao *dao.ImageDao `R0Ioc:"true"`
}

// CreateNewAlbum 创建新图集
func (as *AlbumService) CreateNewAlbum(title, desc string, author string) (*primitive.ObjectID, error) {
	album := &po.Album{
		Title:       title,
		Description: desc,
		Tags:        []string{},
		Visibility:  "private",
		Author:      author,
		CoverImage:  primitive.NilObjectID,
		ImageRefs:   []*po.AlbumImageRef{},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	res, err := as.AlbumDao.CreateAlbum(album)
	if err != nil {
		return nil, err
	}
	id := res.InsertedID.(primitive.ObjectID)
	return &id, nil
}

// GetAlbumDetail 获取图集详情
func (as *AlbumService) GetAlbumDetail(albumID primitive.ObjectID) (*vo.AlbumDetailVo, error) {
	album, err := as.AlbumDao.GetAlbumByID(albumID)
	if err != nil {
		return nil, err
	}
	return &vo.AlbumDetailVo{
		ID:          album.ID,
		Title:       album.Title,
		Description: album.Description,
		CoverImage:  album.CoverImage,
		CreatedAt:   album.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   album.UpdatedAt.Format(time.RFC3339),
		Tags:        album.Tags,
		Visibility:  album.Visibility,
		ImageRefs:   album.ImageRefs,
	}, nil
}

// UpdateImageLayout 更新某张图在图集中的布局
func (as *AlbumService) UpdateImageLayout(albumID, imageID primitive.ObjectID, pos *po.AlbumPosition, caption, desc string) error {
	return as.AlbumDao.UpdateImageLayoutInAlbum(albumID, imageID, pos, caption, desc)
}

// MoveImageBetweenAlbums 图像移动到另一个图集
func (as *AlbumService) MoveImageBetweenAlbums(fromAlbumID, toAlbumID, imageID primitive.ObjectID) error {
	return as.AlbumDao.MoveImageToAnotherAlbum(fromAlbumID, toAlbumID, imageID)
}
