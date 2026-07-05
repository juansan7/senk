# senk-tui (Dependencies TUI)

## Description

**senk-tui** is a Terminal User Interface (TUI) command-line tool written in Go that helps you manage globally installed dependencies across multiple programming languages and package managers on your system. It provides an intuitive, interactive interface to view, inspect, and uninstall global packages seamlessly.

## Features

- **Multi-Ecosystem Support:** Scans and supports global packages from Node.js (npm, bun), Python (pip), Golang (go bin), Rust (cargo), Ruby (gem), PHP (composer), .NET (dotnet), Dart (pub), Lua (luarocks), and macOS (brew).
- **Interactive TUI:** Built with Bubbletea, offering a responsive and easy-to-navigate interface.
- **Package Details:** Fetch deep metadata and sizes for globally installed packages.
- **Uninstall Packages:** Remove packages directly from the interface with a single keystroke.
- **Fast & Asynchronous:** Loads packages and executes actions concurrently without blocking the UI.

## Installation

First, make sure you have [Go](https://golang.org/dl/) installed on your system.

1. Clone this repository:

    ```sh
    git clone https://github.com/your_username/dependencies-tui.git
    cd dependencies-tui
    ```

2. Build the project:

    ```sh
    go build -ldflags="-s -w" -o bin/senk
    ```

3. (Optional) Move the executable to a directory included in your PATH to use it from anywhere:

    ```sh
    mv bin/senk /usr/local/bin/
    ```

## Usage

To use the tool, simply run the executable in your terminal:

```sh
senk
```

### Navigation
- Use **Up/Down/Left/Right** arrows or **k/j/h/l** (Vim bindings) to navigate between panes (Languages, Managers, Packages).
- Press **Enter** on a package to view its details.
- Press **x** inside a package's details modal to uninstall it.
- Press **Esc** to close modals.
- Press **q** or **Ctrl+C** to quit the application.

## Supported Ecosystems

*   **Node.js**: `npm`, `bun`
*   **Python**: `pip`
*   **Golang**: `go` (global binaries)
*   **Rust**: `cargo`
*   **Ruby**: `gem`
*   **PHP**: `composer`
*   **.NET**: `dotnet`
*   **Dart**: `pub`
*   **Lua**: `luarocks`
*   **macOS**: `brew`

## Contribution
Contributions are welcome! If you want to improve this tool, please follow these steps:

1. Fork the repository.
2. Create a new branch (`git checkout -b feature/new-feature`).
3. Make your changes and commit them (`git commit -am 'Add new feature'`).
4. Push your changes (`git push origin feature/new-feature`).
5. Open a Pull Request.

## License
This project is licensed under the MIT License. For more details, see the [LICENSE](LICENSE) file.
