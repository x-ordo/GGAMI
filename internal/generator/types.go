package generator

import "ggami-go/internal/domain"

// Type aliases for backward compatibility
type ProjectConfig = domain.ProjectConfig

// Generator interface defines the contract for code generators
type Generator interface {
	// Scaffold creates the basic folder structure
	Scaffold(path string) error
	// CreateManifest creates config files (go.mod or package.json)
	CreateManifest(config ProjectConfig) error
	// GenerateCode creates source code files
	GenerateCode(config ProjectConfig) error
}
