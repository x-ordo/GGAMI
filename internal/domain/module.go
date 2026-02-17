package domain

// InjectionTarget specifies which file a code snippet targets
type InjectionTarget string

const (
	TargetMainGo   InjectionTarget = "main.go"
	TargetGoMod    InjectionTarget = "go.mod"
	TargetIndexHTML InjectionTarget = "index.html"
)

// CodeSnippet represents a piece of code to inject at a marker
type CodeSnippet struct {
	Target  InjectionTarget `json:"target"`
	Marker  string          `json:"marker"`
	Content string          `json:"content"`
}

// ModuleDef defines a pluggable module with injectable code
type ModuleDef struct {
	ID           string        `json:"id"`
	Name         string        `json:"name"`
	Description  string        `json:"description"`
	Category     string        `json:"category"`     // "feature", "ui", "utils"
	Dependencies []string      `json:"dependencies"` // required module IDs
	Snippets     []CodeSnippet `json:"snippets"`
}
