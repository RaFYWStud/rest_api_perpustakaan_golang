package database

import (
	"database/sql"
	"golang-tutorial/config"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectDB() (*gorm.DB, *sql.DB) {
	cfg := config.Get()

	// set sql logger
	sqlLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
			Colorful:                  false,
		})

	log.Println("connecting to databases")

	// open connection to database
	db, err := gorm.Open(postgres.Open(cfg.DbUri), &gorm.Config{
		Logger:                 sqlLogger,
		SkipDefaultTransaction: true,
		AllowGlobalUpdate:      false,
	})

	// check if connection error
	if err != nil {
		log.Fatalf("error connect sql. error : %v", err)
	}

	log.Println("success connect database")

	log.Println("set database connection configuration")

	sqlDB, err := db.DB()

	if err != nil {
		log.Fatalf("error set database connection config. error : %v", err)
	}

	sqlDB.SetMaxIdleConns(10)

	sqlDB.SetMaxOpenConns(100)

	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, sqlDB
}
