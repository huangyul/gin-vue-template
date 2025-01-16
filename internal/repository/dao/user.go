package dao

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"github.com/huangyul/gin-vue-template/internal/pkg/errno"
	"gorm.io/gorm"
	"time"
)

type UserDao interface {
	InsertByUsername(ctx context.Context, user User) error
	FindByUsername(ctx context.Context, username string) (User, error)
	GetList(ctx context.Context, params UserListQueryParam) ([]User, int64, error)
	UpdateByID(ctx context.Context, id int64, nickname string) error
	FindUserById(ctx context.Context, id int64) (User, error)
	DeleteByID(ctx context.Context, id int64) error
	GetAllUser(ctx context.Context) ([]User, error)
}

type UserDaoImpl struct {
	db *gorm.DB
}

func (u *UserDaoImpl) GetAllUser(ctx context.Context) ([]User, error) {
	result := []User{}
	err := u.db.WithContext(ctx).Find(&result).Error
	return result, err
}

func (u *UserDaoImpl) DeleteByID(ctx context.Context, id int64) error {
	return u.db.WithContext(ctx).Delete(&User{}, "id=?", id).Error
}

func (u *UserDaoImpl) FindUserById(ctx context.Context, id int64) (User, error) {
	var user User
	err := u.db.WithContext(ctx).First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return User{}, errno.UserNotFound
		}
		return User{}, err
	}
	return user, nil
}

func (u *UserDaoImpl) UpdateByID(ctx context.Context, id int64, nickname string) error {
	now := time.Now()
	err := u.db.WithContext(ctx).Model(&User{}).Where("id = ?", id).Updates(map[string]interface{}{
		"nickname":   nickname,
		"updated_at": now,
	}).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errno.UserNotFound
		}
		return err
	}
	return nil
}

type UserListQueryParam struct {
	Nickname string
	UserName string
	Page     int
	PageSize int
}

func (u *UserDaoImpl) GetList(ctx context.Context, params UserListQueryParam) ([]User, int64, error) {
	query := u.db.WithContext(ctx).Model(&User{})

	if params.Nickname != "" {
		query = query.Where("nickname = ?", params.Nickname)
	}
	if params.UserName != "" {
		query = query.Where("username = ?", params.UserName)
	}
	if params.Page == 0 {
		params.Page = 1
	}
	if params.PageSize == 0 {
		params.PageSize = 10
	}
	var count int64
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	var users []User
	if err := query.Offset((params.Page - 1) * params.PageSize).Limit(params.PageSize).Order("id asc").Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, count, nil
}

func (u *UserDaoImpl) FindByUsername(ctx context.Context, username string) (User, error) {
	var user User
	err := u.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return User{}, errno.UserNotFound
	}
	return user, err
}

func (u *UserDaoImpl) InsertByUsername(ctx context.Context, user User) error {
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now
	err := u.db.WithContext(ctx).Create(&user).Error
	if err != nil {
		var er *mysql.MySQLError
		errors.As(err, &er)
		if er.Number == 1062 {
			return errno.UsernameConflict
		}
		return err
	}
	return nil
}

func NewUserDao(db *gorm.DB) UserDao {
	return &UserDaoImpl{
		db: db,
	}
}

type User struct {
	Id        int64 `gorm:"primary_key;AUTO_INCREMENT"`
	Avatar    string
	Username  string `gorm:"unique;type:varchar(255)"`
	Nickname  string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
