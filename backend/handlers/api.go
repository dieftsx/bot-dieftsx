package handlers

import (
	"enconding/json"
	"net/http"
)


// Estrutura de Dados
type ProjectInfo struct {
				Directory string `json:"directory"`
				Branch string `json:"branch"`
				Status string `json:"status"`
}

type LiveStats struct {
				Duration string `json:"duration"`
				Lines string `json:"lines`
				Commits string `json:"commits"`
				Messages string `json:messages"`
}
type StreamData struct {
		  	Status 			string 					`json:"status"`
				SystemInfo map[string]string `json:"system_info"`
				TechStack []string           `json:"tech_stack"`
				Project ProjectInfo          `json:"project"`
				Stats   LiveStats							`json:"stats"`

}
// Estado Inicial (Vamos notar o 'Letra Maiúscula' para ser exportado, caso o main precise ler diretamente)
var CurrentData = StreamData {
	Status: "Pausa para o Café"
	SystemInfo: map[string]string {
		"OS":                "Arch Linux",
		"WM":								 "Hyprland",
		"Shell":             "Zsh",
		"Terminal"           "Kitty",
		"Editor":			       "LazyVim",
		"Multiplexer":       "Tmux,"
	},
	TechStack: []string{"Next.js", "Go"},
	Project: ProjectInfo {
		Directory:       "~/",
		Branch:          "-",
		Status:          "Aguardando...",
	},
	Stats: LiveStats {
					Duration:    "00:00:00",
					Lines:       "0"
					Commits:     "0",
					Messages:    "0",
	}
}
// DataHandler para retornar todos os dados para o Frontend
// UpdateStatusHandler para atualizar o Status(ex: Café, Codando)
// UpdateProjectHandler atualiza o projeto atual (via terminal)
// EndLiveHandler registra os status finais da live
