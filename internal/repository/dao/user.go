package dao

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrUsernameConflict = errors.New("username is already taken")
	ErrUserNotFound     = errors.New("user not found")
)

type UserDao interface {
	InsertByUsername(ctx context.Context, user User) error
	FindByUsername(ctx context.Context, username string) (User, error)
}

type UserDaoImpl struct {
	db *gorm.DB
}

func (u *UserDaoImpl) FindByUsername(ctx context.Context, username string) (User, error) {
	var user User
	err := u.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return User{}, ErrUserNotFound
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
			return ErrUsernameConflict
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
