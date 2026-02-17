package application

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"ggami/internal/domain"
	"ggami/internal/generator"
	"ggami/internal/modules"
)

// --- Step 1: ValidateConfigStep ---

type ValidateConfigStep struct{}

func (s *ValidateConfigStep) Name() string { return "ValidateConfig" }

var validNameRe = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]*$`)

// Go reserved words that cannot be used as model names
var goReservedWords = map[string]bool{
	"break": true, "default": true, "func": true, "interface": true, "select": true,
	"case": true, "defer": true, "go": true, "map": true, "struct": true,
	"chan": true, "else": true, "goto": true, "package": true, "switch": true,
	"const": true, "fallthrough": true, "if": true, "range": true, "type": true,
	"continue": true, "for": true, "import": true, "return": true, "var": true,
}

func (s *ValidateConfigStep) Execute(ctx *domain.PipelineContext) error {
	c := ctx.Config

	if c.ProjectName == "" {
		return fmt.Errorf("project name is required")
	}
	if !validNameRe.MatchString(c.ProjectName) {
		return fmt.Errorf("project name %q contains invalid characters (use letters, digits, hyphens, underscores)", c.ProjectName)
	}
	if c.TargetPath == "" {
		return fmt.Errorf("target path is required")
	}
	if ctx.Language != "go" && ctx.Language != "node" {
		return fmt.Errorf("unsupported language %q (use \"go\" or \"node\")", ctx.Language)
	}

	if c.GormMode && len(c.Models) > 0 {
		modelNames := make(map[string]bool)
		for _, m := range c.Models {
			if m.Name == "" {
				return fmt.Errorf("model name cannot be empty")
			}
			lower := strings.ToLower(m.Name)
			if goReservedWords[lower] {
				return fmt.Errorf("model name %q is a Go reserved word", m.Name)
			}
			if modelNames[lower] {
				return fmt.Errorf("duplicate model name %q", m.Name)
			}
			modelNames[lower] = true

			if len(m.Fields) == 0 {
				return fmt.Errorf("model %q must have at least one field", m.Name)
			}
			for _, f := range m.Fields {
				if f.Name == "" {
					return fmt.Errorf("field name cannot be empty in model %q", m.Name)
				}
				if !isValidFieldType(f.Type) {
					return fmt.Errorf("invalid field type %q for field %q in model %q", f.Type, f.Name, m.Name)
				}
			}
		}

		if c.DBType != "" && c.DBType != domain.DBTypeMSSQL && c.DBType != domain.DBTypePostgres &&
			c.DBType != domain.DBTypeMySQL && c.DBType != domain.DBTypeSQLite {
			return fmt.Errorf("unsupported database type %q", c.DBType)
		}
	}

	if c.RBAC != nil && c.RBAC.Enabled {
		if len(c.RBAC.Roles) == 0 {
			return fmt.Errorf("RBAC enabled but no roles defined")
		}
	}

	return nil
}

func isValidFieldType(t string) bool {
	switch t {
	case "string", "int", "uint", "float64", "bool", "time.Time":
		return true
	}
	return false
}

func (s *ValidateConfigStep) Rollback(ctx *domain.PipelineContext) error { return nil }

// --- Step 2: ResolveModulesStep ---

type ResolveModulesStep struct{}

func (s *ResolveModulesStep) Name() string { return "ResolveModules" }

func (s *ResolveModulesStep) Execute(ctx *domain.PipelineContext) error {
	if len(ctx.Config.Modules) == 0 {
		ctx.Modules = nil
		return nil
	}
	resolved, err := ResolveDependencies(ctx.Config.Modules, modules.Registry)
	if err != nil {
		return err
	}
	ctx.Modules = resolved
	return nil
}

func (s *ResolveModulesStep) Rollback(ctx *domain.PipelineContext) error { return nil }

// --- Step 3: CreateTempDirStep ---

type CreateTempDirStep struct{}

func (s *CreateTempDirStep) Name() string { return "CreateTempDir" }

func (s *CreateTempDirStep) Execute(ctx *domain.PipelineContext) error {
	dir, err := os.MkdirTemp("", "ggami-gen-*")
	if err != nil {
		return fmt.Errorf("create temp dir: %w", err)
	}
	ctx.TempDir = dir
	return nil
}

func (s *CreateTempDirStep) Rollback(ctx *domain.PipelineContext) error {
	if ctx.TempDir != "" {
		return os.RemoveAll(ctx.TempDir)
	}
	return nil
}

// --- Step 4: ScaffoldStep (legacy mode) ---

type ScaffoldStep struct{}

func (s *ScaffoldStep) Name() string { return "Scaffold" }

func (s *ScaffoldStep) Execute(ctx *domain.PipelineContext) error {
	gen, err := generator.NewGenerator(ctx.Language)
	if err != nil {
		return err
	}
	return gen.Scaffold(ctx.TempDir)
}

func (s *ScaffoldStep) Rollback(ctx *domain.PipelineContext) error { return nil }

// --- Step 4b: ScaffoldGormStep (GORM mode) ---

type ScaffoldGormStep struct{}

func (s *ScaffoldGormStep) Name() string { return "ScaffoldGorm" }

func (s *ScaffoldGormStep) Execute(ctx *domain.PipelineContext) error {
	return generator.ScaffoldGorm(ctx.TempDir)
}

func (s *ScaffoldGormStep) Rollback(ctx *domain.PipelineContext) error { return nil }

// --- Step 5: GenerateCoreStep ---

type GenerateCoreStep struct{}

func (s *GenerateCoreStep) Name() string { return "GenerateCore" }

func (s *GenerateCoreStep) Execute(ctx *domain.PipelineContext) error {
	gen, err := generator.NewGenerator(ctx.Language)
	if err != nil {
		return err
	}
	// Create a copy of config pointing to TempDir
	cfg := ctx.Config
	cfg.TargetPath = ctx.TempDir
	return gen.CreateManifest(cfg)
}

func (s *GenerateCoreStep) Rollback(ctx *domain.PipelineContext) error { return nil }

// --- Step 6: GenerateModelsStep (GORM mode) ---

type GenerateModelsStep struct{}

func (s *GenerateModelsStep) Name() string { return "GenerateModels" }

func (s *GenerateModelsStep) Execute(ctx *domain.PipelineContext) error {
	cfg := ctx.Config
	cfg.TargetPath = ctx.TempDir
	return generator.RenderModels(cfg)
}

func (s *GenerateModelsStep) Rollback(ctx *domain.PipelineContext) error { return nil }

// --- Step 7: GenerateHandlersStep (GORM mode) ---

type GenerateHandlersStep struct{}

func (s *GenerateHandlersStep) Name() string { return "GenerateHandlers" }

func (s *GenerateHandlersStep) Execute(ctx *domain.PipelineContext) error {
	cfg := ctx.Config
	cfg.TargetPath = ctx.TempDir
	return generator.RenderHandlers(cfg)
}

func (s *GenerateHandlersStep) Rollback(ctx *domain.PipelineContext) error { return nil }

// --- Step 8: GenerateTemplatesStep (GORM mode) ---

type GenerateTemplatesStep struct{}

func (s *GenerateTemplatesStep) Name() string { return "GenerateTemplates" }

func (s *GenerateTemplatesStep) Execute(ctx *domain.PipelineContext) error {
	cfg := ctx.Config
	cfg.TargetPath = ctx.TempDir
	return generator.RenderHTMLTemplates(cfg)
}

func (s *GenerateTemplatesStep) Rollback(ctx *domain.PipelineContext) error { return nil }

// --- Step 9: GenerateMiddlewareStep (GORM + RBAC) ---

type GenerateMiddlewareStep struct{}

func (s *GenerateMiddlewareStep) Name() string { return "GenerateMiddleware" }

func (s *GenerateMiddlewareStep) Execute(ctx *domain.PipelineContext) error {
	cfg := ctx.Config
	cfg.TargetPath = ctx.TempDir
	return generator.RenderMiddleware(cfg)
}

func (s *GenerateMiddlewareStep) Rollback(ctx *domain.PipelineContext) error { return nil }

// --- Step 10: InjectModulesStep (legacy mode) ---

type InjectModulesStep struct{}

func (s *InjectModulesStep) Name() string { return "InjectModules" }

func (s *InjectModulesStep) Execute(ctx *domain.PipelineContext) error {
	cfg := ctx.Config
	cfg.TargetPath = ctx.TempDir
	cfg.Modules = moduleIDs(ctx.Modules)
	return generator.GenerateLegacyCode(cfg)
}

func moduleIDs(mods []domain.ModuleDef) []string {
	ids := make([]string, len(mods))
	for i, m := range mods {
		ids[i] = m.ID
	}
	return ids
}

func (s *InjectModulesStep) Rollback(ctx *domain.PipelineContext) error { return nil }

// --- Step 11: FinalizeStep ---

type FinalizeStep struct{}

func (s *FinalizeStep) Name() string { return "Finalize" }

func (s *FinalizeStep) Execute(ctx *domain.PipelineContext) error {
	// Remove existing target if present
	if _, err := os.Stat(ctx.FinalDir); err == nil {
		if err := os.RemoveAll(ctx.FinalDir); err != nil {
			return fmt.Errorf("remove existing target: %w", err)
		}
	}

	// Ensure parent directory exists
	if err := os.MkdirAll(filepath.Dir(ctx.FinalDir), 0755); err != nil {
		return fmt.Errorf("create parent dir: %w", err)
	}

	// Move temp to final
	if err := os.Rename(ctx.TempDir, ctx.FinalDir); err != nil {
		// Rename may fail across filesystems; fall back to copy
		if err := copyDir(ctx.TempDir, ctx.FinalDir); err != nil {
			return fmt.Errorf("move to target: %w", err)
		}
		os.RemoveAll(ctx.TempDir)
	}
	ctx.TempDir = "" // prevent rollback from cleaning up
	return nil
}

func (s *FinalizeStep) Rollback(ctx *domain.PipelineContext) error { return nil }

// copyDir recursively copies a directory tree
func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)

		if info.IsDir() {
			return os.MkdirAll(target, info.Mode())
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, data, info.Mode())
	})
}
