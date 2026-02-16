package generator

import (
	"os"
	"path/filepath"
	"strings"

	"ggami/internal/modules"
	"ggami/internal/templates"
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

func (g *GoGenerator) CreateManifest(config ProjectConfig) error {
	goMod := strings.Replace(templates.GoModTemplate(), "{{PROJECT_NAME}}", config.ProjectName, 1)
	return os.WriteFile(filepath.Join(config.TargetPath, "go.mod"), []byte(goMod), 0644)
}

func (g *GoGenerator) GenerateCode(config ProjectConfig) error {
	// Replace template variables
	mainGo := templates.GoMainTemplate()
	mainGo = strings.Replace(mainGo, "{{DB_SERVER}}", config.DBServer, 1)
	mainGo = strings.Replace(mainGo, "{{DB_USER}}", config.DBUser, 1)
	mainGo = strings.Replace(mainGo, "{{DB_PW}}", config.DBPw, 1)
	mainGo = strings.Replace(mainGo, "{{DB_NAME}}", config.DBName, 1)

	indexHTML := templates.HTMLIndexTemplate()
	indexHTML = strings.ReplaceAll(indexHTML, "{{PROJECT_NAME}}", config.ProjectName)

	// Module injection
	activeModules := filterActiveModules(config.Modules)

	for _, mod := range activeModules {
		for _, snippet := range mod.Snippets {
			switch snippet.Target {
			case modules.TargetMainGo:
				mainGo = strings.Replace(mainGo, snippet.Marker, snippet.Content+"\n"+snippet.Marker, 1)
			case modules.TargetIndexHTML:
				indexHTML = strings.Replace(indexHTML, snippet.Marker, snippet.Content+"\n"+snippet.Marker, 1)
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

	return nil
}

func filterActiveModules(selectedIDs []string) []modules.ModuleDef {
	idSet := make(map[string]bool)
	for _, id := range selectedIDs {
		idSet[id] = true
	}

	var active []modules.ModuleDef
	for _, mod := range modules.Registry {
		if idSet[mod.ID] {
			active = append(active, mod)
		}
	}
	return active
}
