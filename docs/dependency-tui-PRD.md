# Global Dependency TUI — Product Requirements Document

| Field | Value |
|---|---|
| **Version** | v1.1 |
| **Status** | Draft |
| **Author** | opencode |
| **Last updated** | 2026-06-26 |
| **Stakeholders** | User / Developer |

---

## 1. Executive Summary

> We are building a Terminal User Interface (TUI) tool in Go that automatically detects and lists globally installed programming languages, their globally installed packages, and system-level packages via OS package managers like Homebrew. This provides developers with a centralized, fast, and visual way right in the terminal to audit their global system environment (e.g., what CLI tools are installed via `npm -g`, `cargo`, or `brew`) without running multiple disparate commands.

---

## 2. Problem Statement

### 2.1 Problem description
> Over time, developers install dozens of global CLI tools and libraries using different package managers (`npm install -g`, `pip install`, `go install`, `cargo install`, `brew install`). Keeping track of what is installed globally, and through which package manager, is difficult. Users often forget what they have installed or end up with duplicate tools across ecosystems.

### 2.2 Current state
> To audit a system, a developer must remember and run `npm list -g`, `pip list`, `cargo install --list`, and `brew list`. There is no unified dashboard to view the entire global state of a developer's machine.

### 2.3 Why now
> With robust TUI frameworks in Go (like `bubbletea`), we can create a beautiful, async terminal dashboard that aggregates this system-wide data into one easy-to-navigate view.

---

## 3. Goals & Success Metrics

### 3.1 Goals
- Provide a unified, interactive terminal interface to view globally installed packages across multiple ecosystems.
- Support core language package managers: **npm (Node), pip (Python), cargo (Rust), go (Golang), gem (Ruby), dotnet (.NET), composer (PHP), bun (Bun), pub (Dart), luarocks (Lua)**.
- Support system/OS-level package managers: **Homebrew (macOS only)**.
- Ensure the tool remains responsive while fetching data asynchronously.
- Always list the package name alongside its installed version.

### 3.2 Non-goals (out of scope)
- Tracking or parsing local (project-level) dependencies (e.g., local `package.json` or `node_modules`).
- Package installation/uninstallation workflows (Read-only for version 1).

### 3.3 Success metrics (KPIs)

| Metric | Baseline | Target | Timeline |
|---|---|---|---|
| Startup time (UI render) | N/A | < 100ms | MVP Launch |
| Supported ecosystems | 0 | 11 (Brew, npm, pip, go, cargo, gem, dotnet, composer, bun, pub, luarocks) | MVP Launch |

---

## 4. Users & Stakeholders

### 4.1 Target users / personas

**Persona 1: The Power User / Polyglot Developer**
- Role / description: Software engineer with a highly customized system environment.
- Main pain point: Losing track of global binaries and tools installed over months/years.
- Key need from this product: A single dashboard to view all global tools and their versions.

---

## 5. Scope & Dependencies

### 5.1 In scope
- Detecting and listing global dependencies for:
  - **Homebrew** (`brew`)
  - **Node.js** (`npm global`)
  - **Python** (`pip global/user`)
  - **Go** (`$GOPATH/bin`)
  - **Rust** (`cargo install`)
  - **Ruby** (`gem local`)
  - **.NET** (`dotnet tool global`)
  - **PHP** (`composer global`)
  - **Bun** (`bun global`)
  - **Dart/Flutter** (`dart pub global`)
  - **Lua** (`luarocks`)
- Interactive TUI with navigation (up/down/enter/esc).

### 5.2 Out of scope
- Local project dependency resolution.
- Dependency vulnerability scanning.

---

## 6. Functional Requirements

> Priority: [P0] = must-have, [P1] = should-have, [P2] = nice-to-have

### 6.1 Package Manager Detection & Data Fetching

**[P0] FR-001 — Detect Homebrew packages and versions**
> The system should list all installed Homebrew formulae and casks.
Acceptance criteria: Executes `brew list --versions` or parses `brew info --json` to display packages and their versions.

**[P0] FR-002 — Detect global Node.js (npm) packages and versions**
> The system should list globally installed npm packages.
Acceptance criteria: Executes `npm list -g --depth=0 --json` to fetch global packages and extract their versions.

**[P0] FR-003 — Detect global Python (pip) packages and versions**
> The system should list globally installed pip packages outside of virtualenvs.
Acceptance criteria: Executes `pip list --format=json` in the global context to extract package names and versions.

**[P0] FR-004 — Detect global Go packages and versions**
> The system should list installed Go binaries.
Acceptance criteria: Scans `$GOPATH/bin` and executes `go version -m <binary>` to extract the module version for each installed package.

**[P0] FR-005 — Detect global Rust (Cargo) packages and versions**
> The system should list globally installed cargo binaries.
Acceptance criteria: Executes `cargo install --list` and parses the output to extract package names and versions (e.g. `ripgrep v13.0.0:`).

**[P1] FR-006 — Detect global Ruby (gem) packages and versions**
> The system should list installed Ruby gems.
Acceptance criteria: Executes `gem list --local` and parses the output for name and version.

**[P1] FR-007 — Detect global .NET (dotnet) tools and versions**
> The system should list global dotnet tools.
Acceptance criteria: Executes `dotnet tool list -g` and parses the output table.

**[P1] FR-008 — Detect global PHP (composer) packages and versions**
> The system should list global composer packages.
Acceptance criteria: Executes `composer global show --format=json` (or parses output) to extract names and versions.

**[P1] FR-009 — Detect global Bun packages and versions**
> The system should list globally installed Bun packages.
Acceptance criteria: Executes `bun pm ls -g` and parses the output.

**[P1] FR-010 — Detect global Dart (pub) packages and versions**
> The system should list global Dart/Flutter packages.
Acceptance criteria: Executes `dart pub global list` and parses the output.

**[P1] FR-011 — Detect global Lua (luarocks) packages and versions**
> The system should list luarocks packages.
Acceptance criteria: Executes `luarocks list` and parses the output.

### 6.2 Interactive TUI 

**[P0] FR-012 — Two-pane layout with Loading States**
> The TUI should have a left pane for package managers (Brew, npm, pip, etc.) and a right pane for the installed packages.

User story: _As a user, I want the UI to load instantly and show spinners for ecosystems that are still being queried, so I don't have to wait for everything to finish before interacting._

Acceptance criteria:
- [ ] Keyboard navigation (Arrow keys, J/K).
- [ ] Async data loading (goroutines) with visual indicators (spinners) while shell commands execute.

---

## 7. Non-Functional Requirements

### 7.1 Performance
- Application UI must launch within 100ms. Because querying global environments (like `brew list`) takes time, data fetching must happen asynchronously in the background.

### 7.2 Compatibility
- Supported operating systems: macOS only (for the MVP/version 1).

---

## 8. Technical Considerations

### 8.1 Tech stack constraints
- Language: Go (Golang)
- Recommended UI Library: `charmbracelet/bubbletea` and `charmbracelet/bubbles`.

### 8.2 Data Fetching Strategy
> Since we are looking at global dependencies, we cannot just parse local files quickly. We must run shell commands (like `brew info --json` or `npm list -g --json`). To keep the TUI fast, we will start the TUI immediately and fire off Go routines (`tea.Cmd`) to fetch data for each ecosystem concurrently.
