package application

import (
	"fmt"

	"ggami-go/internal/domain"
)

// ResolveDependencies resolves module dependencies using topological sort (Kahn's algorithm).
// Returns modules in dependency order. Returns error on missing dependencies or cycles.
func ResolveDependencies(selectedIDs []string, registry []domain.ModuleDef) ([]domain.ModuleDef, error) {
	// Build lookup map
	byID := make(map[string]domain.ModuleDef)
	for _, mod := range registry {
		byID[mod.ID] = mod
	}

	// Validate selected IDs exist
	selected := make(map[string]bool)
	for _, id := range selectedIDs {
		if _, ok := byID[id]; !ok {
			return nil, fmt.Errorf("unknown module %q", id)
		}
		selected[id] = true
	}

	// Check all dependencies are satisfied
	for _, id := range selectedIDs {
		mod := byID[id]
		for _, dep := range mod.Dependencies {
			if !selected[dep] {
				return nil, fmt.Errorf("module %q requires module %q which is not selected", id, dep)
			}
		}
	}

	// Kahn's algorithm for topological sort
	inDegree := make(map[string]int)
	for _, id := range selectedIDs {
		if _, ok := inDegree[id]; !ok {
			inDegree[id] = 0
		}
		mod := byID[id]
		for _, dep := range mod.Dependencies {
			inDegree[id]++
			_ = dep // dep counted as provider
		}
	}

	// Find nodes with no incoming edges
	var queue []string
	for _, id := range selectedIDs {
		if inDegree[id] == 0 {
			queue = append(queue, id)
		}
	}

	var sorted []domain.ModuleDef
	visited := make(map[string]bool)

	for len(queue) > 0 {
		id := queue[0]
		queue = queue[1:]

		if visited[id] {
			continue
		}
		visited[id] = true
		sorted = append(sorted, byID[id])

		// For each module that depends on this one, decrement in-degree
		for _, othID := range selectedIDs {
			if visited[othID] {
				continue
			}
			oth := byID[othID]
			for _, dep := range oth.Dependencies {
				if dep == id {
					inDegree[othID]--
					if inDegree[othID] == 0 {
						queue = append(queue, othID)
					}
				}
			}
		}
	}

	if len(sorted) != len(selectedIDs) {
		return nil, fmt.Errorf("circular dependency detected among modules")
	}

	return sorted, nil
}
