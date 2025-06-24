package server

import (
	"flag"
	"fmt"
	"golang-tutorial/config"
	"golang-tutorial/handler"
	"golang-tutorial/internal/database"
	"golang-tutorial/migrations"
	"golang-tutorial/repository"
	"golang-tutorial/service"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Run() {
	cfg := config.Get()

	db, sqlDB := database.ConnectDB()

	downFlag := flag.Bool("down", false, "Run database migration down")
	downAllFlag := flag.Bool("down-all", false, "Run all database migrations down")

	// Parse the flags
	flag.Parse()

	if *downFlag {
		log.Println("Running database migration down...")
		migrations.Down(sqlDB)
		log.Println("Successfully run database migration down.")
		return
	}

	if *downAllFlag {
		log.Println("Running all database migrations down...")
		migrations.DownAll(sqlDB)
		log.Println("Successfully run all database migrations down.")
		return
	}

	log.Println("Running database migration up...")

	migrations.Up(sqlDB)

	log.Println("Successfully run database migration up.")

	repo := repository.New(db)

	serv := service.New(repo)

	if cfg.IsProduction {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      r,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	handler.New(r, serv)

	// Start the server
	log.Printf("Server is running on port %d", cfg.Port)
	log.Fatal(srv.ListenAndServe())
}
