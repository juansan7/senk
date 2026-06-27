# ADR 0004: Package Detail Modal and Uninstall Action

## 1. Context and Problem Statement
The user requested the ability to select a package, open a modal with detailed information (size, install date, etc.), and perform an **uninstall** action from within the TUI. 
This requires dropping the initial "Read-Only" constraint. Furthermore, we must determine what useful metadata can realistically be extracted across different package managers to populate this modal.

## 2. Decision: Metadata & Mutability
We will expand the scope of `senk-tui` to be a mutating dashboard.

### 2.1 Package Metadata Extraction Strategy
Different package managers expose different metadata natively. When native data is unavailable, we will fallback to filesystem stats. The modal will display a best-effort combination of:
*   **Description/Summary:** What the package does.
*   **Author/Publisher:** Who made it.
*   **Homepage/URL:** Where to find the source.
*   **Install Size:** Disk space consumed.
*   **Path:** Where the binary/library is located on the disk.

**Ecosystem specifics for `FetchDetails()`:**
*   **Homebrew**: `brew info <pkg> --json` (Provides Description, URL, Size, Install Date).
*   **NPM**: `npm view <pkg> description author homepage` (Needs network, or we inspect `package.json` in global `node_modules`).
*   **Pip**: `pip show <pkg>` (Provides Summary, Author, Location).
*   **Cargo**: Cargo doesn't store metadata locally after install. We will provide the binary path (`~/.cargo/bin/<pkg>`) and file size via `os.Stat()`.
*   **Go**: `go version -m` gives the module path. We can use `os.Stat()` for the binary size and modification date.

### 2.2 TUI Modal Design
*   We will use `lipgloss` to render a centered overlay box.
*   Pressing `Enter` on a package triggers `FetchDetails()`. A spinner shows in the modal until the data arrives.
*   Inside the modal, pressing `x` or `backspace` triggers the `Uninstall()` goroutine.
*   The modal blocks interaction with the background panes until closed with `Esc`.

## 3. Consequences
*   **Positive**: The app transforms from a simple list into a powerful control center. The metadata provides context on forgotten packages.
*   **Negative**: Fetching details requires a separate shell command per package (lazy-loaded on `Enter`). Network calls (like `npm view`) might be slow if the user has a poor connection.
