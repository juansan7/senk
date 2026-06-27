# Global Dependency TUI — Technical Specifications

## 1. Architecture Overview
The application will be built as a single-binary CLI tool using **Go (Golang)**. 
The UI will be driven by the **Elm Architecture** (Model, View, Update) using the `charmbracelet/bubbletea` framework. Styling will be handled by `charmbracelet/lipgloss`.

Because fetching global dependencies requires shelling out to various package managers (which can take 1-3 seconds), data fetching will be completely **asynchronous**. The UI will render immediately and use loading spinners for each ecosystem while data is being gathered.

## 2. Project Structure
```text
.
├── main.go               # Entry point, initializes Bubble Tea program
├── tui/                  # UI components
│   ├── model.go          # Bubble Tea model definition
│   ├── update.go         # Message handling and state transitions
│   ├── view.go           # UI rendering logic (lipgloss layouts)
│   └── styles.go         # Color and layout definitions
└── managers/             # Package manager implementations
    ├── manager.go        # Interface definition (PackageManager)
    ├── brew.go
    ├── npm.go
    ├── pip.go
    ├── gobin.go
    ├── cargo.go
    └── ...               # (gem, dotnet, composer, bun, pub, luarocks)
```

## 3. Data Models

### 3.1 Core Interfaces
```go
package managers

// Dependency represents a single installed package.
type Dependency struct {
	Name    string
	Version string
}

// PackageManager is the interface all ecosystem scanners must implement.
type PackageManager interface {
	Name() string
	// Fetch runs the underlying shell commands and parses the output.
	Fetch() ([]Dependency, error)
	// IsInstalled checks if the underlying CLI tool exists in $PATH.
	IsInstalled() bool
}
```

### 3.2 TUI State Model
```go
package tui

type ManagerState int
const (
	StateLoading ManagerState = iota
	StateDone
	StateError
)

type ManagerData struct {
	Manager      managers.PackageManager
	State        ManagerState
	Dependencies []managers.Dependency
	Err          error
}

type Model struct {
	Managers       []ManagerData
	SelectedIndex  int
	ListWidth      int
	ListHeight     int
	// Bubbles components
	Spinner        spinner.Model
}
```

## 4. Concurrency & Message Passing (Bubble Tea)
1. **Init**: The program fires off a `tea.Cmd` for each `PackageManager` that returns `IsInstalled() == true`.
2. **Fetch**: Each command runs in a separate Goroutine (`go exec.Command(...)`).
3. **Message**: When a Goroutine finishes, it returns a `DependenciesFetchedMsg{ Index: int, Deps: []Dependency }` or `FetchErrorMsg`.
4. **Update**: The Bubble Tea `Update` function catches these messages, updates the `Model.Managers[Index].State`, and triggers a UI re-render.

## 5. Ecosystem Fetching Strategies

Since we must extract both **Name** and **Version**, we rely on structured output (`--json`) where possible, or RegEx/String parsing where not.

| Ecosystem | Command | Parsing Strategy |
| :--- | :--- | :--- |
| **Homebrew** | `brew list --versions` | Space separated (`<name> <version>`) |
| **npm** | `npm list -g --depth=0 --json` | Parse JSON `dependencies` map |
| **pip** | `pip list --format=json` | Parse JSON array of `{name, version}` |
| **Go** | `ls $(go env GOPATH)/bin` | Execute `go version -m <bin>` per binary |
| **Cargo** | `cargo install --list` | Regex: `^([^ ]+) v([^:]+):` |
| **Ruby (gem)** | `gem list --local` | Regex: `^([^ ]+) \(([^)]+)\)` |
| **.NET** | `dotnet tool list -g` | Skip 2 header lines, split by spaces |
| **PHP (composer)**| `composer global show --format=json` | Parse JSON `installed` array |
| **Bun** | `bun pm ls -g` | Split by line, regex to extract `name@version` |
| **Dart (pub)** | `dart pub global list` | Space separated (`<name> <version>`) |
| **Lua (luarocks)**| `luarocks list` | Parse tabular output structure |

## 6. UI Layout (Lipgloss)
* **Left Pane (width: 30%)**: A selectable list of active package managers. Shows a spinner `⠋` if `StateLoading`, a checkmark `✓` if `StateDone`, or `✗` if `StateError`.
* **Right Pane (width: 70%)**: A table or list showing the packages for the currently selected manager.
  * *Empty State*: "Select a package manager..."
  * *Loading State*: "Fetching dependencies..."
  * *Populated State*: Two columns: **Name** (left aligned) and **Version** (right aligned).

## 7. Error Handling
- If a package manager CLI tool is not found in `$PATH`, it is omitted from the UI entirely (`IsInstalled() == false`).
- If a package manager fails to fetch data (e.g., permission error, malformed JSON), the Left Pane shows a red `✗` and the Right Pane displays the error message (`Err.Error()`).
- Shell command timeouts (e.g., `context.WithTimeout(..., 5*time.Second)`) will be used to prevent hanging the application if a manager hangs.
