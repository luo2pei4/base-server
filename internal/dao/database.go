package dao

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	SupportedDBSqlite3 = "sqlite3"
)

var baseServerDB *gorm.DB

func InitDB(dbType string) error {
	var err error
	switch dbType {
	case SupportedDBSqlite3:
		baseServerDB, err = initSqlite3DB()
	default:
		return fmt.Errorf("database '%s' is not supported", dbType)
	}

	if err != nil {
		return err
	}

	// 查看users表是否存在
	if !baseServerDB.Migrator().HasTable("users") {
		return fmt.Errorf("table 'users' does not exist")
	}
	return nil
}

func initSqlite3DB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("./base.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
