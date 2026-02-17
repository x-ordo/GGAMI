package generator

import "ggami/internal/domain"

// Type aliases for backward compatibility
type FieldDef = domain.FieldDef
type ModelDef = domain.ModelDef
type DBType = domain.DBType
type RolePermission = domain.RolePermission
type ModelRBAC = domain.ModelRBAC
type RBACConfig = domain.RBACConfig

// Re-export constants
const (
	DBTypeMSSQL    = domain.DBTypeMSSQL
	DBTypePostgres = domain.DBTypePostgres
	DBTypeMySQL    = domain.DBTypeMySQL
	DBTypeSQLite   = domain.DBTypeSQLite
)
