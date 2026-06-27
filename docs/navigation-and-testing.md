# Navigation & Testing Strategy

## 1. UI Navigation & Interaction Design

The interface relies on a dual-pane layout. To ensure smooth keyboard navigation, we will use an **Explicit Focus** model.

### 1.1 Layout Structure
- **Left Pane (Managers List):** Displays the available ecosystems (Brew, npm, pip, etc.). Indicates status (loading, done, error).
- **Right Pane (Packages List):** Displays the parsed dependencies for the currently active manager. Includes version numbers.

### 1.2 Focus States
- **Focus: Left Pane (Default):**
  - `Up` / `k`: Select previous package manager. (Right pane automatically updates).
  - `Down` / `j`: Select next package manager. (Right pane automatically updates).
  - `Tab` / `Right` / `l` / `Enter`: Move focus to the Right Pane.
- **Focus: Right Pane:**
  - `Up` / `k`: Scroll packages list up.
  - `Down` / `j`: Scroll packages list down.
  - `Tab` / `Left` / `h` / `Esc`: Return focus to the Left Pane.

### 1.3 Global Keybindings
- `q` / `ctrl+c`: Quit the application.
- `/`: (Future feature) Open a filter input to search packages in the right pane.

---

## 2. Testing Strategy (TDD Approach)

To follow Test-Driven Development (TDD) and Go best practices, we must separate **side-effects** (executing shell commands) from **business logic** (parsing CLI outputs). 

### 2.1 Separation of Concerns
Instead of writing a monolithic `Fetch()` method, each package manager will have two core parts:
1. **The Fetcher (Side-effect):** Runs `exec.Command("npm", "list", ...)`. Hard to test.
2. **The Parser (Pure Function):** e.g., `ParseNPMOutput(data []byte) ([]Dependency, error)`. Easy to test.

### 2.2 Table-Driven Tests
We will use standard Go Table-Driven Tests for all parsers. We will write the test *before* implementing the parsing logic.

Example structure for `managers/npm_test.go`:
```go
func TestParseNPMOutput(t *testing.T) {
    tests := []struct {
        name     string
        input    []byte
        expected []Dependency
        hasError bool
    }{
        {
            name: "Valid JSON output",
            input: []byte(`{"dependencies":{"express":{"version":"4.17.1"}}}`),
            expected: []Dependency{{Name: "express", Version: "4.17.1"}},
            hasError: false,
        },
        // Write this test first, watch it fail, then write the code!
    }
    // ... test runner logic
}
```

### 2.3 `testdata` Directory
For complex outputs (like `brew list` or `cargo install --list`), we will use Go's standard `testdata/` directory pattern. We will save real outputs from these tools into files (e.g., `testdata/brew_output.txt`) and load them in our tests to ensure our Regex and string manipulation works against real-world data.

### 2.4 TUI Testing
Testing Bubble Tea TUI components can be tricky. We will use the `teatest` package (from `charmbracelet/x/exp/teatest`) to send keystrokes programmatically and assert the view output to ensure navigation (Tab, Up, Down) works as expected.
