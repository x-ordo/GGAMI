# Ggami (까미) - Windows Server Website Builder

**Ggami** is a zero-dependency website builder designed for Windows Server environments. It generates single-binary Go web servers that include everything needed to run (assets, templates, database drivers). Built with **Wails v2 (Go) + HTMX + TailwindCSS**.

## Key Features

- **Zero Dependency**: Generates a standalone `.exe` file.
- **MSSQL Integration**: Built-in support for SQL Server.
- **HTMX + Tailwind**: Modern, fast frontend without complex build steps.
- **Module System**: Inject pre-built features (Login, Hero, etc.) via the UI.
- **Visual Builder**: No-code drag-and-drop website builder with HTML export.
- **Polyglot Architecture**: Supports Go and Node.js code generation.

## Getting Started

### Prerequisites

- **Go**: v1.22+
- **Wails CLI**: v2 (`go install github.com/wailsapp/wails/v2/cmd/wails@latest`)

### Installation & Run

1.  Clone the repository:
    ```bash
    git clone https://github.com/x-ordo/GGAMI.git
    cd GGAMI
    ```
2.  Run in development mode:
    ```bash
    wails dev
    ```
3.  Build for production:
    ```bash
    wails build
    ```

## Usage

1.  **Launch the Builder**: Run `wails dev`.
2.  **Code Generator**: Configure project name, DB connection, select modules, and generate.
3.  **Visual Builder**: Drag-and-drop components, edit properties, and export as static HTML.
4.  **Run Generated Server**: Go to the output folder and run `go mod tidy && go run .`.

## License

MIT
