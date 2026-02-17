package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"ggami-go/internal/builder"
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
			{"heading", "제목", "H", "제목 텍스트"},
			{"paragraph", "문단", "P", "본문 텍스트"},
			{"button", "버튼", "B", "클릭 액션"},
			{"image", "이미지", "I", "이미지 요소"},
			{"container", "컨테이너", "C", "레이아웃 박스"},
			{"hero", "히어로", "H", "히어로 섹션"},
			{"navbar", "내비게이션", "N", "상단 메뉴바"},
			{"card", "카드", "C", "콘텐츠 카드"},
			{"form", "폼", "F", "입력 폼"},
			{"footer", "푸터", "Ft", "하단 영역"},
		}

		var html strings.Builder
		for _, c := range components {
			html.WriteString(fmt.Sprintf(`
<div draggable="true" ondragstart="onDragStart(event, '%s')"
     class="card card-compact bg-base-300 cursor-move hover:border-warning border border-base-content/10 transition mb-2">
    <div class="card-body p-3 flex-row items-center gap-2">
        <span class="badge badge-warning font-bold">%s</span>
        <div>
            <div class="font-bold text-sm">%s</div>
            <div class="text-xs text-base-content/50">%s</div>
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
		fmt.Fprintf(w, `<div id="comp-%s" class="relative group border-2 border-transparent hover:border-primary rounded p-1 mb-2"
     onclick="selectComponent('%s', '%s')">
    <div class="absolute top-1 right-1 hidden group-hover:flex gap-1">
        <button onclick="event.stopPropagation(); deleteComponent('%s')"
                class="btn btn-xs btn-error">X</button>
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
    <div class="form-control">
        <label class="label"><span class="label-text">콘텐츠</span></label>
        <input type="text" name="content" placeholder="내용을 입력하세요..."
               class="input input-bordered input-sm w-full" />
    </div>
    <div class="form-control">
        <label class="label"><span class="label-text">CSS 클래스</span></label>
        <input type="text" name="classes" placeholder="Tailwind 클래스..."
               class="input input-bordered input-sm w-full" />
    </div>
    <button type="submit" class="btn btn-primary btn-sm w-full">
        적용
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
		fmt.Fprintf(w, `<div id="comp-%s" class="relative group border-2 border-transparent hover:border-primary rounded p-1 mb-2"
     onclick="selectComponent('%s', '%s')">
    <div class="absolute top-1 right-1 hidden group-hover:flex gap-1">
        <button onclick="event.stopPropagation(); deleteComponent('%s')"
                class="btn btn-xs btn-error">X</button>
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
		return "text-base-content"
	case "button":
		return "btn btn-primary"
	default:
		return ""
	}
}
