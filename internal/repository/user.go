package repository

import (
	"context"
	"github.com/huangyul/gin-vue-template/internal/domain"
	"github.com/huangyul/gin-vue-template/internal/dto"
	"github.com/huangyul/gin-vue-template/internal/repository/dao"
)

type UserRepository interface {
	CreateByUsername(ctx context.Context, user domain.User) error
	FindByUsername(ctx context.Context, username string) (domain.User, error)
	GetList(ctx context.Context, params dto.UserListQueryParams) ([]domain.User, int64, error)
	UpdateByID(ctx context.Context, id int64, nickname string) error
	FindByID(ctx context.Context, id int64) (domain.User, error)
	DeleteByID(ctx context.Context, id int64) error
}

type UserRepositoryImpl struct {
	dao dao.UserDao
}

func (u *UserRepositoryImpl) UpdateByID(ctx context.Context, id int64, nickname string) error {
	return u.dao.UpdateByID(ctx, id, nickname)
}

func (u *UserRepositoryImpl) FindByID(ctx context.Context, id int64) (domain.User, error) {
	user, err := u.dao.FindUserById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	return u.toDomain(user), nil
}

func (u *UserRepositoryImpl) DeleteByID(ctx context.Context, id int64) error {
	return u.dao.DeleteByID(ctx, id)
}

func (u *UserRepositoryImpl) GetList(ctx context.Context, params dto.UserListQueryParams) ([]domain.User, int64, error) {
	users, count, err := u.dao.GetList(ctx, dao.UserListQueryParam{
		Nickname: params.Nickname,
		UserName: params.Username,
		Page:     params.Page,
		PageSize: params.PageSize,
	})
	if err != nil {
		return nil, 0, err
	}
	var res []domain.User
	for _, user := range users {
		res = append(res, u.toDomain(user))
	}
	return res, count, nil
}

func (u *UserRepositoryImpl) FindByUsername(ctx context.Context, username string) (domain.User, error) {
	user, err := u.dao.FindByUsername(ctx, username)
	return u.toDomain(user), err
}

func (u *UserRepositoryImpl) CreateByUsername(ctx context.Context, user domain.User) error {
	return u.dao.InsertByUsername(ctx, u.toEntity(user))
}

func NewUserRepository(dao dao.UserDao) UserRepository {
	return &UserRepositoryImpl{
		dao: dao,
	}
}

func (u *UserRepositoryImpl) toEntity(user domain.User) dao.User {
	return dao.User{
		Username: user.Username,
		Password: user.Password,
	}
}

func (u *UserRepositoryImpl) toDomain(user dao.User) domain.User {
	return domain.User{
		ID:       user.Id,
		Username: user.Username,
		Password: user.Password,
	}
}
