# ADR 0003: Hierarchical 3-Pane Layout (Languages -> Managers -> Packages)

## 1. Context and Problem Statement
The user pointed out a crucial architectural constraint: a single programming language often has multiple package managers (e.g., Node.js has `npm`, `yarn`, `pnpm`; Python has `pip`, `poetry`, `pipx`). The current flat list of managers does not represent this hierarchy well. The user wants to clearly see Languages, Managers, and Packages, each with their respective versions.

## 2. Decision
We will refactor the data model and the TUI into a strict hierarchical system with a **3-pane layout**:

1. **Data Model Refactor**:
   * Introduce a `Language` concept that groups multiple `PackageManager`s.
   * `Language` handles fetching its own version (e.g., `node -v`).
   * `PackageManager` handles fetching its version (e.g., `npm -v`) and its dependencies.

2. **UI Refactor (3 Panes)**:
   * **Pane 1 (Left - 25%)**: List of installed Languages + Versions.
   * **Pane 2 (Middle - 25%)**: List of Package Managers belonging to the selected Language + Versions.
   * **Pane 3 (Right - 50%)**: List of Packages belonging to the selected Manager + Versions.

## 3. Consequences
* **Positive**: Perfectly scales to support `yarn`, `pnpm`, `poetry`, etc., in the future without cluttering the UI. Provides a very logical, structured flow for polyglot developers.
* **Negative**: Requires a significant refactor of the Bubble Tea model (handling focus across 3 states instead of 2) and grouping our existing 11 managers under their parent languages.
