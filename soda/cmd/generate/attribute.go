package generate

import (
	"fmt"
	"strings"

	"github.com/markbates/inflect"
)

type attribute struct {
	Model        *model
	Name         inflect.Name
	OriginalType string
	GoType       string
	Nullable     bool
}

func (a attribute) String() string {
	var s string
	if a.Model != nil {
		if a.Model.MarshalType == "jsonapi" {
			if a.Name == "id" {
				s = fmt.Sprintf("\t%s %s `%s:\"primary,%s\" db:\"%s\"`", a.Name.Camel(), a.GoType, a.Model.MarshalType, a.Model.Name.PluralUnder(), a.Name.Underscore())
			} else {
				s = fmt.Sprintf("\t%s %s `%s:\"attr,%s\" db:\"%s\"`", a.Name.Camel(), a.GoType, a.Model.MarshalType, a.Name.Underscore(), a.Name.Underscore())
			}
		} else {
			s = fmt.Sprintf("\t%s %s `%s:\"%s\" db:\"%s\"`", a.Name.Camel(), a.GoType, a.Model.MarshalType, a.Name.Underscore(), a.Name.Underscore())
		}
	}
	return s
}

func (a attribute) IsValidable() bool {
	return a.GoType == "string" || a.GoType == "time.Time" || a.GoType == "int"
}

func newAttribute(base string, model *model) attribute {
	col := strings.Split(base, ":")
	if len(col) == 1 {
		col = append(col, "string")
	}

	nullable := strings.HasPrefix(col[1], "nulls.")
	if !model.HasNulls && nullable {
		model.HasNulls = true
		model.Imports = append(model.Imports, "github.com/gobuffalo/pop/nulls")
	} else if !model.HasSlices && strings.HasPrefix(col[1], "slices.") {
		model.HasSlices = true
		model.Imports = append(model.Imports, "github.com/gobuffalo/pop/slices")
	} else if !model.HasUUID && col[1] == "uuid" {
		model.HasUUID = true
		model.Imports = append(model.Imports, "github.com/gobuffalo/uuid")
	}

	got := colType(col[1])
	if len(col) > 2 {
		got = col[2]
	}
	a := attribute{
		Model:        model,
		Name:         inflect.Name(col[0]),
		OriginalType: col[1],
		GoType:       got,
		Nullable:     nullable,
	}

	return a
}

func colType(s string) string {
	switch strings.ToLower(s) {
	case "text":
		return "string"
	case "time", "timestamp", "datetime":
		return "time.Time"
	case "nulls.text":
		return "nulls.String"
	case "uuid":
		return "uuid.UUID"
	case "json", "jsonb":
		return "slices.Map"
	case "[]string":
		return "slices.String"
	case "[]int":
		return "slices.Int"
	case "slices.float", "[]float", "[]float32", "[]float64":
		return "slices.Float"
	case "decimal", "float":
		return "float64"
	case "[]byte", "blob":
		return "[]byte"
	default:
		return s
	}
}
