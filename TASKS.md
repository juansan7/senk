# TASKS (senk-tui)

## Milestone 1: Project Scaffolding & Core Interfaces
- [x] **Task 1.1:** Install `charmbracelet/bubbletea` and `lipgloss`. Set up `main.go` with a basic Bubble Tea model that immediately quits (smoke test).
- [x] **Task 1.2:** Define the core `PackageManager` interface and `Dependency` struct in `managers/manager.go`.

## Milestone 2: Parsers & TDD (Pure Functions)
*Note: Implement the parser and its tests first (using `testdata/`), before writing the shell command execution logic.*
- [x] **Task 2.1:** Implement `npm` parser and tests (`managers/npm_test.go`, `managers/npm.go`).
- [x] **Task 2.2:** Implement `brew` parser and tests (`managers/brew_test.go`, `managers/brew.go`).
- [x] **Task 2.3:** Implement `pip` parser and tests (`managers/pip_test.go`, `managers/pip.go`).
- [x] **Task 2.4:** Implement `go` parser and tests (`managers/gobin_test.go`, `managers/gobin.go`).
- [x] **Task 2.5:** Implement `cargo` parser and tests (`managers/cargo_test.go`, `managers/cargo.go`).

## Milestone 3: Command Execution (Side-Effects)
- [x] **Task 3.1:** Implement the `Fetch()` method for the MVP managers (npm, brew, pip, go, cargo) using `exec.Command` and context timeouts.

## Milestone 4: TUI Integration
- [x] **Task 4.1:** Build the Left Pane (Managers List) with focus state and keyboard navigation (`j`/`k`).
- [x] **Task 4.2:** Implement concurrent fetching using `tea.Cmd` and loading spinners (`charmbracelet/bubbles/spinner`).
- [x] **Task 4.3:** Build the Right Pane (Dependencies List) and implement navigation between left and right panes (`Tab`/`Esc`).
- [x] **Task 4.4:** Apply styling using `lipgloss` to match the two-pane layout specification.

## Milestone 5: Extended Ecosystems (Tier 2)
- [x] **Task 5.1:** Implement `gem` (Ruby).
- [x] **Task 5.2:** Implement `dotnet` (.NET).
- [x] **Task 5.3:** Implement `composer` (PHP).
- [x] **Task 5.4:** Implement `bun`.
- [x] **Task 5.5:** Implement `pub` (Dart).
- [x] **Task 5.6:** Implement `luarocks` (Lua).
