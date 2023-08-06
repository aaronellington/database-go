package database

import (
	"reflect"
	"regexp"
	"strings"
)

func sanitizeStatement(query string) string {
	commentRegex := regexp.MustCompile(`(?m)\s*--.*$`)
	query = string(commentRegex.ReplaceAll([]byte(query), []byte("")))

	spaceRegex := regexp.MustCompile(`\s+`)
	query = string(spaceRegex.ReplaceAll([]byte(query), []byte(" ")))

	query = strings.TrimSpace(query)

	return query
}

func convertStructToParams(v any) map[string]any {
	parameters := map[string]any{}

	loopOverStructFields(reflect.ValueOf(v), func(fieldDefinition reflect.StructField, fieldValue reflect.Value) {
		tag := parseTag(fieldDefinition.Tag)
		if tag.Column == "" {
			return
		}

		parameters[tag.Column] = fieldValue.Interface()
	})

	return parameters
}

func parseTag(tagString reflect.StructTag) dbalTag {
	parts := strings.Split(tagString.Get("db"), ",")

	tag := dbalTag{}

	for i, part := range parts {
		if i == 0 {
			subPart := strings.Split(part, ".")

			if len(subPart) < 2 {
				tag.Column = part

				continue
			}

			tag.Table = subPart[0]
			tag.Column = subPart[1]
		}

		if part == "readOnly" {
			tag.ReadOnly = true

			continue
		}

		if part == "primaryKey" {
			tag.PrimaryKey = true

			continue
		}
	}

	return tag
}

func loopOverStructFields(value reflect.Value, fieldHandler func(fieldDefinition reflect.StructField, fieldValue reflect.Value)) {
	if value.Kind() == reflect.Pointer {
		value = value.Elem()
	}

	for i := 0; i < value.NumField(); i++ {
		fieldValue := value.Field(i)
		fieldDefinition := value.Type().Field(i)

		if !fieldDefinition.IsExported() {
			continue
		}

		if fieldDefinition.Type.Kind() == reflect.Struct && fieldDefinition.Anonymous {
			loopOverStructFields(fieldValue, fieldHandler)

			continue
		}

		fieldHandler(fieldDefinition, fieldValue)
	}
}
