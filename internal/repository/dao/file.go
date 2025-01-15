package dao

import (
	"context"
	"github.com/huangyul/gin-vue-template/internal/dto"
	"github.com/huangyul/gin-vue-template/internal/pkg/errno"
	"gorm.io/gorm"
	"time"
)

type FileDao struct {
	db *gorm.DB
}

func NewFileDao(db *gorm.DB) *FileDao {
	return &FileDao{
		db: db,
	}
}

func (dao *FileDao) Insert(ctx context.Context, fileName string, userName string, userId int64, link string) error {
	now := time.Now()
	file := File{
		Name:       fileName,
		Uploader:   userName,
		UploaderId: userId,
		Link:       link,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	return dao.db.WithContext(ctx).Create(&file).Error
}

func (dao *FileDao) Delete(ctx context.Context, id int64, uId int64) (string, error) {
	var file File
	res := dao.db.WithContext(ctx).Model(&File{}).First(&file).Delete("id = ? AND uploader_id = ?", id, uId).RowsAffected
	if res == 0 {
		return "", errno.FileNotPermisson
	}
	return file.Link, nil
}

func (dao *FileDao) List(ctx context.Context, param dto.FileListQueryParam) ([]File, int64, error) {
	page := param.Page
	pageSize := param.PageSize
	if pageSize == 0 {
		pageSize = 10
	}
	if page == 0 {
		page = 1
	}
	query := dao.db.WithContext(ctx).Model(&File{})
	if param.FileName != "" {
		query = query.Where("name LIKE ?", "%"+param.FileName+"%")
	}
	if param.UserId != "" {
		query = query.Where("uploader_id = ?", param.UserId)
	}
	var total int64
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	var list []File
	err = query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

type File struct {
	Id            int64  `gorm:"primary_key;AUTO_INCREMENT"`
	Name          string `gorm:"type:varchar(255);not null;index:idx_file_name"`
	Uploader      string `gorm:"type:varchar(255);"`
	UploaderId    int64
	Link          string
	DownloadCount int
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
