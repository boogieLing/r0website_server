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
	AlbumDao         *dao.AlbumDao         `R0Ioc:"true"`
	ImageDao         *dao.ImageDao         `R0Ioc:"true"`
	ImageCategoryDao *dao.ImageCategoryDao `R0Ioc:"true"`
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

// AddImageToCategory 添加图片到分类
func (as *AlbumService) AddImageToCategory(imageID primitive.ObjectID, categoryID string, categoryName string, position *po.CategoryPosition) error {
	// 设置默认值
	if position.LayoutMode == "" {
		position.LayoutMode = "freeform"
	}
	if position.GridSize == 0 {
		position.GridSize = 10
	}
	if position.Version == 0 {
		position.Version = 1
	}
	if position.AddedAt.IsZero() {
		position.AddedAt = time.Now()
	}
	if position.UpdatedAt.IsZero() {
		position.UpdatedAt = time.Now()
	}

	// 设置分类信息
	position.CategoryName = categoryName

	return as.ImageDao.UpdateImageCategoryPosition(imageID, categoryID, position)
}

// UpdateImageInCategory 更新图片在分类中的信息
func (as *AlbumService) UpdateImageInCategory(imageID primitive.ObjectID, categoryID string, position *po.CategoryPosition) error {
	// 保持原有分类名称和添加时间
	image, err := as.ImageDao.GetImageByID(imageID)
	if err != nil {
		return err
	}

	if existingPos, exists := image.Positions[categoryID]; exists {
		position.CategoryName = existingPos.CategoryName
		position.AddedAt = existingPos.AddedAt
		position.SortOrder = existingPos.SortOrder
	}

	position.UpdatedAt = time.Now()
	position.Version++

	return as.ImageDao.UpdateImageCategoryPosition(imageID, categoryID, position)
}

// RemoveImageFromCategory 从分类中移除图片
func (as *AlbumService) RemoveImageFromCategory(imageID primitive.ObjectID, categoryID string) error {
	return as.ImageDao.RemoveImageFromCategory(imageID, categoryID)
}

// GetImagesByCategory 获取分类下的图片
func (as *AlbumService) GetImagesByCategory(categoryID string) ([]*po.Image, error) {
	return as.ImageDao.GetImagesByCategory(categoryID)
}

// GetImagesByCategoryWithPagination 分页获取分类下的图片
func (as *AlbumService) GetImagesByCategoryWithPagination(categoryID string, page, pageSize int) ([]*po.Image, int64, error) {
	return as.ImageDao.GetImagesByCategoryWithPagination(categoryID, page, pageSize)
}

// BatchUpdateCategoryPositions 批量更新分类中图片的位置
func (as *AlbumService) BatchUpdateCategoryPositions(categoryID string, updates map[primitive.ObjectID]*po.CategoryPosition) error {
	for imageID, position := range updates {
		if err := as.UpdateImageInCategory(imageID, categoryID, position); err != nil {
			return err
		}
	}
	return nil
}

// InitDefaultCategories 初始化默认分类
func (as *AlbumService) InitDefaultCategories() error {
	return as.ImageCategoryDao.EnsureDefaultCategories()
}
