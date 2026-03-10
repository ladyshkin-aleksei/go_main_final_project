package server

import (
	"go_main_final_project/pkg/api"
	"log"
	"net/http"
	"os"
	"strconv"
)

func Run() error {
	port := 7540

	if portStr := os.Getenv("TODO_PORT"); portStr != "" {
		if p, err := strconv.Atoi(portStr); err == nil {
			port = p
		}
	}

	api.Init()
	log.Printf("the server is running on :%d", port)

	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
