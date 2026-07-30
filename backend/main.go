package main

import (
	"bot-dieftsx/handlers"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func main() {
	handlers.StartHyprlandMonitor(2 * time.Second)

	http.HandleFunc("/api/data", handlers.DataHandler)
	http.HandleFunc("/api/update", handlers.UpdateStatusHandler)
	http.HandleFunc("/api/update_project", handlers.UpdateProjectHandler)
	http.HandleFunc("/api/update_stack", handlers.UpdateStackHandler)
	http.HandleFunc("/api/end_live", handlers.EndLiveHandler)
	http.HandleFunc("/api/alert", handlers.AlertHandler)
	http.HandleFunc("/api/hyprland", handlers.HyprlandStatusHandler)
	http.HandleFunc("/api/inc_chat", handlers.IncChatHandler)

	frontendDir := "./frontend"
	if _, err := os.Stat(frontendDir); os.IsNotExist(err) {
		frontendDir = "../frontend"
	}
	absPath, _ := filepath.Abs(frontendDir)
	log.Printf("📁 Servindo arquivos estáticos de: %s", absPath)

	fs := http.FileServer(http.Dir(frontendDir))
	http.Handle("/", fs)

	log.Println("⚡ Bot rodando em http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}


