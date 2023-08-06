package database

import (
	"fmt"
	"reflect"
	"strings"
)

type dbalTag struct {
	Table      string
	Column     string
	PrimaryKey bool
	ReadOnly   bool
}

func GenerateQuery(e Entity) Query {
	q := Query{
		From: e.TableName(),
		Join: e.Joins(),
	}

	element := reflect.ValueOf(e)

	selects := []string{}

	loopOverStructFields(element, func(fieldDefinition reflect.StructField, fieldValue reflect.Value) {
		tag := parseTag(fieldDefinition.Tag)
		if tag.Column == "" {
			return
		}

		if q.PrimaryKey == "" && tag.PrimaryKey {
			q.PrimaryKey = tag.Column
		}

		if tag.Table != "" {
			selects = append(selects, fmt.Sprintf(
				"`%s`.`%s` as `%s.%s`",
				tag.Table,
				tag.Column,
				tag.Table,
				tag.Column,
			))

			return
		}

		selects = append(selects, fmt.Sprintf("`%s`", tag.Column))
	})

	q.Select = strings.Join(selects, ", ")

	return q
}

func generateInsert(entity Entity) string {
	return generateInsertOrSave(entity, false)
}

func generateInsertOrSave(entity Entity, updateOnDuplicate bool) string {
	columns := []string{}
	valuePlaceholders := []string{}
	updates := []string{}

	element := reflect.ValueOf(entity)
	loopOverStructFields(element, func(fieldDefinition reflect.StructField, fieldValue reflect.Value) {
		tag := parseTag(fieldDefinition.Tag)

		if tag.ReadOnly {
			return
		}

		if tag.Column == "" {
			return
		}

		// Skip over columns that we have joined on other tables to get
		if tag.Table != "" && tag.Table != entity.TableName() {
			return
		}

		columns = append(columns, fmt.Sprintf("`%s`", tag.Column))
		valuePlaceholders = append(valuePlaceholders, fmt.Sprintf(":%s", tag.Column))
		if !tag.PrimaryKey {
			updates = append(updates, fmt.Sprintf("`%s` = :%s", tag.Column, tag.Column))
		}
	})

	queryText := fmt.Sprintf(
		"INSERT INTO `%s` (%s) VALUES (%s)",
		entity.TableName(),
		strings.Join(columns, ", "),
		strings.Join(valuePlaceholders, ", "),
	)

	if updateOnDuplicate {
		queryText += fmt.Sprintf(
			" ON DUPLICATE KEY UPDATE %s",
			strings.Join(updates, ", "),
		)
	}

	return queryText
}

func generateUpdate(entity Entity) string {
	primaryKeys := []string{}
	updates := []string{}

	element := reflect.ValueOf(entity)
	loopOverStructFields(element, func(fieldDefinition reflect.StructField, fieldValue reflect.Value) {
		tag := parseTag(fieldDefinition.Tag)

		if tag.ReadOnly {
			return
		}

		if tag.Column == "" {
			return
		}

		// Skip over columns that we have joined on other tables to get
		if tag.Table != "" && tag.Table != entity.TableName() {
			return
		}

		if tag.PrimaryKey {
			primaryKeys = append(primaryKeys, fmt.Sprintf("`%s` = :%s", tag.Column, tag.Column))

			return
		}

		updates = append(updates, fmt.Sprintf("`%s` = :%s", tag.Column, tag.Column))
	})

	return fmt.Sprintf(
		"UPDATE `%s` SET %s WHERE %s",
		entity.TableName(),
		strings.Join(updates, ", "),
		strings.Join(primaryKeys, " AND "),
	)
}

func generateSave(entity Entity) string {
	return generateInsertOrSave(entity, true)
}

func generateDelete(entity Entity) string {
	deletes := []string{}

	element := reflect.ValueOf(entity)
	loopOverStructFields(element, func(fieldDefinition reflect.StructField, fieldValue reflect.Value) {
		tag := parseTag(fieldDefinition.Tag)
		if !tag.PrimaryKey || tag.Column == "" {
			return
		}

		deletes = append(deletes, fmt.Sprintf("`%s` = :%s", tag.Column, tag.Column))
	})

	return fmt.Sprintf(
		"DELETE FROM `%s` WHERE %s",
		entity.TableName(),
		strings.Join(deletes, " AND "),
	)
}
