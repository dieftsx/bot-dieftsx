#!/usr/bin/env sh

# OBS Live Bot - Terminal Shortcuts & Functions

OBS_BOT_URL="http://localhost:8080"

# Unalias any previous aliases to avoid zsh function/alias name collisions
unalias break_live codar_live pausa cafe obs_break codar obs_codar update_obs set_stack fim_live 2>/dev/null || true

# Altera o status da stream para Pausa para o Café ☕
break_live() {
  local live_status="Pausa para o Café ☕"
  local enc_status
  enc_status=$(python3 -c "import urllib.parse; print(urllib.parse.quote('''$live_status'''))" 2>/dev/null || echo "Pausa")
  
  local res
  res=$(curl -s -w "%{http_code}" -X POST "$OBS_BOT_URL/api/update?status=$enc_status" -o /dev/null)
  if [ "$res" = "200" ]; then
    echo "☕ Status da live alterado para: $live_status"
  else
    echo "❌ Erro ao conectar ao Bot (Porta 8080). Certifique-se de que o backend Go está rodando!"
  fi
}

# Altera o status da stream para Coding Live 💻
codar_live() {
  local live_status="Coding Live 💻"
  local enc_status
  enc_status=$(python3 -c "import urllib.parse; print(urllib.parse.quote('''$live_status'''))" 2>/dev/null || echo "Coding")
  
  local res
  res=$(curl -s -w "%{http_code}" -X POST "$OBS_BOT_URL/api/update?status=$enc_status" -o /dev/null)
  if [ "$res" = "200" ]; then
    echo "💻 Status da live alterado para: $live_status"
  else
    echo "❌ Erro ao conectar ao Bot (Porta 8080). Certifique-se de que o backend Go está rodando!"
  fi
}

pausa() { break_live "$@"; }
cafe() { break_live "$@"; }
obs_break() { break_live "$@"; }
codar() { codar_live "$@"; }
obs_codar() { codar_live "$@"; }

# Atualiza diretório atual, branch git e status git na overlay
update_obs() {
  local current_dir
  current_dir="$(pwd | sed "s|^$HOME|~|")"
  local git_branch="-"
  local git_st="Sem Git"

  if git rev-parse --is-inside-work-tree > /dev/null 2>&1; then
    git_branch="$(git branch --show-current 2>/dev/null)"
    if [ -z "$git_branch" ]; then
      git_branch="$(git rev-parse --short HEAD 2>/dev/null)"
    fi

    if [ -z "$(git status --porcelain 2>/dev/null)" ]; then
      git_st="Clean ✨"
    else
      git_st="Modified 📝"
    fi
  fi

  local enc_dir enc_branch enc_st
  enc_dir=$(python3 -c "import urllib.parse; print(urllib.parse.quote('''$current_dir'''))" 2>/dev/null || echo "$current_dir")
  enc_branch=$(python3 -c "import urllib.parse; print(urllib.parse.quote('''$git_branch'''))" 2>/dev/null || echo "$git_branch")
  enc_st=$(python3 -c "import urllib.parse; print(urllib.parse.quote('''$git_st'''))" 2>/dev/null || echo "$git_st")

  local res
  res=$(curl -s -w "%{http_code}" -X POST "$OBS_BOT_URL/api/update_project?dir=$enc_dir&branch=$enc_branch&status=$enc_st" -o /dev/null)
  if [ "$res" = "200" ]; then
    echo "⚡ Overlay atualizada: Dir [$current_dir] | Branch [$git_branch] | Status [$git_st]"
  else
    echo "❌ Erro ao conectar ao Bot (Porta 8080)."
  fi
}

# Atualiza a tech stack exibida na overlay (ex: set_stack "Go, Next.js, Arch Linux")
set_stack() {
  if [ -z "$1" ]; then
    echo "Uso: set_stack \"Go, Next.js, Arch Linux\""
    return 1
  fi
  local stack="$1"
  local enc_stack
  enc_stack=$(python3 -c "import urllib.parse; print(urllib.parse.quote('''$stack'''))" 2>/dev/null || echo "$stack")
  
  local res
  res=$(curl -s -w "%{http_code}" -X POST "$OBS_BOT_URL/api/update_stack?stack=$enc_stack" -o /dev/null)
  if [ "$res" = "200" ]; then
    echo "🥞 Tech stack atualizada para: $stack"
  else
    echo "❌ Erro ao conectar ao Bot (Porta 8080)."
  fi
}

# Incrementa mensagem no contador de chat
inc_chat() {
  local count="${1:-1}"
  curl -s -X POST "$OBS_BOT_URL/api/inc_chat?count=$count" -o /dev/null
}

# Assistente de fim de live 100% automatizado
fim_live() {
  echo "📊 Finalizando live e consolidando estatísticas automáticas..."

  local res
  res=$(curl -s -w "%{http_code}" -X POST "$OBS_BOT_URL/api/end_live" -o /tmp/obs_end_live.json)
  
  if [ "$res" = "200" ]; then
    echo "🎉 Live finalizada com sucesso! (Estatísticas 100% automatizadas)"
    echo "📹 Overlay alternada para Tela de Encerramento (Webcam desligada na overlay)."
    
    if command -v obs-cmd >/dev/null 2>&1; then
      obs-cmd source hide "Webcam" 2>/dev/null || true
      obs-cmd source hide "Cam" 2>/dev/null || true
      obs-cmd source hide "Camera" 2>/dev/null || true
    fi
  else
    echo "❌ Erro ao conectar ao Bot (Porta 8080)."
  fi
}
