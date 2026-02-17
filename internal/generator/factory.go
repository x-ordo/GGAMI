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
