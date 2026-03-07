package server

import (
	"go_main_final_project/pkg/api"
	"log"
	"net/http"
)

func Run() error {
	api.Init()

	log.Println("Сервер запущен на :7540")
	return http.ListenAndServe(":7540", nil)
}
