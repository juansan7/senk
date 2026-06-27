# ADR 0002: Filter Default/System Dependencies (All Ecosystems)

## 1. Context and Problem Statement
When viewing global dependencies, users see a massive list of packages. Many of these are "system/default" packages (shipped with the OS) or "implicit/transitive" dependencies (installed automatically because a top-level tool needed them). The user only wants to see **top-level packages they explicitly installed**.

## 2. Analysis of Ecosystems & Decision
We evaluated all 11 supported ecosystems to ensure we only display top-level/user-installed packages.

### Actions Required:
*   **Homebrew**: Change strategy. We will execute `brew leaves` (which returns only explicitly installed formulas that nothing else depends on) and `brew list --cask`, then cross-reference those names with `brew list --versions`.
*   **Ruby (gem)**: Change shell command to `gem list --local --no-default` to hide built-in OS gems.
*   **Python (pip)**: Change shell command to `pip list --not-required --format=json`. The `--not-required` flag hides packages that are merely dependencies of other packages.
*   **PHP (composer)**: Change shell command to `composer global show --direct --format=json`. The `--direct` flag limits the output to packages explicitly requested by the user.

### No Action Required (Already filtering correctly):
*   **Node.js (npm)**: Already uses `--depth=0`, which correctly limits output to top-level installs.
*   **Go (GOPATH)**: Scanning the `bin/` directory inherently only shows explicitly built/installed binaries (Go statically compiles transitive dependencies into the binary).
*   **Rust (Cargo)**: `cargo install --list` inherently only lists user-installed binaries.
*   **.NET (dotnet)**: `dotnet tool list -g` inherently only lists top-level tools.
*   **Dart (pub)**: `dart pub global list` inherently only lists activated top-level packages.
*   **Bun / Lua**: Their global behavior acts primarily on explicitly requested tools/binaries.

## 3. Consequences
*   **Positive**: The lists across all applicable managers will be drastically smaller and highly relevant to the user's explicit intent. 
*   **Negative**: Homebrew's fetching logic will become slightly more complex (requiring two shell commands), but it will remain asynchronous.
