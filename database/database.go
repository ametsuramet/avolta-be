package database

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"avolta/config"

	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var DB = &gorm.DB{}

// InitDB initializes the database connection
func InitDB(ctx context.Context) (*gorm.DB, error) {
	config.LoadConfig()

	if err := createDatabase(); err != nil {
		log.Println(err)
		return nil, err
	}

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
	schema.RegisterSerializer("json", JSONSerializer{})

	return db, nil
}

func createDatabase() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local",
		config.App.Database.DBUser, config.App.Database.DBPassword, config.App.Database.DBHost, config.App.Database.DBPort)
	DB, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	createDB := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci", config.App.Database.DBName)
	err := DB.Exec(createDB).Error
	if err != nil {
		return err
	}
	return nil
}

type JSONSerializer struct {
}

// Scan implements serializer interface
func (JSONSerializer) Scan(ctx context.Context, field *schema.Field, dst reflect.Value, dbValue interface{}) (err error) {
	fieldValue := reflect.New(field.FieldType)

	if dbValue != nil {
		var bytes []byte
		switch v := dbValue.(type) {
		case []byte:
			bytes = v
		case string:
			bytes = []byte(v)
		default:
			return fmt.Errorf("failed to unmarshal JSONB value: %#v", dbValue)
		}

		err = json.Unmarshal(bytes, fieldValue.Interface())
	}

	field.ReflectValueOf(ctx, dst).Set(fieldValue.Elem())
	return
}

// Value implements serializer interface
func (JSONSerializer) Value(ctx context.Context, field *schema.Field, dst reflect.Value, fieldValue interface{}) (interface{}, error) {
	return json.Marshal(fieldValue)
}
