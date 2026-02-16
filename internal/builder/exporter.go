package builder

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Exporter handles exporting builder projects to static HTML
type Exporter struct {
	project *BuilderProject
}

// Export writes static HTML files to the output directory
func (e *Exporter) Export(outputPath string) error {
	if err := os.MkdirAll(outputPath, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	for _, page := range e.project.Pages {
		html := e.renderPage(page)
		filename := page.ID + ".html"
		if page.ID == "page-1" {
			filename = "index.html"
		}
		if err := os.WriteFile(filepath.Join(outputPath, filename), []byte(html), 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", filename, err)
		}
	}
	return nil
}

func (e *Exporter) renderPage(page Page) string {
	var body strings.Builder
	for _, comp := range page.Components {
		body.WriteString("    ")
		body.WriteString(comp.ToHTML())
		body.WriteString("\n")
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="ko">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>%s - %s</title>
    <script src="https://unpkg.com/htmx.org@1.9.12"></script>
    <script src="https://cdn.tailwindcss.com"></script>
</head>
<body class="bg-gray-100">
%s
</body>
</html>`, e.project.Name, page.Name, body.String())
}
