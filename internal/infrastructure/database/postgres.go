package database

import (
	"fmt"
	"goproject/internal/domain/model"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	Host, Port, DBName, User, Password string
}
type Database struct {
	DB *gorm.DB
}

func NewPostgresDB() (*Database, error) {
	config := Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		DBName:   os.Getenv("DB_DATABASE"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable timezone=Asia/Jakarta",
		config.Host,
		config.User,
		config.Password,
		config.DBName,
		config.Port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(
		&model.List{},
		&model.User{},
		&model.Blog{},
		&model.Post{},
		&model.Comment{},
		&model.ListPost{},
	)
	if err != nil {
		return nil, err
	}

	return &Database{
		DB: db,
	}, nil
}
