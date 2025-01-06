package service

import "github.com/huangyul/gin-vue-template/internal/repository"

type UserService interface{}

type UserServiceImpl struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &UserServiceImpl{
		repo: repo,
	}
}
