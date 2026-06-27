# AI Agent Rules & Guidelines (Codename: senk-tui)

When working on this codebase, the AI agent must strictly adhere to the following general rules to ensure high-quality, maintainable, and verifiable code.

## 1. Development Methodology (TDD)
*   **Test-First Approach:** Always write tests before implementing the actual logic. Follow the Red-Green-Refactor cycle.
*   **Table-Driven Tests:** Use idiomatic table-driven test structures (e.g., in Go) for unit tests.
*   **Testdata Directory:** For parsing complex outputs or files, save raw mock outputs in a `testdata/` directory to ensure tests run against real-world string formats.

## 2. Architecture & Separation of Concerns
*   **Pure Functions:** Never mix side-effects (I/O, database calls, shell command execution) with data processing or parsing logic. 
    *   *Good:* A function that takes a raw byte array (the data) and returns parsed structures (Pure function, easily testable). The I/O runs in a separate, thin wrapper.
*   **Interface Driven:** Use interfaces to decouple components and make the system easily extensible without modifying core logic.

## 3. Concurrency & Performance
*   **Never Block the Main Thread/Loop:** Whether building a TUI, GUI, or web app, ensure the UI remains responsive. Heavy I/O or shell commands must be executed asynchronously in the background.
*   **Message Passing:** Communicate between background processes and the main application state using established concurrency patterns (e.g., events, messages, channels).

## 4. Code Quality & Conventions
*   **Formatting:** Follow the standard formatting tools for the chosen language (e.g., `gofmt` for Go).
*   **Error Handling:** Handle all errors explicitly. Do not silently ignore errors.
*   **Self-Contained:** Rely on standard library packages as much as possible to minimize third-party dependencies, unless a specific framework is approved.

## 5. Version Control: Atomic Commits & PR Limits
*   **Atomic Commits:** Commits must do exactly *one* specific thing. They must have a clear, descriptive commit message and include a minimal number of modified files. Do not mix refactoring with new feature additions in the same commit.
*   **Pull Request Size:** PRs should contain a logical grouping of changes but must remain small and reviewable. A PR **must not exceed 12 to 15 changed files** (including tests). If a feature requires more files, break it down into smaller, sequential PRs.

## 6. Scope Constraints
*   **Strict Scope:** Do not implement features or side-effects outside the explicitly defined scope (PRD or Technical Specs) without user approval.

## 7. Project Management & Documentation
*   **Task Management:** All tasks and progress must be added, tracked, and managed strictly within the `TASKS.md` file.
*   **Documentation Location:** All project documentation must reside strictly within the `docs/` directory.
*   **Post-MVP Architecture Decisions (ADRs):** After the MVP phase, any significant architectural or feature changes must first be documented in the `docs/adrs/` folder following a numbered sequence (e.g., `0001-feature-name.md`). The ADR must include an "Alternatives Considered" section. **Execution of the code changes can only proceed after the ADR is approved.**