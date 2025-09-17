package models

import (
	"album/configs"
	"errors"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	InstanceSqlLite int = iota
	InstanceMySQL
)

var (
	DB                            *gorm.DB
	errInvalidSQLDatabaseInstance = errors.New("invalid sql db instance")
)

func GetModels() []interface{} {
	return []interface{}{&Album{}, &Category{}}
}

func NewDatabaseSQLFactory(instance int) (db *gorm.DB, err error) {
	switch instance {
	case InstanceMySQL:
		dns := fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
			configs.Config.DBUser,
			configs.Config.DBPassword,
			configs.Config.DBHost,
			configs.Config.DBPort,
			configs.Config.DBName)
		db, err = gorm.Open(mysql.Open(dns), &gorm.Config{})
	case InstanceSqlLite:
		db, err = gorm.Open(sqlite.Open(configs.Config.DBName), &gorm.Config{})
	default:
		return nil, errInvalidSQLDatabaseInstance
	}
	return db, err
}

func SetDatabase(instance int) (err error) {
	db, err := NewDatabaseSQLFactory(instance)
	if err != nil {
		return err
	}
	DB = db
	return err
}
