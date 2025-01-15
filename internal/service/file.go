package service

import (
	"context"
	"github.com/huangyul/gin-vue-template/internal/domain"
	"github.com/huangyul/gin-vue-template/internal/dto"
	"github.com/huangyul/gin-vue-template/internal/repository"
	"mime/multipart"
)

var (
	staticPath = "static"
)

type FileService struct {
	repo    *repository.FileRepository
	userSvc UserService
}

func NewFileService(repo *repository.FileRepository, userSvc UserService) *FileService {
	return &FileService{repo: repo, userSvc: userSvc}
}

func (svc *FileService) Save(ctx context.Context, file *multipart.FileHeader, link string, uId int64) error {
	user, err := svc.userSvc.GetByID(ctx, uId)
	if err != nil {
		return err
	}
	return svc.repo.Insert(ctx, file.Filename, user.Username, user.ID, link)
}

func (svc *FileService) Delete(ctx context.Context, id int64, uId int64) (string, error) {
	return svc.repo.Delete(ctx, id, uId)
}

func (svc *FileService) List(ctx context.Context, param dto.FileListQueryParam) ([]domain.File, int64, error) {
	res, total, err := svc.repo.List(ctx, param)
	if err != nil {
		return []domain.File{}, 0, err
	}
	return res, total, nil
}

func (svc *FileService) GetOption(ctx context.Context) ([]dto.QuerySelectOption, error) {
	return svc.userSvc.GetOptions(ctx)
}
