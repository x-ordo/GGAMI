package generator

import (
	"os"
	"path/filepath"
	"strings"

	"ggami-go/internal/modules"
	"ggami-go/internal/templates"
)

// GoGenerator implements the Generator interface for Go projects
type GoGenerator struct{}

func (g *GoGenerator) Scaffold(path string) error {
	dirs := []string{
		path,
		filepath.Join(path, "templates"),
		filepath.Join(path, "assets"),
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	return nil
}

// scaffoldGorm creates additional directories for GORM projects
func (g *GoGenerator) scaffoldGorm(path string) error {
	dirs := []string{
		path,
		filepath.Join(path, "models"),
		filepath.Join(path, "handlers"),
		filepath.Join(path, "middleware"),
		filepath.Join(path, "templates"),
		filepath.Join(path, "assets"),
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	return nil
}

func (g *GoGenerator) CreateManifest(config ProjectConfig) error {
	if config.GormMode && len(config.Models) > 0 {
		return nil // GormCodeGenerator handles go.mod via template
	}
	goMod := strings.Replace(templates.GoModTemplate(), "{{PROJECT_NAME}}", config.ProjectName, 1)
	return os.WriteFile(filepath.Join(config.TargetPath, "go.mod"), []byte(goMod), 0644)
}

func (g *GoGenerator) GenerateCode(config ProjectConfig) error {
	if config.GormMode && len(config.Models) > 0 {
		if err := g.scaffoldGorm(config.TargetPath); err != nil {
			return err
		}
		return (&GormCodeGenerator{}).Generate(config)
	}

	// Legacy: Replace template variables
	mainGo := templates.GoMainTemplate()
	mainGo = strings.Replace(mainGo, "{{DB_SERVER}}", config.DBServer, 1)
	mainGo = strings.Replace(mainGo, "{{DB_USER}}", config.DBUser, 1)
	mainGo = strings.Replace(mainGo, "{{DB_PW}}", config.DBPw, 1)
	mainGo = strings.Replace(mainGo, "{{DB_NAME}}", config.DBName, 1)

	indexHTML := templates.HTMLIndexTemplate()
	indexHTML = strings.ReplaceAll(indexHTML, "{{PROJECT_NAME}}", config.ProjectName)

	// Read go.mod for potential module injection
	goModPath := filepath.Join(config.TargetPath, "go.mod")
	goModBytes, err := os.ReadFile(goModPath)
	if err != nil {
		return err
	}
	goMod := string(goModBytes)

	// Module injection
	activeModules := filterActiveModules(config.Modules)

	for _, mod := range activeModules {
		for _, snippet := range mod.Snippets {
			switch snippet.Target {
			case modules.TargetMainGo:
				mainGo = strings.Replace(mainGo, snippet.Marker, snippet.Content+"\n"+snippet.Marker, 1)
			case modules.TargetIndexHTML:
				indexHTML = strings.Replace(indexHTML, snippet.Marker, snippet.Content+"\n"+snippet.Marker, 1)
			case modules.TargetGoMod:
				goMod = strings.Replace(goMod, snippet.Marker, snippet.Content+"\n"+snippet.Marker, 1)
			}
		}
	}

	// Write files
	if err := os.WriteFile(filepath.Join(config.TargetPath, "main.go"), []byte(mainGo), 0644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(config.TargetPath, "templates", "index.html"), []byte(indexHTML), 0644); err != nil {
		return err
	}
	if err := os.WriteFile(goModPath, []byte(goMod), 0644); err != nil {
		return err
	}

	return nil
}

// ScaffoldGorm creates the directory structure for a GORM project
func ScaffoldGorm(path string) error {
	return (&GoGenerator{}).scaffoldGorm(path)
}

// GenerateLegacyCode generates legacy (non-GORM) Go project files with module injection.
// This is the code generation portion only — scaffold and manifest must be done separately.
func GenerateLegacyCode(config ProjectConfig) error {
	return (&GoGenerator{}).GenerateCode(config)
}

func filterActiveModules(selectedIDs []string) []modules.ModuleDef {
	byID := make(map[string]modules.ModuleDef)
	for _, mod := range modules.Registry {
		byID[mod.ID] = mod
	}

	var active []modules.ModuleDef
	for _, id := range selectedIDs {
		if mod, ok := byID[id]; ok {
			active = append(active, mod)
		}
	}
	return active
}
