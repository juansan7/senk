# ADR 0001: Show Language/Manager Versions

## 1. Context and Problem Statement
Currently, the UI displays the name of the package manager (e.g., "Node.js (npm)", "Homebrew"). However, developers often need to know the specific version of the tool they are using (e.g., "Node.js (npm v10.2.4)") directly in the UI without having to run a separate `--version` command in their terminal.

## 2. Decision
We will extend the `PackageManager` interface to include a `Version()` method. When the application initializes, it will fetch the version of the CLI tool (e.g., `npm --version`, `brew --version`) and append it to the name displayed in the Left Pane of the TUI.

## 3. Alternatives Considered
* **Alternative 1**: Add a separate UI element (like a footer) to show the version of the currently selected manager.
  * *Why rejected*: It's less scannable. Having it directly next to the manager's name in the left list provides immediate context for all installed managers at a glance.
* **Alternative 2**: Fetch the version synchronously in `IsInstalled()`.
  * *Why rejected*: Running `--version` shell commands can take 10-50ms each. Doing this sequentially for 11 managers before starting the UI would block the main thread and violate our `< 100ms` startup goal.

## 4. Consequences
* **Positive**: Better UX and immediate visibility of toolchain versions.
* **Negative**: Requires firing off an additional shell command per ecosystem. To mitigate performance impacts, the version fetching must happen asynchronously along with the dependency fetching.
