package main

import (
	"log"
	"net/http"

	// Importa o seu pacote de handlers local.

	"bot-dieftsx/handlers"
)

func main() {
	// Mapeia os endpoints chamando as funções do pacote handlers
	http.HandleFunc("/api/data", handlers.DataHandler)
	http.HandleFunc("/api/update", handlers.UpdateStatusHandler)
	http.HandleFunc("/api/update_project", handlers.UpdateProjectHandler)
	http.HandleFunc("/api/end_live", handlers.EndLiveHandler)

	// Servidor de Arquivos Estáticos (Apontando para a pasta frontend)
	fs := http.FileServer(http.Dir("../frontend"))
	http.Handle("/", fs)

	log.Println("⚡ Bot rodando em http://localhost:8080")
	log.Println("📡 Servindo frontend da pasta: ../frontend")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
