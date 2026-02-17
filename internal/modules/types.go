package modules

import "ggami/internal/domain"

// Type aliases for backward compatibility
type InjectionTarget = domain.InjectionTarget
type CodeSnippet = domain.CodeSnippet
type ModuleDef = domain.ModuleDef

// Re-export constants
const (
	TargetMainGo   = domain.TargetMainGo
	TargetGoMod    = domain.TargetGoMod
	TargetIndexHTML = domain.TargetIndexHTML
)
