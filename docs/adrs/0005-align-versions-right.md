# 0005: Align Versions Right in Languages and Managers Panes

## Status
Proposed

## Context
Currently, the versions for programming languages and package managers are displayed directly next to their names in parentheses (e.g., `Node (v18.0.0)`). 
The user requested that the versions for languages and managers be aligned to the right edge of their respective panels, matching the styling already used in the "Packages" pane (which uses dot padding `...` to align versions to the right).

## Decision
We will modify the layout of the "Languages" and "Managers" panes to align versions to the right. We will pass the pane width down to the `viewLeftPane` and `viewMiddlePane` rendering functions. We will then calculate the width of the name and version strings, and fill the remaining space with dot padding (`.`), so the version appears at the far right. For the managers, the status indicator (e.g., the spinner or checkmark) will be appended immediately after the version.

## Alternatives Considered
*   **Keep current layout**: Leave the version in parentheses next to the name. This is easier to implement but creates an inconsistent UI experience compared to the Packages pane.
*   **Use flexbox/grid equivalent layouts**: Refactor the TUI framework layout entirely to support flex spaces. This is over-engineering for a simple padding string and introduces unnecessary complexity to the `lipgloss` rendering logic.

## Consequences
*   **Pros**: Creates a uniform look across all three panes in the TUI, improving visual hierarchy and readability.
*   **Cons**: Requires passing pane widths into the rendering functions and calculating string widths manually, adding slight complexity to `tui/view.go`.
