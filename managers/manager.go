package managers

// Dependency represents a single installed package.
type Dependency struct {
	Name    string
	Version string
}

// PackageManager is the interface all ecosystem scanners must implement.
type PackageManager interface {
	Name() string
	// ManagerVersion fetches the CLI tool version
	ManagerVersion() (string, error)
	// Fetch runs the underlying shell commands and parses the output.
	Fetch() ([]Dependency, error)
	// IsInstalled checks if the underlying CLI tool exists in $PATH.
	IsInstalled() bool
}
