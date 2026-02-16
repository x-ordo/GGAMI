package builder

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Page represents a single page in the builder project
type Page struct {
	ID         string      `json:"id"`
	Name       string      `json:"name"`
	Components []Component `json:"components"`
}

// BuilderProject represents a visual builder project
type BuilderProject struct {
	Name      string `json:"name"`
	Pages     []Page `json:"pages"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// ProjectManager manages builder projects
type ProjectManager struct {
	current *BuilderProject
}

// NewProjectManager creates a new ProjectManager
func NewProjectManager() *ProjectManager {
	return &ProjectManager{}
}

// GetCurrentProject returns the current project
func (pm *ProjectManager) GetCurrentProject() *BuilderProject {
	return pm.current
}

// CreateProject creates a new builder project
func (pm *ProjectManager) CreateProject(name string) (*BuilderProject, error) {
	now := time.Now().Format(time.RFC3339)
	pm.current = &BuilderProject{
		Name: name,
		Pages: []Page{
			{
				ID:         "page-1",
				Name:       "Home",
				Components: []Component{},
			},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
	return pm.current, nil
}

// SaveProject saves the project to a JSON file
func (pm *ProjectManager) SaveProject(path string) error {
	if pm.current == nil {
		return fmt.Errorf("no project to save")
	}
	pm.current.UpdatedAt = time.Now().Format(time.RFC3339)
	data, err := json.MarshalIndent(pm.current, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// LoadProject loads a project from a JSON file
func (pm *ProjectManager) LoadProject(path string) (*BuilderProject, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read project file: %w", err)
	}
	var project BuilderProject
	if err := json.Unmarshal(data, &project); err != nil {
		return nil, fmt.Errorf("failed to parse project: %w", err)
	}
	pm.current = &project
	return pm.current, nil
}

// AddComponent adds a component to a page
func (pm *ProjectManager) AddComponent(pageID string, comp Component) error {
	page, err := pm.findPage(pageID)
	if err != nil {
		return err
	}
	page.Components = append(page.Components, comp)
	return nil
}

// UpdateComponent updates a component on a page
func (pm *ProjectManager) UpdateComponent(pageID, compID string, updates Component) error {
	page, err := pm.findPage(pageID)
	if err != nil {
		return err
	}
	for i, c := range page.Components {
		if c.ID == compID {
			if updates.Content != "" {
				page.Components[i].Content = updates.Content
			}
			if updates.Styles != nil {
				page.Components[i].Styles = updates.Styles
			}
			if updates.Type != "" {
				page.Components[i].Type = updates.Type
			}
			return nil
		}
	}
	return fmt.Errorf("component %s not found", compID)
}

// DeleteComponent deletes a component from a page
func (pm *ProjectManager) DeleteComponent(pageID, compID string) error {
	page, err := pm.findPage(pageID)
	if err != nil {
		return err
	}
	for i, c := range page.Components {
		if c.ID == compID {
			page.Components = append(page.Components[:i], page.Components[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("component %s not found", compID)
}

// ReorderComponents reorders components on a page
func (pm *ProjectManager) ReorderComponents(pageID string, orderedIDs []string) error {
	page, err := pm.findPage(pageID)
	if err != nil {
		return err
	}

	compMap := make(map[string]Component)
	for _, c := range page.Components {
		compMap[c.ID] = c
	}

	reordered := make([]Component, 0, len(orderedIDs))
	for _, id := range orderedIDs {
		if comp, ok := compMap[id]; ok {
			reordered = append(reordered, comp)
		}
	}
	page.Components = reordered
	return nil
}

// ExportToHTML exports the project as static HTML files
func (pm *ProjectManager) ExportToHTML(outputPath string) error {
	if pm.current == nil {
		return fmt.Errorf("no project to export")
	}
	exporter := &Exporter{project: pm.current}
	return exporter.Export(outputPath)
}

func (pm *ProjectManager) findPage(pageID string) (*Page, error) {
	if pm.current == nil {
		return nil, fmt.Errorf("no project loaded")
	}
	for i := range pm.current.Pages {
		if pm.current.Pages[i].ID == pageID {
			return &pm.current.Pages[i], nil
		}
	}
	return nil, fmt.Errorf("page %s not found", pageID)
}
