package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"go_main_final_project/pkg/db"
	"go_main_final_project/pkg/api"
	"go_main_final_project/pkg/handlers"
)

var Port = 7540

func init() {

	if portStr := os.Getenv("TODO_PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			Port = port
		}
	}
}


func main() {
	err := db.Init("scheduler.db")
	if err != nil {
		log.Fatalf("database initialization error: %v", err)
	}

	api.Init()

	http.HandleFunc("/api/nextdate", handlers.NextDateHandler)

	webDir := "./web"
	fileServer := http.FileServer(http.Dir(webDir))
	http.Handle("/", fileServer)

	log.Printf("running the web server on the port %d...\n", Port)
	log.Printf("open it in a browser http://localhost:%d/\n", Port)

	err = http.ListenAndServe(fmt.Sprintf(":%d", Port), nil)
	if err != nil {
		log.Fatalf("server startup error: %v", err)
	}
}