# Ggami (까미) - Windows Server Website Builder

**Ggami** is a zero-dependency website builder designed for Windows Server environments. It generates single-binary Go web servers that include everything needed to run (assets, templates, database drivers).

## Key Features

- **Zero Dependency**: Generates a standalone `.exe` file.
- **MSSQL Integration**: Built-in support for SQL Server.
- **HTMX + Tailwind**: Modern, fast frontend without complex build steps.
- **Module System**: Inject pre-built features (Login, Hero, etc.) via the UI.
- **Polyglot Architecture**: Designed to support multiple languages (Go implemented, Node.js planned).

## Getting Started

### Prerequisites

- **Node.js**: v18+
- **Rust**: v1.70+ (with Visual Studio Build Tools)
- **Go**: v1.22+ (for running generated projects)

### Installation & Run

1.  Clone the repository:
    ```bash
    git clone https://github.com/x-ordo/GGAMI.git
    cd GGAMI
    ```
2.  Install dependencies:
    ```bash
    npm install
    npm run tauri dev
    ```

## Usage

1.  **Launch the Builder**: Run `npm run tauri dev`.
2.  **Configure**: Enter project name and DB connection details.
3.  **Select Modules**: Choose features like "Simple Login Form" from the sidebar.
4.  **Generate**: Click "Generate Project".
5.  **Run Server**: Go to the output folder and run `go run .`.

## License

MIT
