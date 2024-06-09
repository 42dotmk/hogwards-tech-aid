package db

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	dbServer "gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DB SERVER CONN STR
func sqlite_conn_str() string {
	dbPath, isFound := os.LookupEnv("DB_PATH")
	if !isFound {
		dbPath = "db.sqlite"
	}
	return dbPath
}

var DB = func() (db *gorm.DB) {

	godotenv.Load()
	if db, err := gorm.Open(dbServer.Open(sqlite_conn_str()), &gorm.Config{}); err != nil {
		fmt.Println("Connection to database failed", err)
		panic(err)
	} else {
		fmt.Println("Connected to database")
		return db
	}
}()
