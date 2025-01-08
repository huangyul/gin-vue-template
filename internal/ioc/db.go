package ioc

import (
	"fmt"
	"github.com/huangyul/gin-vue-template/internal/repository/dao"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

func InitDB() *gorm.DB {
	type DBConfig struct {
		Host     string
		Port     int
		Username string
		Password string
		Database string
	}
	var cfg DBConfig
	if err := viper.UnmarshalKey("db", &cfg); err != nil {
		panic(err)
	}
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		})
	db, err := gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}
	dao.InitTable(db)
	return db
}
