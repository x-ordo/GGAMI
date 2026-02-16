package generator

// ProjectConfig holds all configuration for project generation
type ProjectConfig struct {
	ProjectName string   `json:"projectName"`
	TargetPath  string   `json:"targetPath"`
	DBServer    string   `json:"dbServer"`
	DBUser      string   `json:"dbUser"`
	DBPw        string   `json:"dbPw"`
	DBName      string   `json:"dbName"`
	Port        int      `json:"port,omitempty"`
	Modules     []string `json:"modules"`

	// GORM full-stack generation fields
	GormMode bool       `json:"gormMode,omitempty"`
	Models   []ModelDef `json:"models,omitempty"`
	DBType   DBType     `json:"dbType,omitempty"`
	RBAC     *RBACConfig `json:"rbac,omitempty"`
}

// Generator interface defines the contract for code generators
type Generator interface {
	// Scaffold creates the basic folder structure
	Scaffold(path string) error
	// CreateManifest creates config files (go.mod or package.json)
	CreateManifest(config ProjectConfig) error
	// GenerateCode creates source code files
	GenerateCode(config ProjectConfig) error
}
