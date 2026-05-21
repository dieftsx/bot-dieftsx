# ⚡ OBS Live Bot

> **Automatize sua overlay de stream direto do terminal usando Go, JS Vanilla e Arch Linux.**

Um bot local extremamente leve e rápido projetado para desenvolvedores. Ele extrai informações em tempo real do seu ambiente de trabalho (diretórios, status do Git, métricas da sessão) e projeta uma interface cyberpunk/neon na sua stream através do _Browser Source_ do OBS Studio.

---

## ✨ Funcionalidades em Destaque

- **🔄 Sincronização em Tempo Real:** Atualize status e stack tecnológica direto do terminal (ex: alternando de `Go` para `Next.js`).
- **🐙 Git Aware:** O painel "Now Working On" rastreia automaticamente seu diretório, branch atual e informa se a árvore de trabalho está limpa ou com modificações.
- **📊 Live Stats Automático:** Ao encerrar a live, um script consolida seus commits do dia e gera um relatório visual instantâneo para o chat.
- **🚀 Zero Impacto na CPU:** Escrito em Go e servindo arquivos HTML estáticos. Não consome recursos que deveriam ir para sua IDE ou compilador.

---

## 🏗️ Arquitetura e Estrutura

O projeto divide claramente a lógica de dados da apresentação visual, facilitando a manutenção e futuras migrações.

| Camada        | Tecnologia        | Descrição                                                       |
| :------------ | :---------------- | :-------------------------------------------------------------- |
| **Backend**   | `Go 1.21+`        | Servidor HTTP e gerenciamento de estado da live em memória.     |
| **Frontend**  | `HTML / CSS / JS` | Telas estáticas otimizadas para o OBS com design em CSS Grid.   |
| **Automação** | `Zsh / Shell`     | Scripts locais para comunicar o seu terminal com o servidor Go. |

<details>
<summary><strong>Ver Árvore de Diretórios</strong></summary>

```text
obs-live-bot/
├── backend/
│   ├── go.mod
│   ├── main.go
│   └── handlers/
│       └── api.go
├── frontend/
│   ├── index.html           # Cena: Pausa para o Café
│   ├── coding.html          # Widget: Bate Papo
│   ├── main_scene.html      # Cena: Desktop Principal
│   └── ending.html          # Cena: Encerramento
└── scripts/
    └── zsh_aliases.sh       # Funções para o ~/.zshrc
```

</details>

## 🚀 Instalação e Uso

### 1. Preparando o Terreno

Certifique-se de ter o Go instalado. Testado nativamente em ambiente Arch Linux / Hyprland.

### 2. Inicializando o Servidor

Clone o repositório e inicie o backend em Go:

```bash
git clone https://github.com/dieftsx/obs-live-bot.git
cd obs-live-bot/backend

# Inicia o servidor na porta 8080
go run main.go
```

> **Nota:** Mantenha esta aba do terminal rodando em segundo plano (ou utilize o Tmux) durante a sua live.

### 3. Configurando o OBS Studio

Adicione as cenas abaixo utilizando a fonte **Navegador (Browser)**:

| Cena               | URL                                             | Resolução            |
| :----------------- | :---------------------------------------------- | :------------------- |
| Cena Principal     | `http://localhost:8080/main_scene.html`         | 1920x1080            |
| Cena de Pausa      | `http://localhost:8080/index.html`              | Customizada          |
| Cena de Encerramento | `http://localhost:8080/ending.html`           | Customizada          |

> **💡 Dica do OBS:** Lembre-se de limpar qualquer CSS customizado na janela de propriedades da fonte e marcar a opção para não renderizar o fundo se ele não estiver sendo exibido.

## 💻 Comandos de Terminal (Zsh)

Para não precisar sair do seu fluxo de código, o bot é controlado via atalhos no terminal. Adicione o conteúdo de `scripts/zsh_aliases.sh` ao seu arquivo `~/.zshrc`.

| Comando       | Ação Executada na Stream                                    |
| :------------ | :---------------------------------------------------------- |
| `break`       | Altera o status superior para "Pausa para o Café ☕".        |
| `codar`       | Retorna o status superior para "Coding Live 💻".            |
| `update_obs`  | Envia o path atual (pwd) e a branch ativa para a overlay.   |
| `fim_live`    | Inicia o assistente de encerramento, colhendo commits diários. |

## 🔌 Referência da API (Endpoints REST)

O backend roda em `localhost:8080` e aceita as seguintes requisições:

| Método | Rota                                | Descrição                                      |
| :----- | :---------------------------------- | :--------------------------------------------- |
| `GET`  | `/api/data`                         | Retorna o JSON completo de estado.             |
| `POST` | `/api/update?status={texto}`        | Sobrescreve o status da stream.                |
| `POST` | `/api/update_project?dir={X}&branch={Y}&status={Z}` | Atualiza métricas do Git. |
| `POST` | `/api/end_live?duration={X}&commits={Y}...` | Alimenta a tela de estatísticas.  |

## 🤝 Let's Connect & Contribute

Este projeto foi construído publicamente. Sinta-se à vontade para fazer um fork, melhorar a estilização ou adicionar novas integrações!

- **🐙 GitHub:** [github.com/dieftsx](https://github.com/dieftsx)
- **💼 LinkedIn:** [linkedin.com/in/diefersonsoares](https://linkedin.com/in/diefersonsoares)
