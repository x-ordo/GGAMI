package application

import "ggami-go/internal/domain"

// GenerateProject orchestrates the full project generation using the pipeline
func GenerateProject(config domain.ProjectConfig, language string) (string, error) {
	ctx := &domain.PipelineContext{
		Config:   config,
		Language: language,
		FinalDir: config.TargetPath,
	}

	steps := buildSteps(config, language)
	if err := NewPipeline(steps...).Run(ctx); err != nil {
		return "", err
	}

	return "Generation complete: " + config.TargetPath, nil
}

func buildSteps(config domain.ProjectConfig, language string) []PipelineStep {
	var steps []PipelineStep

	// Common: validate, resolve modules, create temp dir
	steps = append(steps,
		&ValidateConfigStep{},
		&ResolveModulesStep{},
		&CreateTempDirStep{},
	)

	isGorm := config.GormMode && len(config.Models) > 0 && language == "go"

	if isGorm {
		// GORM mode: scaffold gorm dirs, then generate per-concern
		steps = append(steps,
			&ScaffoldGormStep{},
			&GenerateCoreStep{},
			&GenerateModelsStep{},
			&GenerateHandlersStep{},
			&GenerateTemplatesStep{},
		)
		if config.RBAC != nil && config.RBAC.Enabled {
			steps = append(steps, &GenerateMiddlewareStep{})
		}
	} else {
		// Legacy mode: scaffold, manifest, code + module injection
		steps = append(steps,
			&ScaffoldStep{},
			&GenerateCoreStep{},
			&InjectModulesStep{},
		)
	}

	// Common: finalize
	steps = append(steps, &FinalizeStep{})
	return steps
}
