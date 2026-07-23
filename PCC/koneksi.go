package main

import (
	"os"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func koneksi() *gorm.DB {
	// membaca file .env
	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "cloud_udb"
	}

	// Menggunakan SQLite (file database akan otomatis terbuat dengan nama sesuai DB_NAME)
	db, err := gorm.Open(sqlite.Open(dbname+".db"), &gorm.Config{})
	
	if err != nil {
		panic(err)
	}
	return db
}