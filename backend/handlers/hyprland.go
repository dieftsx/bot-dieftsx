package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type HyprWorkspace struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type HyprClient struct {
	Pid          int           `json:"pid"`
	Class        string        `json:"class"`
	Title        string        `json:"title"`
	InitialClass string        `json:"initialClass"`
	InitialTitle string        `json:"initialTitle"`
	Workspace    HyprWorkspace `json:"workspace"`
}

type TmuxPane struct {
	Pid     int
	Command string
	Path    string
}

type NvimStatusResponse struct {
	WM                string         `json:"wm"`
	TerminalWS2Active bool           `json:"terminal_ws2_active"`
	TerminalName      string         `json:"terminal_name,omitempty"`
	NvimWS2Active     bool           `json:"nvim_ws2_active"`
	FirefoxActive     bool           `json:"firefox_active"`
	NvimPid           int            `json:"nvim_pid,omitempty"`
	Cwd               string         `json:"cwd,omitempty"`
	ActiveFile        ActiveFileInfo `json:"active_file,omitempty"`
}

func detectWM() string {
	if os.Getenv("HYPRLAND_INSTANCE_SIGNATURE") != "" {
		return "Hyprland"
	}
	xdg := strings.ToLower(os.Getenv("XDG_CURRENT_DESKTOP"))
	if strings.Contains(xdg, "kde") || strings.Contains(xdg, "plasma") || os.Getenv("KDE_FULL_SESSION") != "" {
		return "KDE Plasma"
	}
	session := strings.ToLower(os.Getenv("DESKTOP_SESSION"))
	if strings.Contains(session, "plasma") || strings.Contains(session, "kde") {
		return "KDE Plasma"
	}
	if strings.Contains(session, "hyprland") {
		return "Hyprland"
	}
	if isProcessRunning("hyprland") {
		return "Hyprland"
	}
	if isProcessRunning("kwin_wayland") || isProcessRunning("kwin_x11") || isProcessRunning("plasmashell") {
		return "KDE Plasma"
	}
	return "Linux"
}

func isProcessRunning(name string) bool {
	entries, err := os.ReadDir("/proc")
	if err != nil {
		return false
	}
	nameLower := strings.ToLower(name)
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		pid, err := strconv.Atoi(entry.Name())
		if err != nil {
			continue
		}
		commPath := fmt.Sprintf("/proc/%d/comm", pid)
		data, err := os.ReadFile(commPath)
		if err != nil {
			continue
		}
		comm := strings.ToLower(strings.TrimSpace(string(data)))
		if comm == nameLower {
			return true
		}
	}
	return false
}

func getDescendants(parentPid int) []int {
	var descendants []int
	entries, err := os.ReadDir("/proc")
	if err != nil {
		return descendants
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		pid, err := strconv.Atoi(entry.Name())
		if err != nil {
			continue
		}

		statPath := fmt.Sprintf("/proc/%d/stat", pid)
		data, err := os.ReadFile(statPath)
		if err != nil {
			continue
		}

		content := string(data)
		lastParen := strings.LastIndex(content, ")")
		if lastParen == -1 || len(content) <= lastParen+2 {
			continue
		}

		fields := strings.Fields(content[lastParen+2:])
		if len(fields) >= 2 {
			ppid, err := strconv.Atoi(fields[1])
			if err == nil && ppid == parentPid {
				descendants = append(descendants, pid)
				descendants = append(descendants, getDescendants(pid)...)
			}
		}
	}
	return descendants
}

func isNvimProcess(pid int) bool {
	commPath := fmt.Sprintf("/proc/%d/comm", pid)
	data, err := os.ReadFile(commPath)
	if err != nil {
		return false
	}
	comm := strings.ToLower(strings.TrimSpace(string(data)))
	return comm == "nvim" || comm == "neovim"
}

func getTmuxPanes() []TmuxPane {
	var panes []TmuxPane
	cmd := exec.Command("tmux", "list-panes", "-a", "-F", "#{pane_pid} #{pane_current_command} #{pane_current_path}")
	out, err := cmd.Output()
	if err != nil {
		return panes
	}

	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, " ", 3)
		if len(parts) == 3 {
			pid, err := strconv.Atoi(parts[0])
			if err == nil {
				panes = append(panes, TmuxPane{
					Pid:     pid,
					Command: strings.ToLower(parts[1]),
					Path:    parts[2],
				})
			}
		}
	}
	return panes
}

func InspectEnvironment() (string, bool, bool, bool, string, int, string) {
	wm := detectWM()
	termClasses := map[string]string{
		"kitty": "Kitty", "konsole": "Konsole", "alacritty": "Alacritty",
		"foot": "Foot", "ghostty": "Ghostty", "wezterm": "WezTerm",
		"st": "ST", "xterm": "XTerm", "urxvt": "URxvt", "rio": "Rio",
		"gnome-terminal": "GNOME Terminal", "tilix": "Tilix", "xfce4-terminal": "XFCE Terminal",
	}

	foundTerm := false
	foundNvim := false
	foundFirefox := false
	terminalName := ""
	nvimPid := 0
	nvimCwd := ""

	// 1. Try Hyprland client inspection if running Hyprland
	if wm == "Hyprland" {
		cmd := exec.Command("hyprctl", "clients", "-j")
		if output, err := cmd.Output(); err == nil {
			var clients []HyprClient
			if err := json.Unmarshal(output, &clients); err == nil {
				tmuxPanes := getTmuxPanes()
				for _, c := range clients {
					classLower := strings.ToLower(c.Class)
					titleLower := strings.ToLower(c.Title)

					if strings.Contains(classLower, "firefox") || strings.Contains(titleLower, "firefox") || strings.Contains(c.InitialClass, "firefox") {
						foundFirefox = true
					}

					if tName, ok := termClasses[classLower]; ok {
						foundTerm = true
						if terminalName == "" {
							terminalName = tName
						}
					} else if classLower != "" && !strings.Contains(classLower, "firefox") && !strings.Contains(classLower, "chrome") {
						foundTerm = true
						if terminalName == "" {
							terminalName = c.Class
						}
					}

					if strings.Contains(titleLower, "nvim") || strings.Contains(titleLower, "neovim") ||
						strings.Contains(classLower, "nvim") || strings.Contains(classLower, "neovim") {
						foundNvim = true
						if nvimPid == 0 {
							nvimPid = c.Pid
						}
						if nvimCwd == "" {
							nvimCwd = getNvimCwd(c.Pid)
						}
					}

					if c.Pid > 0 {
						if isNvimProcess(c.Pid) {
							foundNvim = true
							if nvimPid == 0 {
								nvimPid = c.Pid
							}
							if nvimCwd == "" {
								nvimCwd = getNvimCwd(c.Pid)
							}
						}
						for _, dPid := range getDescendants(c.Pid) {
							if isNvimProcess(dPid) {
								foundNvim = true
								if nvimPid == 0 {
									nvimPid = dPid
								}
								if nvimCwd == "" {
									nvimCwd = getNvimCwd(dPid)
								}
							}
						}
					}
				}
				for _, pane := range tmuxPanes {
					if pane.Command == "nvim" || pane.Command == "neovim" {
						foundNvim = true
						if nvimPid == 0 {
							nvimPid = pane.Pid
						}
						if nvimCwd == "" {
							nvimCwd = pane.Path
						}
					}
				}
			}
		}
	}

	// 2. Scan /proc for running processes (Universal for KDE Plasma & Hyprland)
	entries, err := os.ReadDir("/proc")
	if err == nil {
		myPid := os.Getpid()
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			pid, err := strconv.Atoi(entry.Name())
			if err != nil || pid == myPid {
				continue
			}

			commPath := fmt.Sprintf("/proc/%d/comm", pid)
			data, err := os.ReadFile(commPath)
			if err != nil {
				continue
			}
			comm := strings.ToLower(strings.TrimSpace(string(data)))

			if comm == "firefox" || comm == "firefox-bin" || comm == "geckomain" {
				foundFirefox = true
			}

			if tName, ok := termClasses[comm]; ok {
				foundTerm = true
				if terminalName == "" {
					terminalName = tName
				}
			}

			if comm == "nvim" || comm == "neovim" {
				foundNvim = true
				if nvimPid == 0 {
					nvimPid = pid
				}
				if nvimCwd == "" {
					nvimCwd = getNvimCwd(pid)
				}
			}
		}
	}

	// Check tmux panes universally
	tmuxPanes := getTmuxPanes()
	for _, pane := range tmuxPanes {
		if pane.Command == "nvim" || pane.Command == "neovim" {
			foundNvim = true
			if nvimPid == 0 {
				nvimPid = pane.Pid
			}
			if nvimCwd == "" {
				nvimCwd = pane.Path
			}
		}
	}

	// Check active file / typing signal
	if !foundNvim {
		if _, ok := checkTypingFile(); ok {
			foundNvim = true
		} else if data, err := os.ReadFile("/tmp/obs_active_file.txt"); err == nil && len(strings.TrimSpace(string(data))) > 0 {
			foundNvim = true
		}
	}

	return wm, foundTerm, foundNvim, foundFirefox, terminalName, nvimPid, nvimCwd
}

func InspectWorkspace2() (bool, bool, string, int, string) {
	_, foundTerm, foundNvim, _, termName, nvimPid, nvimCwd := InspectEnvironment()
	return foundTerm, foundNvim, termName, nvimPid, nvimCwd
}

func getNvimCwd(pid int) string {
	if pid <= 0 {
		return ""
	}
	cwdPath := fmt.Sprintf("/proc/%d/cwd", pid)
	target, err := os.Readlink(cwdPath)
	if err != nil {
		return ""
	}
	return target
}

func isCodeOrConfigFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	validExts := map[string]bool{
		".go": true, ".js": true, ".ts": true, ".jsx": true, ".tsx": true,
		".py": true, ".rs": true, ".c": true, ".cpp": true, ".h": true,
		".hpp": true, ".java": true, ".php": true, ".rb": true, ".html": true,
		".css": true, ".scss": true, ".vue": true, ".svelte": true, ".sh": true,
		".json": true, ".yaml": true, ".yml": true, ".toml": true, ".md": true,
		".sql": true, ".lua": true, ".zig": true, ".nim": true, ".dart": true,
	}
	base := strings.ToLower(filepath.Base(path))
	if base == "dockerfile" || base == "makefile" || base == "cmakelists.txt" {
		return true
	}
	return validExts[ext]
}

func checkTypingFile() (ActiveFileInfo, bool) {
	data, err := os.ReadFile("/tmp/obs_typing.txt")
	if err != nil {
		return ActiveFileInfo{}, false
	}

	contentStr := string(data)
	parts := strings.SplitN(contentStr, "\n---OBS_SPLIT---\n", 2)
	if len(parts) < 2 {
		return ActiveFileInfo{}, false
	}

	fileName := strings.TrimSpace(parts[0])
	content := parts[1]
	ext := filepath.Ext(fileName)
	lang := strings.TrimPrefix(ext, ".")

	return ActiveFileInfo{
		FileName: fileName,
		Content:  content,
		Language: lang,
	}, true
}

func findActiveFile(nvimPid int, cwd string) ActiveFileInfo {
	// 0. Check real-time typing sync from /tmp/obs_typing.txt
	if typingInfo, ok := checkTypingFile(); ok {
		return typingInfo
	}

	var targetFile string

	// 1. Check /tmp/obs_active_file.txt (from Neovim Lua autocmd)
	if data, err := os.ReadFile("/tmp/obs_active_file.txt"); err == nil {
		fp := strings.TrimSpace(string(data))
		if info, err := os.Stat(fp); err == nil && !info.IsDir() {
			if cwd == "" || strings.HasPrefix(fp, cwd) {
				targetFile = fp
			}
		}
	}

	// 2. Check RPC sockets specifically for nvimPid and its descendant PIDs
	if targetFile == "" && nvimPid > 0 {
		uid := os.Getuid()
		pidsToCheck := append([]int{nvimPid}, getDescendants(nvimPid)...)

		for _, p := range pidsToCheck {
			sock := fmt.Sprintf("/run/user/%d/nvim.%d.0", uid, p)
			if info, err := os.Stat(sock); err == nil && !info.IsDir() {
				cmd := exec.Command("nvim", "--server", sock, "--remote-expr", "expand('%:p')")
				out, err := cmd.Output()
				if err == nil {
					fp := strings.TrimSpace(string(out))
					if fp != "" && !strings.HasPrefix(fp, "term://") {
						if info, err := os.Stat(fp); err == nil && !info.IsDir() {
							targetFile = fp
							break
						}
					}
				}
			}
		}
	}

	// 3. Check /proc/<pid>/cmdline for nvimPid or its child PIDs
	if targetFile == "" && nvimPid > 0 {
		pidsToCheck := append([]int{nvimPid}, getDescendants(nvimPid)...)
		for _, p := range pidsToCheck {
			cmdlinePath := fmt.Sprintf("/proc/%d/cmdline", p)
			data, err := os.ReadFile(cmdlinePath)
			if err != nil {
				continue
			}
			args := strings.Split(string(data), "\x00")
			for _, arg := range args {
				arg = strings.TrimSpace(arg)
				if arg == "" || strings.HasPrefix(arg, "-") {
					continue
				}
				fullPath := arg
				if !filepath.IsAbs(fullPath) && cwd != "" {
					fullPath = filepath.Join(cwd, arg)
				}
				if info, err := os.Stat(fullPath); err == nil && !info.IsDir() {
					targetFile = fullPath
					break
				}
			}
			if targetFile != "" {
				break
			}
		}
	}

	// 4. Check open file descriptors in /proc/<pid>/fd under cwd
	if targetFile == "" && nvimPid > 0 {
		pidsToCheck := append([]int{nvimPid}, getDescendants(nvimPid)...)
		var newestFile string
		var newestTime time.Time

		for _, p := range pidsToCheck {
			fdDir := fmt.Sprintf("/proc/%d/fd", p)
			entries, err := os.ReadDir(fdDir)
			if err != nil {
				continue
			}

			for _, entry := range entries {
				linkPath := filepath.Join(fdDir, entry.Name())
				target, err := os.Readlink(linkPath)
				if err != nil {
					continue
				}
				if cwd != "" && !strings.HasPrefix(target, cwd) {
					continue
				}
				if info, err := os.Stat(target); err == nil && !info.IsDir() {
					if isCodeOrConfigFile(target) && info.ModTime().After(newestTime) {
						newestTime = info.ModTime()
						newestFile = target
					}
				}
			}
		}
		if newestFile != "" {
			targetFile = newestFile
		}
	}

	// 5. Check git diff modified files in cwd
	if targetFile == "" && cwd != "" {
		cmd := exec.Command("git", "-C", cwd, "diff", "--name-only")
		out, err := cmd.Output()
		if err == nil {
			lines := strings.Split(strings.TrimSpace(string(out)), "\n")
			for _, l := range lines {
				l = strings.TrimSpace(l)
				if l != "" {
					fp := filepath.Join(cwd, l)
					if info, err := os.Stat(fp); err == nil && !info.IsDir() {
						targetFile = fp
						break
					}
				}
			}
		}
	}

	// 6. Universal directory walk: find most recently modified source/config file in cwd
	if targetFile == "" && cwd != "" {
		var newestFile string
		var newestTime time.Time

		_ = filepath.Walk(cwd, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if info.IsDir() {
				name := info.Name()
				if strings.HasPrefix(name, ".") || name == "node_modules" || name == "target" || name == "vendor" || name == "dist" || name == "build" {
					return filepath.SkipDir
				}
				return nil
			}
			if isCodeOrConfigFile(path) && info.ModTime().After(newestTime) {
				newestTime = info.ModTime()
				newestFile = path
			}
			return nil
		})

		if newestFile != "" {
			targetFile = newestFile
		}
	}

	if targetFile == "" {
		return ActiveFileInfo{
			FileName: "Terminal",
			Content:  "// Terminal ativo na Workspace 2 (Aguardando arquivo)...",
			Language: "bash",
		}
	}

	contentBytes, err := os.ReadFile(targetFile)
	if err != nil {
		return ActiveFileInfo{
			FileName: filepath.Base(targetFile),
			FilePath: targetFile,
			Content:  "// Erro ao ler arquivo: " + err.Error(),
			Language: "txt",
		}
	}

	lines := strings.Split(string(contentBytes), "\n")
	if len(lines) > 80 {
		lines = lines[:80]
	}
	liveContent := strings.Join(lines, "\n")

	ext := filepath.Ext(targetFile)
	lang := strings.TrimPrefix(ext, ".")

	return ActiveFileInfo{
		FilePath: targetFile,
		FileName: filepath.Base(targetFile),
		Content:  liveContent,
		Language: lang,
	}
}

func getGitInfo(dir string) (string, string) {
	cmdBranch := exec.Command("git", "-C", dir, "branch", "--show-current")
	outBranch, err := cmdBranch.Output()
	branch := strings.TrimSpace(string(outBranch))
	if err != nil || branch == "" {
		cmdRev := exec.Command("git", "-C", dir, "rev-parse", "--short", "HEAD")
		outRev, errRev := cmdRev.Output()
		if errRev == nil {
			branch = strings.TrimSpace(string(outRev))
		} else {
			branch = "-"
		}
	}

	cmdStatus := exec.Command("git", "-C", dir, "status", "--porcelain")
	outStatus, errStatus := cmdStatus.Output()
	if errStatus != nil {
		return branch, "Sem Git"
	}

	if len(strings.TrimSpace(string(outStatus))) == 0 {
		return branch, "Clean ✨"
	}
	return branch, "Modified 📝"
}

func extractInsertionsDeletions(stat string) int {
	linesSum := 0
	parts := strings.Split(stat, ",")
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if strings.Contains(p, "insertion") || strings.Contains(p, "deletion") {
			fields := strings.Fields(p)
			if len(fields) > 0 {
				n, err := strconv.Atoi(fields[0])
				if err == nil {
					linesSum += n
				}
			}
		}
	}
	return linesSum
}

func isBinaryPath(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	binaryExts := map[string]bool{
		".exe": true, ".bin": true, ".png": true, ".jpg": true, ".jpeg": true,
		".gif": true, ".ico": true, ".pdf": true, ".zip": true, ".tar": true,
		".gz": true, ".7z": true, ".o": true, ".a": true, ".so": true,
	}
	if binaryExts[ext] {
		return true
	}
	base := filepath.Base(path)
	if base == "bot-dieftsx" || base == "server" {
		return true
	}
	return false
}

func getProjectFiles(dir string) []string {
	if dir == "" {
		return nil
	}
	homeDir, _ := os.UserHomeDir()
	cleanDir := dir
	if strings.HasPrefix(cleanDir, "~") && homeDir != "" {
		cleanDir = filepath.Join(homeDir, strings.TrimPrefix(cleanDir, "~"))
	}

	var files []string
	cmd := exec.Command("git", "-C", cleanDir, "ls-files")
	out, err := cmd.Output()
	if err == nil && len(out) > 0 {
		lines := strings.Split(strings.TrimSpace(string(out)), "\n")
		for _, l := range lines {
			l = strings.TrimSpace(l)
			if l != "" && !isBinaryPath(l) {
				files = append(files, l)
				if len(files) >= 60 {
					break
				}
			}
		}
		if len(files) > 0 {
			return files
		}
	}

	_ = filepath.Walk(cleanDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			name := info.Name()
			if strings.HasPrefix(name, ".") || name == "node_modules" || name == "target" || name == "vendor" || name == "dist" || name == "build" {
				return filepath.SkipDir
			}
			return nil
		}
		rel, err := filepath.Rel(cleanDir, path)
		if err == nil && !strings.HasPrefix(rel, ".") && !isBinaryPath(rel) {
			files = append(files, rel)
			if len(files) >= 60 {
				return filepath.SkipDir
			}
		}
		return nil
	})

	return files
}

func updateAutoStats(dir string) {
	if dir == "" {
		return
	}
	cmdCommits := exec.Command("git", "-C", dir, "log", "--since=midnight", "--oneline")
	outCommits, err := cmdCommits.Output()
	commitsCount := 0
	if err == nil {
		lines := strings.Split(strings.TrimSpace(string(outCommits)), "\n")
		for _, l := range lines {
			if strings.TrimSpace(l) != "" {
				commitsCount++
			}
		}
	}

	totalLines := 0
	cmdDiff := exec.Command("git", "-C", dir, "diff", "--shortstat")
	outDiff, err := cmdDiff.Output()
	if err == nil {
		totalLines += extractInsertionsDeletions(string(outDiff))
	}

	cmdLogStat := exec.Command("git", "-C", dir, "log", "--since=midnight", "--shortstat")
	outLogStat, err := cmdLogStat.Output()
	if err == nil {
		totalLines += extractInsertionsDeletions(string(outLogStat))
	}

	CurrentData.Stats.Commits = strconv.Itoa(commitsCount)
	if totalLines > 0 {
		CurrentData.Stats.Lines = strconv.Itoa(totalLines)
	}
}

func StartHyprlandMonitor(interval time.Duration) {
	// High-speed 50ms typing monitor
	go func() {
		for {
			if typingInfo, ok := checkTypingFile(); ok {
				mu.Lock()
				CurrentData.ActiveFile = typingInfo
				mu.Unlock()
			}
			time.Sleep(50 * time.Millisecond)
		}
	}()

	// 1s Workspace, Environment & Git monitor
	go func() {
		log.Println("🔍 Environment Monitor iniciado (KDE Plasma & Hyprland)...")
		for {
			wmName, foundTerm, foundNvim, foundFirefox, termName, nvimPid, nvimCwd := InspectEnvironment()

			mu.Lock()
			CurrentData.SystemInfo["WM"] = wmName

			if foundTerm {
				if termName != "" {
					CurrentData.SystemInfo["Terminal"] = termName + " 🟢"
				} else {
					CurrentData.SystemInfo["Terminal"] = "Ativo 🟢"
				}
			} else {
				CurrentData.SystemInfo["Terminal"] = "Ausente 🔴"
			}

			if foundFirefox {
				CurrentData.SystemInfo["Browser"] = "Firefox 🟢"
			} else {
				CurrentData.SystemInfo["Browser"] = "Ausente 🔴"
			}

			if foundNvim {
				CurrentData.SystemInfo["Editor"] = "Neovim 🟢"
				if foundFirefox {
					CurrentData.SystemInfo["Sessão"] = "Neovim + Firefox 🚀"
				} else {
					CurrentData.SystemInfo["Sessão"] = wmName + ": Neovim ⚡"
				}

				activeFile := findActiveFile(nvimPid, nvimCwd)
				CurrentData.ActiveFile = activeFile

				if nvimCwd != "" {
					homeDir, _ := os.UserHomeDir()
					displayDir := nvimCwd
					if homeDir != "" && strings.HasPrefix(nvimCwd, homeDir) {
						displayDir = "~" + strings.TrimPrefix(nvimCwd, homeDir)
					}
					CurrentData.Project.Directory = displayDir

					branch, gitStatus := getGitInfo(nvimCwd)
					CurrentData.Project.Branch = branch
					CurrentData.Project.Status = gitStatus
					CurrentData.Project.Files = getProjectFiles(nvimCwd)

					updateAutoStats(nvimCwd)
				}
			} else if foundFirefox {
				CurrentData.SystemInfo["Editor"] = "Firefox 🌐 (Testes)"
				CurrentData.SystemInfo["Sessão"] = "Testando Projetos 🧪"

				if typingInfo, ok := checkTypingFile(); ok {
					CurrentData.ActiveFile = typingInfo
				} else {
					CurrentData.ActiveFile = ActiveFileInfo{
						FileName: "Mozilla Firefox",
						FilePath: "firefox",
						Content:  "// Testando projeto no Mozilla Firefox 🌐\n// Navegador aberto para testes de UI / Frontend...",
						Language: "html",
					}
				}

				if nvimCwd != "" {
					CurrentData.Project.Files = getProjectFiles(nvimCwd)
					updateAutoStats(nvimCwd)
				} else if CurrentData.Project.Directory != "" && CurrentData.Project.Directory != "~/" {
					cleanDir := CurrentData.Project.Directory
					homeDir, _ := os.UserHomeDir()
					if strings.HasPrefix(cleanDir, "~") {
						cleanDir = filepath.Join(homeDir, strings.TrimPrefix(cleanDir, "~"))
					}
					CurrentData.Project.Files = getProjectFiles(cleanDir)
					updateAutoStats(cleanDir)
				}
			} else if foundTerm {
				CurrentData.SystemInfo["Editor"] = "Terminal 💻"
				CurrentData.SystemInfo["Sessão"] = wmName + ": Terminal Ativo"

				if typingInfo, ok := checkTypingFile(); ok {
					CurrentData.ActiveFile = typingInfo
				} else {
					CurrentData.ActiveFile = ActiveFileInfo{
						FileName: "Terminal",
						FilePath: "",
						Content:  "// Terminal ativo em " + wmName + " (Aguardando Neovim ou Firefox)...",
						Language: "bash",
					}
				}
				if nvimCwd != "" {
					CurrentData.Project.Files = getProjectFiles(nvimCwd)
					updateAutoStats(nvimCwd)
				}
			} else {
				CurrentData.SystemInfo["Editor"] = "Ausente 🔴"
				CurrentData.SystemInfo["Sessão"] = wmName + ": Aguardando..."

				CurrentData.ActiveFile = ActiveFileInfo{
					FileName: "Sistema",
					FilePath: "",
					Content:  "// Aguardando abertura do Neovim ou Mozilla Firefox em " + wmName + "...",
					Language: "txt",
				}
			}
			mu.Unlock()

			time.Sleep(interval)
		}
	}()
}

func HyprlandStatusHandler(w http.ResponseWriter, r *http.Request) {
	if setCORS(w, r) {
		return
	}
	wm, foundTerm, foundNvim, foundFirefox, termName, nvimPid, cwd := InspectEnvironment()
	activeFile := findActiveFile(nvimPid, cwd)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(NvimStatusResponse{
		WM:                wm,
		TerminalWS2Active: foundTerm,
		TerminalName:      termName,
		NvimWS2Active:     foundNvim,
		FirefoxActive:     foundFirefox,
		NvimPid:           nvimPid,
		Cwd:               cwd,
		ActiveFile:        activeFile,
	})
}
