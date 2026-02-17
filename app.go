package main

import (
	"context"
	"fmt"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"ggami-go/internal/application"
	"ggami-go/internal/builder"
	"ggami-go/internal/generator"
	"ggami-go/internal/modules"
)

// App struct
type App struct {
	ctx context.Context
	pm  *builder.ProjectManager
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		pm: builder.NewProjectManager(),
	}
}

// startup is called when the app starts
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// SelectDirectory opens a native directory picker dialog
func (a *App) SelectDirectory() (string, error) {
	dir, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title:            "Select Target Directory",
		DefaultDirectory: "C:\\",
	})
	if err != nil {
		return "", err
	}
	return dir, nil
}

// GetModules returns all available modules
func (a *App) GetModules() []modules.ModuleDef {
	return modules.Registry
}

// GenerateProject generates a project with the given config and language
func (a *App) GenerateProject(config generator.ProjectConfig, lang string) map[string]interface{} {
	result, err := application.GenerateProject(config, lang)
	if err != nil {
		return map[string]interface{}{
			"success": false,
			"message": fmt.Sprintf("Generation failed: %v", err),
		}
	}
	return map[string]interface{}{
		"success": true,
		"message": result,
	}
}

// --- Builder Methods ---

// CreateBuilderProject creates a new builder project
func (a *App) CreateBuilderProject(name string) (*builder.BuilderProject, error) {
	return a.pm.CreateProject(name)
}

// SaveBuilderProject saves the current project to a file
func (a *App) SaveBuilderProject() error {
	path, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:           "Save Builder Project",
		DefaultFilename: "project.ggami.json",
	})
	if err != nil {
		return err
	}
	if path == "" {
		return fmt.Errorf("no file selected")
	}
	return a.pm.SaveProject(path)
}

// LoadBuilderProject loads a project from a file
func (a *App) LoadBuilderProject() (*builder.BuilderProject, error) {
	path, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Open Builder Project",
		Filters: []runtime.FileFilter{
			{DisplayName: "Ggami Projects", Pattern: "*.ggami.json"},
		},
	})
	if err != nil {
		return nil, err
	}
	if path == "" {
		return nil, fmt.Errorf("no file selected")
	}
	return a.pm.LoadProject(path)
}

// ExportBuilderHTML exports the current project as static HTML
func (a *App) ExportBuilderHTML() error {
	dir, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Export Directory",
	})
	if err != nil {
		return err
	}
	if dir == "" {
		return fmt.Errorf("no directory selected")
	}
	return a.pm.ExportToHTML(dir)
}

// AddComponent adds a component to a page
func (a *App) AddComponent(pageID string, comp builder.Component) error {
	return a.pm.AddComponent(pageID, comp)
}

// UpdateComponent updates a component on a page
func (a *App) UpdateComponent(pageID, compID string, updates builder.Component) error {
	return a.pm.UpdateComponent(pageID, compID, updates)
}

// DeleteComponent deletes a component from a page
func (a *App) DeleteComponent(pageID, compID string) error {
	return a.pm.DeleteComponent(pageID, compID)
}

// ReorderComponents reorders components on a page
func (a *App) ReorderComponents(pageID string, orderedIDs []string) error {
	return a.pm.ReorderComponents(pageID, orderedIDs)
}

// GetCurrentProject returns the current builder project
func (a *App) GetCurrentProject() *builder.BuilderProject {
	return a.pm.GetCurrentProject()
}
