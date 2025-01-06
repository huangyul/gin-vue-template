package repository

import (
	"context"
	"github.com/huangyul/gin-vue-template/internal/domain"
	"github.com/huangyul/gin-vue-template/internal/repository/dao"
)

type UserRepository interface {
	CreateByUsername(ctx context.Context, user domain.User) error
	FindByUsername(ctx context.Context, username string) (domain.User, error)
}

type UserRepositoryImpl struct {
	dao dao.UserDao
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
