package generator

import "fmt"

// NewGenerator creates a generator for the specified language
func NewGenerator(lang string) (Generator, error) {
	switch lang {
	case "go":
		return &GoGenerator{}, nil
	case "node":
		return &NodeGenerator{}, nil
	default:
		return nil, fmt.Errorf("unsupported language type: %s", lang)
	}
}

// GenerateProject orchestrates the full project generation
func GenerateProject(config ProjectConfig, language string) (string, error) {
	gen, err := NewGenerator(language)
	if err != nil {
		return "", err
	}

	// 1. Scaffold
	if err := gen.Scaffold(config.TargetPath); err != nil {
		return "", fmt.Errorf("scaffold failed: %w", err)
	}

	// 2. Manifest
	if err := gen.CreateManifest(config); err != nil {
		return "", fmt.Errorf("manifest creation failed: %w", err)
	}

	// 3. Code generation
	if err := gen.GenerateCode(config); err != nil {
		return "", fmt.Errorf("code generation failed: %w", err)
	}

	return "Generation complete: " + config.TargetPath, nil
}
