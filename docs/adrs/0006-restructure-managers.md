# 0006: Restructure Managers Package by Ecosystem

## Status
Proposed

## Context
The `managers` directory currently contains a flat list of 28 files, including all language definitions, manager implementations, extended methods, test files, and testdata. Large global files like `languages.go` group unrelated language definitions, and `fallback_extended.go` aggregates implementations across vastly different package managers (Go, Cargo, Gem, Dotnet, etc.). This makes the codebase difficult to navigate, violates the principle of separation of concerns, and hinders scalability as we add more ecosystems.

## Decision
We will refactor the `managers` package to use a modular, directory-based structure organized by ecosystem.

1.  **Core Package (`managers/`)**: Retain only the interfaces, shared data structures, and common utility functions (`manager.go`). 
2.  **Ecosystem Sub-packages**: Create subdirectories for each language ecosystem (e.g., `managers/node`, `managers/python`, `managers/golang`). 
    *   Each subdirectory will contain the `Language` implementation for that ecosystem.
    *   Each subdirectory will contain the respective `PackageManager` implementations (e.g., `managers/node/npm.go`, `managers/node/bun.go`).
    *   Test files and `testdata` will be moved closer to the code they test within these subdirectories.
    *   The `FetchDetails` and `Uninstall` methods (currently in `fallback_extended.go` and `*_extended.go`) will be merged into their respective manager's primary file.
3.  **Registration**: Update `main.go` to import and instantiate the languages from these new sub-packages.

## Alternatives Considered
*   **Keep current flat structure**: Easiest short-term, but unsustainable long-term. Finding specific logic requires jumping across multiple files (e.g., `cargo.go`, `languages.go`, `fallback_extended.go`).
*   **Split by concern (Languages vs. Managers)**: Having a `managers/langs` folder and a `managers/tools` folder. This separates related concepts; a language and its manager are highly cohesive and should live together.

## Consequences
*   **Pros**: 
    *   Drastically improves maintainability and navigation.
    *   Encapsulates logic by domain.
    *   Simplifies adding new languages (just add a new folder).
*   **Cons**:
    *   Requires a significant initial refactor of file paths and package declarations.
    *   `main.go` will require multiple imports instead of a single `senk-tui/managers` import (though an index/registry file could mitigate this).
