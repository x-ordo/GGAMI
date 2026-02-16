package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"ggami/internal/builder"
)

// NewHandler creates a Chi router for HTMX partial HTML responses
func NewHandler(pm *builder.ProjectManager) http.Handler {
	r := chi.NewRouter()

	// GET /api/components - Component library cards
	r.Get("/api/components", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		components := []struct {
			Type string
			Name string
			Icon string
			Desc string
		}{
			{"heading", "Heading", "H", "Title text"},
			{"paragraph", "Paragraph", "P", "Body text"},
			{"button", "Button", "B", "Click action"},
			{"image", "Image", "I", "Image element"},
			{"container", "Container", "C", "Layout box"},
			{"hero", "Hero", "H", "Hero section"},
			{"navbar", "Navbar", "N", "Navigation bar"},
			{"card", "Card", "C", "Content card"},
			{"form", "Form", "F", "Input form"},
			{"footer", "Footer", "Ft", "Page footer"},
		}

		var html strings.Builder
		for _, c := range components {
			html.WriteString(fmt.Sprintf(`
<div draggable="true" ondragstart="onDragStart(event, '%s')"
     class="p-3 bg-gray-700 rounded border border-gray-600 cursor-move hover:border-yellow-400 transition mb-2">
    <div class="flex items-center gap-2">
        <span class="w-8 h-8 bg-gray-600 rounded flex items-center justify-center text-xs font-bold text-yellow-400">%s</span>
        <div>
            <div class="font-bold text-sm">%s</div>
            <div class="text-xs text-gray-400">%s</div>
        </div>
    </div>
</div>`, c.Type, c.Icon, c.Name, c.Desc))
		}
		fmt.Fprint(w, html.String())
	})

	// POST /api/canvas/add - Render a dropped component as HTML
	r.Post("/api/canvas/add", func(w http.ResponseWriter, r *http.Request) {
		compType := r.FormValue("type")
		compID := r.FormValue("id")

		comp := builder.Component{
			ID:   compID,
			Type: compType,
			Styles: map[string]string{
				"class": getDefaultClass(compType),
			},
		}

		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `<div id="comp-%s" class="relative group border-2 border-transparent hover:border-blue-400 rounded p-1 mb-2"
     onclick="selectComponent('%s', '%s')">
    <div class="absolute top-1 right-1 hidden group-hover:flex gap-1">
        <button onclick="event.stopPropagation(); deleteComponent('%s')"
                class="bg-red-500 text-white text-xs px-2 py-0.5 rounded">X</button>
    </div>
    %s
</div>`, compID, compID, compType, compID, comp.ToHTML())
	})

	// GET /api/properties/{id} - Property edit form
	r.Get("/api/properties/{id}", func(w http.ResponseWriter, r *http.Request) {
		compID := chi.URLParam(r, "id")
		compType := r.URL.Query().Get("type")

		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `
<form hx-post="/api/properties/%s" hx-target="#comp-%s" hx-swap="outerHTML" class="space-y-3">
    <input type="hidden" name="id" value="%s" />
    <input type="hidden" name="type" value="%s" />
    <div>
        <label class="block text-sm font-medium text-gray-300 mb-1">Content</label>
        <input type="text" name="content" placeholder="Enter content..."
               class="w-full bg-gray-700 border border-gray-600 rounded p-2 text-sm" />
    </div>
    <div>
        <label class="block text-sm font-medium text-gray-300 mb-1">CSS Classes</label>
        <input type="text" name="classes" placeholder="Tailwind classes..."
               class="w-full bg-gray-700 border border-gray-600 rounded p-2 text-sm" />
    </div>
    <button type="submit" class="w-full bg-blue-600 text-white py-2 rounded text-sm hover:bg-blue-700">
        Update
    </button>
</form>`, compID, compID, compID, compType)
	})

	// POST /api/properties/{id} - Update component and re-render
	r.Post("/api/properties/{id}", func(w http.ResponseWriter, r *http.Request) {
		compID := r.FormValue("id")
		compType := r.FormValue("type")
		content := r.FormValue("content")
		classes := r.FormValue("classes")

		comp := builder.Component{
			ID:      compID,
			Type:    compType,
			Content: content,
			Styles: map[string]string{
				"class": classes,
			},
		}

		// Also update in project manager if project exists
		pm := pm
		if proj := pm.GetCurrentProject(); proj != nil {
			pm.UpdateComponent("page-1", compID, comp)
		}

		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `<div id="comp-%s" class="relative group border-2 border-transparent hover:border-blue-400 rounded p-1 mb-2"
     onclick="selectComponent('%s', '%s')">
    <div class="absolute top-1 right-1 hidden group-hover:flex gap-1">
        <button onclick="event.stopPropagation(); deleteComponent('%s')"
                class="bg-red-500 text-white text-xs px-2 py-0.5 rounded">X</button>
    </div>
    %s
</div>`, compID, compID, compType, compID, comp.ToHTML())
	})

	// DELETE /api/canvas/{id} - Remove component
	r.Delete("/api/canvas/{id}", func(w http.ResponseWriter, r *http.Request) {
		compID := chi.URLParam(r, "id")

		if proj := pm.GetCurrentProject(); proj != nil {
			pm.DeleteComponent("page-1", compID)
		}

		w.Header().Set("Content-Type", "text/html")
		// Return empty to remove from DOM
		w.WriteHeader(http.StatusOK)
	})

	return r
}

func getDefaultClass(compType string) string {
	switch compType {
	case "heading":
		return "text-3xl font-bold"
	case "paragraph":
		return "text-gray-600"
	case "button":
		return "bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700"
	case "hero":
		return ""
	case "navbar":
		return ""
	case "card":
		return ""
	case "form":
		return ""
	case "footer":
		return ""
	default:
		return ""
	}
}
