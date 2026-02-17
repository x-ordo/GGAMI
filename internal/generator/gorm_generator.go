package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	gormtmpl "ggami-go/internal/templates/gorm"
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

	// Base dashboard pages
	if err := g.renderBasePages(config.TargetPath, data); err != nil {
		return fmt.Errorf("base pages: %w", err)
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
		if err := g.renderHTMLFile(filepath.Join(config.TargetPath, "templates"), "register.html", "register.html.tmpl", data); err != nil {
			return fmt.Errorf("register template: %w", err)
		}
		if err := g.renderHTMLFile(filepath.Join(config.TargetPath, "templates"), "forgot_password.html", "forgot_password.html.tmpl", data); err != nil {
			return fmt.Errorf("forgot_password template: %w", err)
		}
	}

	return nil
}

// RenderModels generates GORM model files (models/*.go)
func RenderModels(config ProjectConfig) error {
	g := &GormCodeGenerator{}
	data := buildTemplateData(config)

	for _, model := range data.Models {
		modelData := struct {
			TemplateData
			Model ModelTmplData
		}{data, model}

		modelFile := strings.ToLower(model.Name) + ".go"
		if err := g.renderGoFile(filepath.Join(config.TargetPath, "models"), modelFile, "model.go.tmpl", modelData); err != nil {
			return fmt.Errorf("model %s: %w", model.Name, err)
		}
	}
	return nil
}

// RenderHandlers generates handler files (handlers/*.go + helpers.go)
func RenderHandlers(config ProjectConfig) error {
	g := &GormCodeGenerator{}
	data := buildTemplateData(config)
	if config.Port > 0 {
		data.Port = fmt.Sprintf("%d", config.Port)
	} else {
		data.Port = "8080"
	}

	// main.go and go.mod
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
	// Base handler
	if err := g.renderGoFile(filepath.Join(config.TargetPath, "handlers"), "base.go", "base_handler.go.tmpl", data); err != nil {
		return fmt.Errorf("base handler: %w", err)
	}

	for _, model := range data.Models {
		modelData := struct {
			TemplateData
			Model ModelTmplData
		}{data, model}

		handlerFile := strings.ToLower(model.Name) + ".go"
		if err := g.renderGoFile(filepath.Join(config.TargetPath, "handlers"), handlerFile, "handler.go.tmpl", modelData); err != nil {
			return fmt.Errorf("handler %s: %w", model.Name, err)
		}
	}
	return nil
}

// RenderHTMLTemplates generates HTML template files (templates/*.html)
func RenderHTMLTemplates(config ProjectConfig) error {
	g := &GormCodeGenerator{}
	data := buildTemplateData(config)

	if err := g.renderHTMLFile(filepath.Join(config.TargetPath, "templates"), "layout.html", "layout.html.tmpl", data); err != nil {
		return fmt.Errorf("layout.html: %w", err)
	}

	// Base dashboard pages
	if err := g.renderBasePages(config.TargetPath, data); err != nil {
		return fmt.Errorf("base pages: %w", err)
	}

	for _, model := range data.Models {
		modelData := struct {
			TemplateData
			Model ModelTmplData
		}{data, model}

		listFile := strings.ToLower(model.Name) + "_list.html"
		if err := g.renderHTMLFile(filepath.Join(config.TargetPath, "templates"), listFile, "list.html.tmpl", modelData); err != nil {
			return fmt.Errorf("list template %s: %w", model.Name, err)
		}

		formFile := strings.ToLower(model.Name) + "_form.html"
		if err := g.renderHTMLFile(filepath.Join(config.TargetPath, "templates"), formFile, "form.html.tmpl", modelData); err != nil {
			return fmt.Errorf("form template %s: %w", model.Name, err)
		}
	}
	return nil
}

// RenderMiddleware generates RBAC middleware and auth files
func RenderMiddleware(config ProjectConfig) error {
	g := &GormCodeGenerator{}
	data := buildTemplateData(config)
	if config.Port > 0 {
		data.Port = fmt.Sprintf("%d", config.Port)
	} else {
		data.Port = "8080"
	}

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
	if err := g.renderHTMLFile(filepath.Join(config.TargetPath, "templates"), "register.html", "register.html.tmpl", data); err != nil {
		return fmt.Errorf("register template: %w", err)
	}
	if err := g.renderHTMLFile(filepath.Join(config.TargetPath, "templates"), "forgot_password.html", "forgot_password.html.tmpl", data); err != nil {
		return fmt.Errorf("forgot_password template: %w", err)
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

// renderBasePages generates all dashboard base page templates and the base handler
func (g *GormCodeGenerator) renderBasePages(targetPath string, data TemplateData) error {
	// Base handler (Go file)
	if err := g.renderGoFile(filepath.Join(targetPath, "handlers"), "base.go", "base_handler.go.tmpl", data); err != nil {
		return fmt.Errorf("base handler: %w", err)
	}

	// Base HTML pages (use << >> delimiters)
	basePages := []struct {
		output string
		tmpl   string
	}{
		{"dashboard.html", "dashboard.html.tmpl"},
		{"leads.html", "leads.html.tmpl"},
		{"transactions.html", "transactions.html.tmpl"},
		{"charts.html", "charts.html.tmpl"},
		{"integration.html", "integration.html.tmpl"},
		{"calendar.html", "calendar.html.tmpl"},
		{"profile_settings.html", "profile_settings.html.tmpl"},
		{"team.html", "team.html.tmpl"},
		{"billing.html", "billing.html.tmpl"},
		{"welcome.html", "welcome.html.tmpl"},
		{"blank.html", "blank.html.tmpl"},
		{"404.html", "404.html.tmpl"},
	}

	tmplDir := filepath.Join(targetPath, "templates")
	for _, p := range basePages {
		if err := g.renderHTMLFile(tmplDir, p.output, p.tmpl, data); err != nil {
			return fmt.Errorf("%s: %w", p.output, err)
		}
	}

	return nil
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
