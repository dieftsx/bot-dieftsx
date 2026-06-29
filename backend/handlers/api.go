package handlers

import (
	"encoding/json"
	"net/http"
)

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
	Status         string            `json:"status"`
	SystemInfo     map[string]string `json:"system_info"`
	TechStack      []string          `json:"tech_stack"`
	Project        ProjectInfo       `json:"project"`
	Stats          LiveStats         `json:"stats"`
	LatestFollower string            `json:"latest_follower"`
	LatestSub      string            `json:"latest_sub"`
	LatestBitsUser string            `json:"latest_bits_user"`
	LatestBits     string            `json:"latest_bits"`
	HypeTrainLevel string            `json:"hype_train_level"`
}

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
	LatestFollower: "Aguardando...",
	LatestSub:      "Aguardando...",
	LatestBitsUser: "Aguardando...",
	LatestBits:     "0",
	HypeTrainLevel: "0",
}

func DataHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(CurrentData)
}

func UpdateStatusHandler(w http.ResponseWriter, r *http.Request) {
	if status := r.URL.Query().Get("status"); status != "" {
		CurrentData.Status = status
	}
	w.WriteHeader(http.StatusOK)
}

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

func AlertHandler(w http.ResponseWriter, r *http.Request) {
	if follower := r.URL.Query().Get("follower"); follower != "" {
		CurrentData.LatestFollower = follower
	}
	if sub := r.URL.Query().Get("sub"); sub != "" {
		CurrentData.LatestSub = sub
	}
	if bitsUser := r.URL.Query().Get("bits_user"); bitsUser != "" {
		CurrentData.LatestBitsUser = bitsUser
	}
	if bits := r.URL.Query().Get("bits"); bits != "" {
		CurrentData.LatestBits = bits
	}
	if hypeLevel := r.URL.Query().Get("hype_level"); hypeLevel != "" {
		CurrentData.HypeTrainLevel = hypeLevel
	}
	w.WriteHeader(http.StatusOK)
}
