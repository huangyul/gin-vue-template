package service

import (
	"context"
	"github.com/huangyul/gin-vue-template/internal/domain"
	"github.com/huangyul/gin-vue-template/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(ctx context.Context, username string, password string) error
	Login(ctx context.Context, username string, password string) (domain.User, error)
}

type UserServiceImpl struct {
	repo repository.UserRepository
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
