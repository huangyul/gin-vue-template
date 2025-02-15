package service

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/huangyul/gin-vue-template/internal/domain"
	"github.com/huangyul/gin-vue-template/internal/dto"
	"github.com/huangyul/gin-vue-template/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"strconv"
)

type UserService interface {
	Register(ctx context.Context, username string, password string) error
	Login(ctx context.Context, username string, password string) (domain.User, error)
	List(ctx *gin.Context, params dto.UserListQueryParams) ([]domain.User, int64, error)
	GetByID(ctx context.Context, id int64) (domain.User, error)
	Update(ctx context.Context, id int64, nickname string) error
	DeleteByID(ctx context.Context, id int64) error
	Create(ctx *gin.Context, username string, password string, nickname string, avatar string) error
	GetOptions(ctx context.Context) ([]dto.QuerySelectOption, error)
}

type UserServiceImpl struct {
	repo repository.UserRepository
}

func (u *UserServiceImpl) GetOptions(ctx context.Context) ([]dto.QuerySelectOption, error) {
	users, err := u.repo.GetAllUser(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]dto.QuerySelectOption, len(users))
	for i, user := range users {
		result[i] = dto.QuerySelectOption{
			Label: user.Username,
			Value: strconv.Itoa(int(user.ID)),
		}
	}
	return result, nil
}

func (u *UserServiceImpl) Create(ctx *gin.Context, username string, password string, nickname string, avatar string) error {
	passStr, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return u.repo.CreateByUsername(ctx, domain.User{Username: username, Password: string(passStr), Nickname: nickname, Avatar: avatar})
}

func (u *UserServiceImpl) GetByID(ctx context.Context, id int64) (domain.User, error) {
	return u.repo.FindByID(ctx, id)
}

func (u *UserServiceImpl) Update(ctx context.Context, id int64, nickname string) error {
	return u.repo.UpdateByID(ctx, id, nickname)
}

func (u *UserServiceImpl) DeleteByID(ctx context.Context, id int64) error {
	return u.repo.DeleteByID(ctx, id)
}

func (u *UserServiceImpl) List(ctx *gin.Context, params dto.UserListQueryParams) ([]domain.User, int64, error) {
	return u.repo.GetList(ctx, params)
}

func (u *UserServiceImpl) Login(ctx context.Context, username string, password string) (domain.User, error) {
	user, err := u.repo.FindByUsername(ctx, username)
	if err != nil {
		return domain.User{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (u *UserServiceImpl) Register(ctx context.Context, username string, password string) error {
	passStr, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return u.repo.CreateByUsername(ctx, domain.User{
		Username: username,
		Password: string(passStr),
	})
}

func NewUserService(repo repository.UserRepository) UserService {
	return &UserServiceImpl{
		repo: repo,
	}
}
