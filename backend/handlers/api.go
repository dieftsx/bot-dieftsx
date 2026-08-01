package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type ProjectInfo struct {
	Directory string   `json:"directory"`
	Branch    string   `json:"branch"`
	Status    string   `json:"status"`
	Files     []string `json:"files,omitempty"`
}

type LiveStats struct {
	Duration string `json:"duration"`
	Lines    string `json:"lines"`
	Commits  string `json:"commits"`
	Messages string `json:"messages"`
}

type ActiveFileInfo struct {
	FilePath string `json:"file_path"`
	FileName string `json:"file_name"`
	Content  string `json:"content"`
	Language string `json:"language"`
}

type StreamData struct {
	Status         string            `json:"status"`
	SystemInfo     map[string]string `json:"system_info"`
	TechStack      []string          `json:"tech_stack"`
	Project        ProjectInfo       `json:"project"`
	ActiveFile     ActiveFileInfo    `json:"active_file"`
	Stats          LiveStats         `json:"stats"`
	LatestFollower string            `json:"latest_follower"`
	LatestSub      string            `json:"latest_sub"`
	LatestBitsUser string            `json:"latest_bits_user"`
	LatestBits     string            `json:"latest_bits"`
	HypeTrainLevel string            `json:"hype_train_level"`
}

var (
	mu               sync.RWMutex
	LiveStartTime    = time.Now()
	ChatMessageCount = 0
	CurrentData      = StreamData{
		Status: "Coding Live 💻",
		SystemInfo: map[string]string{
			"OS":          "Arch Linux",
			"WM":          "Detectando...",
			"Shell":       "Zsh",
			"Terminal":    "Detectando...",
			"Editor":      "Detectando...",
			"Browser":     "Detectando...",
			"Multiplexer": "Tmux",
		},
		TechStack: []string{"Next.js", "Go"},
		Project: ProjectInfo{
			Directory: "~/",
			Branch:    "-",
			Status:    "Aguardando...",
		},
		ActiveFile: ActiveFileInfo{
			FileName: "terminal",
			Content:  "// Aguardando código ou sessão do Neovim...",
			Language: "bash",
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
)

func setCORS(w http.ResponseWriter, r *http.Request) bool {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return true
	}
	return false
}

func DataHandler(w http.ResponseWriter, r *http.Request) {
	if setCORS(w, r) {
		return
	}
	w.Header().Set("Content-Type", "application/json")

	mu.Lock()
	elapsed := time.Since(LiveStartTime)
	h := int(elapsed.Hours())
	m := int(elapsed.Minutes()) % 60
	s := int(elapsed.Seconds()) % 60
	CurrentData.Stats.Duration = fmt.Sprintf("%02d:%02d:%02d", h, m, s)
	mu.Unlock()

	mu.RLock()
	defer mu.RUnlock()
	json.NewEncoder(w).Encode(CurrentData)
}

func IncChatHandler(w http.ResponseWriter, r *http.Request) {
	if setCORS(w, r) {
		return
	}
	count := 1
	if cStr := r.URL.Query().Get("count"); cStr != "" {
		if c, err := strconv.Atoi(cStr); err == nil && c > 0 {
			count = c
		}
	}
	mu.Lock()
	ChatMessageCount += count
	CurrentData.Stats.Messages = strconv.Itoa(ChatMessageCount)
	mu.Unlock()
	w.WriteHeader(http.StatusOK)
}

func UpdateStatusHandler(w http.ResponseWriter, r *http.Request) {
	if setCORS(w, r) {
		return
	}
	if status := r.URL.Query().Get("status"); status != "" {
		mu.Lock()
		CurrentData.Status = status
		mu.Unlock()
		log.Printf("⚡ [API] Status alterado para: %s", status)
	}
	w.WriteHeader(http.StatusOK)
}

func UpdateProjectHandler(w http.ResponseWriter, r *http.Request) {
	if setCORS(w, r) {
		return
	}
	mu.Lock()
	if dir := r.URL.Query().Get("dir"); dir != "" {
		CurrentData.Project.Directory = dir
	}
	if branch := r.URL.Query().Get("branch"); branch != "" {
		CurrentData.Project.Branch = branch
	}
	if status := r.URL.Query().Get("status"); status != "" {
		CurrentData.Project.Status = status
	}
	mu.Unlock()
	w.WriteHeader(http.StatusOK)
}

func UpdateStackHandler(w http.ResponseWriter, r *http.Request) {
	if setCORS(w, r) {
		return
	}
	if stackStr := r.URL.Query().Get("stack"); stackStr != "" {
		parts := strings.Split(stackStr, ",")
		var cleaned []string
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p != "" {
				cleaned = append(cleaned, p)
			}
		}
		if len(cleaned) > 0 {
			mu.Lock()
			CurrentData.TechStack = cleaned
			mu.Unlock()
		}
	}
	w.WriteHeader(http.StatusOK)
}

func EndLiveHandler(w http.ResponseWriter, r *http.Request) {
	if setCORS(w, r) {
		return
	}
	mu.Lock()
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
	CurrentData.Status = "Live Finalizada 📊"
	mu.Unlock()
	log.Printf("📊 [API] Live encerrada! Duração: %s | Commits: %s | Linhas: %s | Mensagens: %s",
		CurrentData.Stats.Duration, CurrentData.Stats.Commits, CurrentData.Stats.Lines, CurrentData.Stats.Messages)
	w.WriteHeader(http.StatusOK)
}

func AlertHandler(w http.ResponseWriter, r *http.Request) {
	if setCORS(w, r) {
		return
	}
	mu.Lock()
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
	mu.Unlock()
	w.WriteHeader(http.StatusOK)
}
