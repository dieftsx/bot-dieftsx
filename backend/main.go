package main

import (
	"bot-dieftsx/handlers"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/api/data", handlers.DataHandler)
	http.HandleFunc("/api/update", handlers.UpdateStatusHandler)
	http.HandleFunc("/api/update_project", handlers.UpdateProjectHandler)
	http.HandleFunc("/api/end_live", handlers.EndLiveHandler)
	http.HandleFunc("/api/alert", handlers.AlertHandler)

	fs := http.FileServer(http.Dir("../frontend"))
	http.Handle("/", fs)

	log.Println("⚡ Bot rodando em http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
