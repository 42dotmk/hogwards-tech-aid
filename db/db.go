package db

import (
	"fmt"
	"os"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DB SERVER CONN STR
func GetEnvConnStringOrDefault() string {
	dbPath, isFound := os.LookupEnv("DB_PATH")
	if !isFound {
		dbPath = "db.sqlite"
	}
	return dbPath
}



func Connect() (db *gorm.DB) {
	if db, err := gorm.Open(sqlite.Open(GetEnvConnStringOrDefault()), &gorm.Config{}); err != nil {
		fmt.Println("Connection to database failed", err)
		panic(err)
	} else {
		fmt.Println("Connected to database")
		return db
	}
}
