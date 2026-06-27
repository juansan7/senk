package managers

// Dependency represents a single installed package.
type Dependency struct {
	Name    string
	Version string
}

// PackageManager is the interface all ecosystem scanners must implement.
type PackageManager interface {
	Name() string
	ManagerVersion() (string, error)
	Fetch() ([]Dependency, error)
	IsInstalled() bool
}

// Language represents a programming language/runtime that groups package managers.
type Language interface {
	Name() string
	LanguageVersion() (string, error)
	Managers() []PackageManager
	IsInstalled() bool
}
