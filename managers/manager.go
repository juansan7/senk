package managers

// Dependency represents a single installed package.
type Dependency struct {
	Name    string
	Version string
}

// PackageDetails contains extended metadata for a modal view.
type PackageDetails struct {
	Name        string
	Version     string
	Description string
	Author      string
	Homepage    string
	Path        string
	Size        string
}

// PackageManager is the interface all ecosystem scanners must implement.
type PackageManager interface {
	Name() string
	ManagerVersion() (string, error)
	Fetch() ([]Dependency, error)
	IsInstalled() bool

	// FetchDetails retrieves deep metadata for a single package.
	FetchDetails(dep Dependency) (PackageDetails, error)
	// Uninstall removes the package from the system.
	Uninstall(dep Dependency) error
}

// Language represents a programming language/runtime that groups package managers.
type Language interface {
	Name() string
	LanguageVersion() (string, error)
	Managers() []PackageManager
	IsInstalled() bool
}
