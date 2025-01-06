package repository

import "github.com/huangyul/gin-vue-template/internal/repository/dao"

type UserRepository interface {
}

type UserRepositoryImpl struct {
	dao dao.UserDao
}

func NewUserRepository(dao dao.UserDao) UserRepository {
	return &UserRepositoryImpl{
		dao: dao,
	}
}
