package generator

// FieldDef defines a single field in a GORM model
type FieldDef struct {
	Name       string   `json:"name"`       // PascalCase: "Title"
	Type       string   `json:"type"`       // "string","int","uint","float64","bool","time.Time"
	GormTags   []string `json:"gormTags"`   // ["primaryKey","unique","not null","index"]
	DefaultVal string   `json:"defaultVal"` // default value
	JsonName   string   `json:"jsonName"`   // auto snake_case from frontend
}

// ModelDef defines a GORM model with its fields
type ModelDef struct {
	Name   string     `json:"name"`   // PascalCase: "Product"
	Fields []FieldDef `json:"fields"`
}

// DBType represents a supported database type
type DBType string

const (
	DBTypeMSSQL    DBType = "mssql"
	DBTypePostgres DBType = "postgres"
	DBTypeMySQL    DBType = "mysql"
	DBTypeSQLite   DBType = "sqlite"
)

// RolePermission defines CRUD permissions for a single role
type RolePermission struct {
	Role   string `json:"role"`   // "admin","editor","viewer"
	Create bool   `json:"create"`
	Read   bool   `json:"read"`
	Update bool   `json:"update"`
	Delete bool   `json:"delete"`
}

// ModelRBAC defines per-model RBAC permissions
type ModelRBAC struct {
	ModelName   string           `json:"modelName"`
	Permissions []RolePermission `json:"permissions"`
}

// RBACConfig holds all RBAC/JWT configuration
type RBACConfig struct {
	Enabled    bool        `json:"enabled"`
	Roles      []string    `json:"roles"`
	JWTSecret  string      `json:"jwtSecret"`
	ModelPerms []ModelRBAC `json:"modelPerms"`
}
