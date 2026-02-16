# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Run Commands

```bash
wails dev          # Development mode with hot reload
wails build        # Production build → build/bin/ggami.exe
go build ./...     # Verify Go compilation (no Wails window)
```

No frontend build step is needed — the frontend is pure HTML/CSS/JS served via `//go:embed all:frontend`.

## Architecture

Ggami is a **Wails v2** desktop app (Go backend + HTML/HTMX frontend) that does two things:

1. **Code Generator** — Generates standalone Go/Node.js web server projects with MSSQL integration
2. **Visual Builder** — Drag-and-drop website builder that exports static HTML

### Two communication paths between frontend and backend

- **Wails Bindings (IPC):** `App` struct methods in `app.go` are exposed to JS as `window.go.main.App.*`. Used for data operations (file dialogs, project CRUD, code generation).
- **HTMX via Chi router:** `api/handler.go` serves HTML fragments at `/api/*` through `AssetServer.Handler`. Used for dynamic UI rendering (component library, canvas rendering, property forms).

### Key packages

| Package | Purpose |
|---------|---------|
| `app.go` | Wails-bound methods — the bridge between frontend JS and Go |
| `api/` | Chi HTTP handler returning HTMX partial HTML fragments |
| `internal/generator/` | Strategy pattern: `Generator` interface with `GoGenerator` and `NodeGenerator` |
| `internal/modules/` | Module registry — injectable code snippets with marker-based injection |
| `internal/templates/` | Embedded `.tmpl` files for generated project output (go_main, go_mod, html_index) |
| `internal/builder/` | Visual builder: Component model with `ToHTML()`, ProjectManager (CRUD + JSON persistence), Exporter |
| `frontend/` | Pure HTML + vanilla JS + HTMX + Tailwind CDN. No bundler. |

### Module injection system

Modules in `internal/modules/registry.go` define `CodeSnippet` entries with a `Target` file, `Marker` string (e.g., `// @INJECT_ROUTES`), and `Content` to inject. During generation (`go_strategy.go`), each snippet is inserted at its marker using `strings.Replace(marker, content + "\n" + marker)` — the marker is preserved so multiple modules can inject at the same point.

### Frontend structure

`frontend/index.html` uses iframe-based view switching between `generator.html` and `builder.html`. Each page has its own JS file in `frontend/assets/js/`. The builder uses HTMX for dynamic content (component library loads via `hx-get`, property forms via `hx-get/hx-post`), while the generator uses direct Wails binding calls.

## Conventions

- Korean comments are common in templates and some source files (this is intentional)
- Builder projects are saved as `.ggami.json` files
- Generated output templates use `{{PLACEHOLDER}}` syntax (not Go `text/template` — just `strings.Replace`)
- The `internal/templates/*.tmpl` files are raw text files embedded via `//go:embed`, not Go template files despite the `.tmpl` extension
