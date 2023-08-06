package database

import (
	"fmt"
	"strings"
)

type Query struct {
	Select     string
	From       string
	Join       string
	Where      string
	OrderBy    string
	GroupBy    string
	Having     string
	Limit      string
	PrimaryKey string
}

func (q Query) String() string {
	where := ""
	if q.Where != "" {
		where = fmt.Sprintf("WHERE %s", q.Where)
	}

	orderBy := ""
	if q.OrderBy != "" {
		orderBy = fmt.Sprintf("ORDER BY %s", q.OrderBy)
	}

	groupBy := ""
	if q.GroupBy != "" {
		groupBy = fmt.Sprintf("GROUP BY %s", q.GroupBy)
	}

	having := ""
	if q.Having != "" {
		having = fmt.Sprintf("HAVING %s", q.Having)
	}

	limit := ""
	if q.Limit != "" {
		limit = fmt.Sprintf("LIMIT %s", q.Limit)
	}

	return sanitizeStatement(
		strings.Join([]string{
			fmt.Sprintf("SELECT %s", q.Select),
			fmt.Sprintf("FROM `%s`", q.From),
			q.Join,
			where,
			groupBy,
			having,
			orderBy,
			limit,
		}, " "),
	)
}
