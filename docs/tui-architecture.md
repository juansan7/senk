# TUI Architecture (Bubble Tea & Lipgloss)

`senk-tui` is built using the **Elm Architecture** (Model, Update, View) provided by `charmbracelet/bubbletea`. It uses `charmbracelet/lipgloss` for styling and layout.

This document explains how the code in the `tui/` folder is structured.

## 1. The Elm Architecture

Bubble Tea programs revolve around three core concepts:
1.  **Model**: The state of your application.
2.  **Update**: A function that handles events (key presses, timer ticks, network responses) and updates the Model.
3.  **View**: A function that renders the UI based solely on the current state of the Model.

## 2. Directory Structure (`tui/`)

```text
tui/
├── model.go          # Defines the State (Model) and its initialization.
├── model_modal.go    # Defines the specific state structures for the Modal overlay.
├── messages.go       # Defines custom events (Msgs) used for async communication.
├── messages_modal.go # Defines custom events (Msgs) for the Modal (FetchDetails, Uninstall).
├── update.go         # The main Update() loop and the tea.Cmd generators (side-effects).
├── view.go           # The main View() rendering logic and layout generation.
└── styles.go         # Lipgloss style definitions (colors, margins, borders).
```

## 3. The Model (`tui/model.go`)

The state is strictly hierarchical to support the 3-pane layout:

```go
type Model struct {
    // 1. Data Hierarchy
    Languages   []*LanguageData // Contains Languages -> Managers -> Dependencies
    
    // 2. Navigation State
    LangIndex   int             // Which Language is selected (Left Pane)
    MgrIndex    int             // Which Manager is selected (Middle Pane)
    DepIndex    int             // Which Package is selected (Right Pane)
    
    // 3. UI Context
    Focus       FocusState      // Where is the user's cursor? (Left, Middle, or Right)
    Width       int             // Terminal width
    Height      int             // Terminal height
    Spinner     spinner.Model   // Loading animation state
    
    // 4. Modal State
    Modal       ModalData       // Is the modal open? What package is it showing?
}
```

## 4. The Update Loop (`tui/update.go`)

The `Update(msg tea.Msg) (tea.Model, tea.Cmd)` function receives events. 
We handle three main types of messages:

1.  **`tea.KeyMsg` (Keyboard)**: Updates `Focus`, `LangIndex`, `MgrIndex`, etc. If `Enter` is pressed on a package, it changes the `Modal.State` and returns a command to fetch details.
2.  **Custom Data Messages (e.g., `DepsFetchedMsg`)**: These are returned by background goroutines when a shell command finishes. The `Update` function takes this data, places it into the correct spot in `Model.Languages`, and changes the loading state to `StateDone`.
3.  **`tea.WindowSizeMsg`**: Fired when the terminal is resized. We save the new `Width` and `Height`.

### Asynchronous Commands (`tea.Cmd`)
We *never* run `exec.Command` directly inside `Update`. Instead, we return a `tea.Cmd`. A `tea.Cmd` is simply a function that runs in a Goroutine and returns a `tea.Msg` when it finishes. 
Example: `fetchMgrCmd` runs `manager.Fetch()` in the background and returns a `DepsFetchedMsg`.

## 5. The View (`tui/view.go`)

The `View() string` function takes the `Model` and turns it into a string of ANSI characters using `lipgloss`.

1.  **Layout Calculation**: It calculates the widths (`paneWidth1 = 25%`, `paneWidth3 = 50%`).
2.  **Style Application**: It checks `m.Focus`. If `m.Focus == FocusLeft`, it applies `activeStyle` (purple border) to the left pane; otherwise, it applies `baseStyle` (grey border).
3.  **Content Rendering**: It calls `viewLeftPane()`, `viewMiddlePane()`, and `viewRightPane()`.
4.  **Joining**: It joins the three columns horizontally using `lipgloss.JoinHorizontal`.
5.  **Modal Overlay**: Finally, if `m.Modal.State != ModalClosed`, it takes the entire 3-pane UI string, renders the modal box, and uses `lipgloss.Place` to draw the modal exactly in the center over the top of the background.
