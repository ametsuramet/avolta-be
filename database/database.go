package database

import (
	"context"
	"fmt"

	"avolta/config"

	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB = &gorm.DB{}

// InitDB initializes the database connection
func InitDB(ctx context.Context) (*gorm.DB, error) {
	config.LoadConfig()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.App.Database.DBUser, config.App.Database.DBPassword, config.App.Database.DBHost, config.App.Database.DBPort, config.App.Database.DBName)
	var cfg gorm.Config

	if config.App.Server.Mode != "release" {
		cfg = gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		}
	}
	// Open a connection to the database
	db, err := gorm.Open(mysql.Open(dsn), &cfg)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	DB = db.WithContext(ctx)

	return db, nil
}
