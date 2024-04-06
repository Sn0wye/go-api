package db

import (
	"fmt"

	"github.com/Sn0wye/go-api/pkg/config"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

func GetDB() *gorm.DB {
	conf := config.GetConfig()
	driver := conf.Get("db.driver")
	conn := getConnectionString(conf)

	if driver == "postgres" {
		db, err := gorm.Open(postgres.Open(conn), &gorm.Config{
			Logger: gormLogger.Default.LogMode(gormLogger.Info),
		})

		if err != nil {
			panic("Failed to connect to database")
		}

		return db
	}

	// sqlite (default)
	db, err := gorm.Open(sqlite.Open("db.sqlite"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	return db
}

func getConnectionString(conf *viper.Viper) string {
	connRaw := conf.Get("db.connectionString")
	conn, ok := connRaw.(string)

	if !ok {
		fmt.Println("Connection string is not defined in db.connectionString")
	}

	return conn
}
