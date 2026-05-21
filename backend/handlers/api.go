package handlers

import (
	"encoding/json"
	"net/http"
)

// Estruturas de Dados
type ProjectInfo struct {
	Directory string `json:"directory"`
	Branch    string `json:"branch"`
	Status    string `json:"status"`
}

type LiveStats struct {
	Duration string `json:"duration"`
	Lines    string `json:"lines"`
	Commits  string `json:"commits"`
	Messages string `json:"messages"`
}

type StreamData struct {
	Status     string            `json:"status"`
	SystemInfo map[string]string `json:"system_info"`
	TechStack  []string          `json:"tech_stack"`
	Project    ProjectInfo       `json:"project"`
	Stats      LiveStats         `json:"stats"`
}

// Estado Inicial (Note o 'C' maiúsculo para ser exportado, caso o main precise ler diretamente)
var CurrentData = StreamData{
	Status: "Pausa para o Café",
	SystemInfo: map[string]string{
		"OS":          "Arch Linux",
		"WM":          "Hyprland",
		"Shell":       "Zsh",
		"Terminal":    "Kitty",
		"Editor":      "LazyVim",
		"Multiplexer": "Tmux",
	},
	TechStack: []string{"Next.js", "Go"},
	Project: ProjectInfo{
		Directory: "~/",
		Branch:    "-",
		Status:    "Aguardando...",
	},
	Stats: LiveStats{
		Duration: "00:00:00",
		Lines:    "0",
		Commits:  "0",
		Messages: "0",
	},
}

// DataHandler retorna todos os dados para o Frontend
func DataHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(CurrentData)
}

// UpdateStatusHandler atualiza o Status (ex: Café, Codando)
func UpdateStatusHandler(w http.ResponseWriter, r *http.Request) {
	if status := r.URL.Query().Get("status"); status != "" {
		CurrentData.Status = status
	}
	w.WriteHeader(http.StatusOK)
}

// UpdateProjectHandler atualiza o projeto atual (via terminal)
func UpdateProjectHandler(w http.ResponseWriter, r *http.Request) {
	if dir := r.URL.Query().Get("dir"); dir != "" {
		CurrentData.Project.Directory = dir
	}
	if branch := r.URL.Query().Get("branch"); branch != "" {
		CurrentData.Project.Branch = branch
	}
	if status := r.URL.Query().Get("status"); status != "" {
		CurrentData.Project.Status = status
	}
	w.WriteHeader(http.StatusOK)
}

// EndLiveHandler registra os stats finais da live
func EndLiveHandler(w http.ResponseWriter, r *http.Request) {
	if duration := r.URL.Query().Get("duration"); duration != "" {
		CurrentData.Stats.Duration = duration
	}
	if lines := r.URL.Query().Get("lines"); lines != "" {
		CurrentData.Stats.Lines = lines
	}
	if commits := r.URL.Query().Get("commits"); commits != "" {
		CurrentData.Stats.Commits = commits
	}
	if messages := r.URL.Query().Get("messages"); messages != "" {
		CurrentData.Stats.Messages = messages
	}
	w.WriteHeader(http.StatusOK)
}
