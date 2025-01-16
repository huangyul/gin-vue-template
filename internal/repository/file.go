package repository

import (
	"context"
	"github.com/huangyul/gin-vue-template/internal/domain"
	"github.com/huangyul/gin-vue-template/internal/dto"
	"github.com/huangyul/gin-vue-template/internal/repository/dao"
)

type FileRepository struct {
	dao *dao.FileDao
}

func NewFileRepository(dao *dao.FileDao) *FileRepository {
	return &FileRepository{dao: dao}
}

func (repo *FileRepository) Insert(ctx context.Context, fileName string, userName string, userId int64, link string) error {
	return repo.dao.Insert(ctx, fileName, userName, userId, link)
}

func (repo *FileRepository) Delete(ctx context.Context, id int64, uId int64) (domain.File, error) {
	f, err := repo.dao.Delete(ctx, id, uId)
	if err != nil {
		return domain.File{}, err
	}
	return domain.File{
		Id:         f.Id,
		FileName:   f.Name,
		Link:       f.Link,
		Uploader:   f.Uploader,
		UploaderId: f.UploaderId,
		CreateAt:   f.CreatedAt,
	}, nil
}

func (repo *FileRepository) List(ctx context.Context, param dto.FileListQueryParam) ([]domain.File, int64, error) {
	res, total, err := repo.dao.List(ctx, param)
	files := make([]domain.File, 0)
	if err != nil {
		return files, 0, err
	}
	for _, file := range res {
		files = append(files, domain.File{
			Id:         file.Id,
			FileName:   file.Name,
			Uploader:   file.Uploader,
			UploaderId: file.UploaderId,
			Link:       file.Link,
			CreateAt:   file.CreatedAt,
		})
	}
	return files, total, nil
}
