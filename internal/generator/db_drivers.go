package generator

// DBDriverInfo holds import and connection info for a database driver
type DBDriverInfo struct {
	GormDriver string // GORM driver import path
	DialFunc   string // gorm.Open dialector function
	DSNFormat  string // DSN format string with placeholders
	GoModDep   string // go.mod dependency line
}

// DBDriverMap maps DBType to driver configuration
var DBDriverMap = map[DBType]DBDriverInfo{
	DBTypeMSSQL: {
		GormDriver: "gorm.io/driver/sqlserver",
		DialFunc:   "sqlserver.Open",
		DSNFormat:  `"sqlserver://%s:%s@%s?database=%s"`,
		GoModDep:   "gorm.io/driver/sqlserver v1.5.4",
	},
	DBTypePostgres: {
		GormDriver: "gorm.io/driver/postgres",
		DialFunc:   "postgres.Open",
		DSNFormat:  `"host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable"`,
		GoModDep:   "gorm.io/driver/postgres v1.5.11",
	},
	DBTypeMySQL: {
		GormDriver: "gorm.io/driver/mysql",
		DialFunc:   "mysql.Open",
		DSNFormat:  `"%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local"`,
		GoModDep:   "gorm.io/driver/mysql v1.5.7",
	},
	DBTypeSQLite: {
		GormDriver: "gorm.io/driver/sqlite",
		DialFunc:   "sqlite.Open",
		DSNFormat:  `"%s.db"`,
		GoModDep:   "gorm.io/driver/sqlite v1.5.7",
	},
}
