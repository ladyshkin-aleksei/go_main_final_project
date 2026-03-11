package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"go_main_final_project/pkg/api"
	"go_main_final_project/pkg/config"
	"go_main_final_project/pkg/db"
	"go_main_final_project/pkg/handlers"
)

func main() {
	cfg := config.Load()

	err := db.Init(cfg.DBFile)
	if err != nil {
		log.Fatalf("database initialization error: %v", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("error closing database: %v", err)
		}
	}()

	err = api.Init()
	if err != nil {
		log.Fatalf("API initialization error: %v", err)
	}

	http.HandleFunc("/api/nextdate", handlers.NextDateHandler)

	webDir := "./web"
	if customWebDir := os.Getenv("TODO_WEB_DIR"); customWebDir != "" {
		webDir = customWebDir
	}

	fileServer := http.FileServer(http.Dir(webDir))
	http.Handle("/", fileServer)

	log.Printf("running the web server on the port %d...\n", cfg.Port)
	log.Printf("open it in a browser http://localhost:%d/\n", cfg.Port)

	err = http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), nil)
	if err != nil {
		log.Fatalf("server startup error: %v", err)
	}
}
