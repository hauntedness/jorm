package jormgen

import "strings"

func AddArray[T any](column string, methodParams []T, sqlParams []any) string {
	var query []string = make([]string, 0, len(methodParams))
	for _, v := range methodParams {
		sqlParams = append(sqlParams, v)
		query = append(query, "?")
	}
	return column + " in ( " + strings.Join(query, ",") + " )"
}
