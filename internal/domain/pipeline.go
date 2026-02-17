package domain

// PipelineContext carries state through the generation pipeline
type PipelineContext struct {
	Config   ProjectConfig
	Language string
	TempDir  string      // temporary directory during generation
	FinalDir string      // final output path
	Modules  []ModuleDef // dependency-sorted active modules
}
