package database

import (
	"fmt"
	"reflect"
	"strconv"
)

type SQLBuild struct {
}

func (query SQLBuild) Insert(table string, v any) string {
	columns, valuesCount := query.columnsAndValuesCountToInsert(v)
	return fmt.Sprintf(`INSERT INTO %s(%s) VALUES(%s)`, table, columns, valuesCount)
}

func (query SQLBuild) Select(table string, v any) string {
	columns, _ := query.columnsAndValuesCount(v)
	return fmt.Sprintf(`SELECT %s FROM %s`, columns, table)
}

func (query SQLBuild) Pagination(sql, pageIndex, pageSize string) string {
	if pageIndex != "" {
		sql = fmt.Sprintf(`%s limit %s`, sql, pageSize)
	}

	if pageSize != "" {
		sql = fmt.Sprintf(`%s offset %s`, sql, pageIndex)
	}

	return sql
}

func (query SQLBuild) columnsAndValuesCount(v any) (string, string) {
	fields := reflect.ValueOf(v)
	selecto := ""
	count := ""
	declaredFieldsSize := fields.NumField()

	for i := 0; i < declaredFieldsSize; i++ {
		fieldValue := fields.Type().Field(i).Tag
		valor := fieldValue.Get("column")
		if valor != "" {
			selecto = selecto + valor
			count = count + "$" + strconv.Itoa(i+1)

			if i != declaredFieldsSize-1 {
				selecto = selecto + ", "
				count = count + ", "
			}

		}

	}
	return selecto, count
}

func (query SQLBuild) columnsAndValuesCountToInsert(v any) (string, string) {
	fields := reflect.ValueOf(v)
	selecto := ""
	count := ""
	declaredFieldsSize := fields.NumField()

	for i := 0; i < declaredFieldsSize; i++ {
		fieldValue := fields.Type().Field(i).Tag
		valor := fieldValue.Get("column")

		if valor != "" && valor != "id" {
			selecto = selecto + valor
			count = count + "$" + strconv.Itoa(i)

			if i != declaredFieldsSize-1 {
				selecto = selecto + ", "
				count = count + ", "
			}

		}

	}
	return selecto, count
}
