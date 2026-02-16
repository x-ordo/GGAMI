package generator

import (
	"strings"
	"text/template"
	"unicode"
)

// TemplateData is the root data passed to all GORM templates
type TemplateData struct {
	ProjectName string
	DBType      DBType
	DBServer    string
	DBUser      string
	DBPw        string
	DBName      string
	Port        string
	Models      []ModelTmplData
	Driver      DBDriverInfo
	RBAC        *RBACConfig
	HasRBAC     bool
	RBACMatrix  string // Pre-built Go source for permission matrix
}

// ModelTmplData is per-model data for templates
type ModelTmplData struct {
	Name       string // PascalCase
	NameLower  string // camelCase first char lower
	NameSnake  string // snake_case
	NamePlural string // simple plural
	Fields     []FieldTmplData
}

// FieldTmplData is per-field data for templates
type FieldTmplData struct {
	Name       string
	Type       string // Go type
	GormTag    string // full GORM tag string
	JsonName   string
	InputType  string // HTML input type
	DefaultVal string
	IsID       bool
}

// GormFuncMap provides template helper functions
var GormFuncMap = template.FuncMap{
	"snake":   toSnakeCase,
	"lower":   strings.ToLower,
	"lower1":  lowerFirst,
	"plural":  simplePlural,
	"goType":  goType,
	"inputType": htmlInputType,
	"join":    strings.Join,
	"joinGormTags": joinGormTags,
}

func buildTemplateData(config ProjectConfig) TemplateData {
	driver := DBDriverMap[config.DBType]
	if config.DBType == "" {
		driver = DBDriverMap[DBTypeSQLite]
	}

	port := "8080"
	if config.Port > 0 {
		port = strings.TrimLeft(strings.Replace(string(rune(config.Port+'0')), "", "", -1), "")
	}

	var models []ModelTmplData
	for _, m := range config.Models {
		mtd := ModelTmplData{
			Name:       m.Name,
			NameLower:  lowerFirst(m.Name),
			NameSnake:  toSnakeCase(m.Name),
			NamePlural: simplePlural(m.Name),
		}
		for _, f := range m.Fields {
			jsonName := f.JsonName
			if jsonName == "" {
				jsonName = toSnakeCase(f.Name)
			}
			ftd := FieldTmplData{
				Name:       f.Name,
				Type:       goType(f.Type),
				GormTag:    joinGormTags(f.GormTags),
				JsonName:   jsonName,
				InputType:  htmlInputType(f.Type),
				DefaultVal: f.DefaultVal,
				IsID:       containsTag(f.GormTags, "primaryKey"),
			}
			mtd.Fields = append(mtd.Fields, ftd)
		}
		models = append(models, mtd)
	}

	hasRBAC := config.RBAC != nil && config.RBAC.Enabled

	rbacMatrix := ""
	if hasRBAC {
		rbacMatrix = buildRBACMatrixSource(config.RBAC)
	}

	return TemplateData{
		ProjectName: config.ProjectName,
		DBType:      config.DBType,
		DBServer:    config.DBServer,
		DBUser:      config.DBUser,
		DBPw:        config.DBPw,
		DBName:      config.DBName,
		Port:        port,
		Models:      models,
		Driver:      driver,
		RBAC:        config.RBAC,
		HasRBAC:     hasRBAC,
		RBACMatrix:  rbacMatrix,
	}
}

func buildRBACMatrixSource(rbac *RBACConfig) string {
	var b strings.Builder
	b.WriteString("map[string]map[string]Permission{\n")
	for _, role := range rbac.Roles {
		b.WriteString("\t\"" + role + "\": {\n")
		for _, mp := range rbac.ModelPerms {
			for _, perm := range mp.Permissions {
				if perm.Role == role {
					b.WriteString("\t\t\"" + mp.ModelName + "\": {")
					b.WriteString("Create: " + boolStr(perm.Create) + ", ")
					b.WriteString("Read: " + boolStr(perm.Read) + ", ")
					b.WriteString("Update: " + boolStr(perm.Update) + ", ")
					b.WriteString("Delete: " + boolStr(perm.Delete))
					b.WriteString("},\n")
				}
			}
		}
		b.WriteString("\t},\n")
	}
	b.WriteString("}")
	return b.String()
}

func boolStr(v bool) string {
	if v {
		return "true"
	}
	return "false"
}

func toSnakeCase(s string) string {
	var result []rune
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				result = append(result, '_')
			}
			result = append(result, unicode.ToLower(r))
		} else {
			result = append(result, r)
		}
	}
	return string(result)
}

func lowerFirst(s string) string {
	if s == "" {
		return s
	}
	return strings.ToLower(s[:1]) + s[1:]
}

func simplePlural(s string) string {
	if s == "" {
		return s
	}
	lower := strings.ToLower(s)
	if strings.HasSuffix(lower, "s") || strings.HasSuffix(lower, "x") || strings.HasSuffix(lower, "ch") || strings.HasSuffix(lower, "sh") {
		return s + "es"
	}
	if strings.HasSuffix(lower, "y") && len(s) > 1 {
		prev := s[len(s)-2]
		if prev != 'a' && prev != 'e' && prev != 'i' && prev != 'o' && prev != 'u' {
			return s[:len(s)-1] + "ies"
		}
	}
	return s + "s"
}

func goType(t string) string {
	switch t {
	case "string":
		return "string"
	case "int":
		return "int"
	case "uint":
		return "uint"
	case "float64":
		return "float64"
	case "bool":
		return "bool"
	case "time.Time":
		return "time.Time"
	default:
		return "string"
	}
}

func htmlInputType(t string) string {
	switch t {
	case "string":
		return "text"
	case "int", "uint":
		return "number"
	case "float64":
		return "number"
	case "bool":
		return "checkbox"
	case "time.Time":
		return "datetime-local"
	default:
		return "text"
	}
}

func joinGormTags(tags []string) string {
	if len(tags) == 0 {
		return ""
	}
	return strings.Join(tags, ";")
}

func containsTag(tags []string, tag string) bool {
	for _, t := range tags {
		if strings.EqualFold(t, tag) {
			return true
		}
	}
	return false
}
