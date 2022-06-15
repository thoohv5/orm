package adapter

import (
	_ "github.com/go-sql-driver/mysql"

	"github.com/thoohv5/orm/impl/gorm"
	"github.com/thoohv5/orm/standard"
)

type typeName string

const (
	Database typeName = "database"
	Gorm     typeName = "gorm"
)

// GetConnect 数据库连接适配
func GetConnect(typeName typeName) standard.IConnect {
	var connect standard.IConnect
	switch typeName {
	case Database:
		// connect = impl.NewSDb()
	case Gorm:
		connect = gorm.NewGDb()
	default:
		connect = gorm.NewGDb()
	}
	return connect
}
