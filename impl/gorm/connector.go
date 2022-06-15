package gorm

import (
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"

	"github.com/thoohv5/orm/standard"
)

// Connect 连接
func (g *gDb) Connect(config *standard.Config, sos ...standard.ServerOption) (standard.IBuilder, error) {
	var (
		err error
		gdb *gorm.DB
	)

	gdb, err = gorm.Open(getDial(config.Driver, config.Dsn), &gorm.Config{
		QueryFields: true,
	})
	if err != nil {
		return nil, err
	}
	_ = gdb.Use(&OpentracingPlugin{})
	dsn := CopyGDb(gdb, append(sos, standard.WithIsWrite(true))...)
	db, err := gdb.DB()
	if err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(time.Duration(config.ConnMaxLifeTime) * time.Second)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetMaxOpenConns(config.MaxOpenConns)

	slaveDial := make([]gorm.Dialector, 0, len(config.Slave))
	for _, cs := range config.Slave {
		slaveDial = append(slaveDial, getDial(config.Driver, cs.Dsn))
	}

	// 读写分离
	err = gdb.Use(dbresolver.Register(dbresolver.Config{
		Replicas: slaveDial,
		Policy:   dbresolver.RandomPolicy{},
	}))
	if err != nil {
		return nil, err
	}

	return dsn, err
}

func getDial(driver string, dsn string) (dial gorm.Dialector) {
	switch strings.ToUpper(driver) {
	case "MYSQL":
		dial = mysql.New(mysql.Config{
			DSN: dsn,
		})
	case "POSTGRESSQL":
		dial = postgres.New(postgres.Config{
			DSN: dsn,
		})
	}

	return
}
