package application

import (
	"fmt"

	"ggami/internal/domain"
)

// PipelineStep defines a single step in the generation pipeline
type PipelineStep interface {
	Name() string
	Execute(ctx *domain.PipelineContext) error
	Rollback(ctx *domain.PipelineContext) error
}

// Pipeline runs steps sequentially with rollback on failure
type Pipeline struct {
	steps     []PipelineStep
	completed []PipelineStep
}

// NewPipeline creates a pipeline with the given steps
func NewPipeline(steps ...PipelineStep) *Pipeline {
	return &Pipeline{steps: steps}
}

// Run executes all steps. On failure, rolls back completed steps in reverse order.
func (p *Pipeline) Run(ctx *domain.PipelineContext) error {
	for _, step := range p.steps {
		if err := step.Execute(ctx); err != nil {
			// Rollback completed steps in reverse
			for i := len(p.completed) - 1; i >= 0; i-- {
				if rbErr := p.completed[i].Rollback(ctx); rbErr != nil {
					return fmt.Errorf("step %q failed: %w (rollback of %q also failed: %v)",
						step.Name(), err, p.completed[i].Name(), rbErr)
				}
			}
			return fmt.Errorf("step %q failed: %w", step.Name(), err)
		}
		p.completed = append(p.completed, step)
	}
	return nil
}
