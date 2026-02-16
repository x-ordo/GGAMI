package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	gormtmpl "ggami/internal/templates/gorm"
)

// GormCodeGenerator generates multi-file GORM-based Go projects
type GormCodeGenerator struct{}

// Generate creates a full GORM project from config
func (g *GormCodeGenerator) Generate(config ProjectConfig) error {
	data := buildTemplateData(config)

	// Fix port from int
	if config.Port > 0 {
		data.Port = fmt.Sprintf("%d", config.Port)
	} else {
		data.Port = "8080"
	}

	// Single Go files (use standard {{ }} delimiters)
	if err := g.renderGoFile(config.TargetPath, "main.go", "main.go.tmpl", data); err != nil {
		return fmt.Errorf("main.go: %w", err)
	}
	if err := g.renderGoFile(config.TargetPath, "go.mod", "go_mod.go.tmpl", data); err != nil {
		return fmt.Errorf("go.mod: %w", err)
	}
	// Helpers
	if err := g.renderGoFile(filepath.Join(config.TargetPath, "handlers"), "helpers.go", "helpers.go.tmpl", data); err != nil {
		return fmt.Errorf("helpers.go: %w", err)
	}
	// HTML templates (use << >> delimiters so {{ }} passes through to output)
	if err := g.renderHTMLFile(filepath.Join(config.TargetPath, "templates"), "layout.html", "layout.html.tmpl", data); err != nil {
		return fmt.Errorf("layout.html: %w", err)
	}

	// Per-model files
	for _, model := range data.Models {
		modelData := struct {
			TemplateData
			Model ModelTmplData
		}{data, model}

		modelFile := strings.ToLower(model.Name) + ".go"
		if err := g.renderGoFile(filepath.Join(config.TargetPath, "models"), modelFile, "model.go.tmpl", modelData); err != nil {
			return fmt.Errorf("model %s: %w", model.Name, err)
		}

		handlerFile := strings.ToLower(model.Name) + ".go"
		if err := g.renderGoFile(filepath.Join(config.TargetPath, "handlers"), handlerFile, "handler.go.tmpl", modelData); err != nil {
			return fmt.Errorf("handler %s: %w", model.Name, err)
		}

		listFile := strings.ToLower(model.Name) + "_list.html"
		if err := g.renderHTMLFile(filepath.Join(config.TargetPath, "templates"), listFile, "list.html.tmpl", modelData); err != nil {
			return fmt.Errorf("list template %s: %w", model.Name, err)
		}

		formFile := strings.ToLower(model.Name) + "_form.html"
		if err := g.renderHTMLFile(filepath.Join(config.TargetPath, "templates"), formFile, "form.html.tmpl", modelData); err != nil {
			return fmt.Errorf("form template %s: %w", model.Name, err)
		}
	}

	// RBAC templates (Phase 2)
	if data.HasRBAC {
		// User model for auth
		if err := g.renderGoFile(filepath.Join(config.TargetPath, "models"), "user.go", "user_model.go.tmpl", data); err != nil {
			return fmt.Errorf("user model: %w", err)
		}
		if err := g.renderGoFile(filepath.Join(config.TargetPath, "middleware"), "auth.go", "middleware_auth.go.tmpl", data); err != nil {
			return fmt.Errorf("middleware auth: %w", err)
		}
		if err := g.renderGoFile(filepath.Join(config.TargetPath, "middleware"), "rbac.go", "middleware_rbac.go.tmpl", data); err != nil {
			return fmt.Errorf("middleware rbac: %w", err)
		}
		if err := g.renderGoFile(filepath.Join(config.TargetPath, "handlers"), "auth.go", "auth_handler.go.tmpl", data); err != nil {
			return fmt.Errorf("auth handler: %w", err)
		}
		if err := g.renderHTMLFile(filepath.Join(config.TargetPath, "templates"), "login.html", "login.html.tmpl", data); err != nil {
			return fmt.Errorf("login template: %w", err)
		}
	}

	return nil
}

// renderGoFile renders a Go source template using standard {{ }} delimiters
func (g *GormCodeGenerator) renderGoFile(dir, filename, tmplName string, data interface{}) error {
	content, err := gormtmpl.FS.ReadFile(tmplName)
	if err != nil {
		return fmt.Errorf("read template %s: %w", tmplName, err)
	}

	tmpl, err := template.New(tmplName).Funcs(GormFuncMap).Parse(string(content))
	if err != nil {
		return fmt.Errorf("parse template %s: %w", tmplName, err)
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("execute template %s: %w", tmplName, err)
	}

	return os.WriteFile(filepath.Join(dir, filename), []byte(buf.String()), 0644)
}

// renderHTMLFile renders an HTML template using << >> delimiters
// so that {{ }} in the output is preserved for the generated project's html/template
func (g *GormCodeGenerator) renderHTMLFile(dir, filename, tmplName string, data interface{}) error {
	content, err := gormtmpl.FS.ReadFile(tmplName)
	if err != nil {
		return fmt.Errorf("read template %s: %w", tmplName, err)
	}

	tmpl, err := template.New(tmplName).Delims("<<", ">>").Funcs(GormFuncMap).Parse(string(content))
	if err != nil {
		return fmt.Errorf("parse template %s: %w", tmplName, err)
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("execute template %s: %w", tmplName, err)
	}

	return os.WriteFile(filepath.Join(dir, filename), []byte(buf.String()), 0644)
}
