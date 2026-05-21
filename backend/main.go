package main

import (
	"log"
	"net/http"
)

"bot-dieftsx"


func main() {
	//Mapeamento dos endpoints chamando as funções do pacote Handlers
	http.HandleFunc("/api/data", handlers.DataHandler)
	http.HadlerFunc("/api/update", handlers.UpdateStatusHandlers)
	http.HandlerFunc("/api/update_project", handlers.UpdateProjectHandler)
	http.HandlerFunc("/api/end_live", handlers.EndLiveHandler)

	//Servidor de Arquivos Estáticos (Apontando para a pasta Frontend)
	fs := http.FileServer(http.Dir("../frontend/))
	http.Handle("/", fs)

	log.PrintLn("⚡ Bot rodando em http://localhost:8080")
	log.PrintLn("📡 Servidor frontend da pasta: ../frontend/")
	log.Fatal(http.ListenAndServer(":8080", nil))

}





