package store

import (
	"strings"

	"github.com/jinzhu/gorm"
)

type (
	bulkRow     []interface{}
	bulkRowFunc func(int) bulkRow
)

func bulkPlaceholders(cols int, rows int) string {
	lines := make([]string, rows)

	for i := 0; i < rows; i++ {
		l := make([]string, cols)
		for j := 0; j < cols; j++ {
			l[j] = "?"
		}
		lines[i] = "(" + strings.Join(l, ",") + ")"
	}

	return strings.Join(lines, ",")
}

func bulkInsert(db *gorm.DB, query string, n int, rowfunc bulkRowFunc) error {
	var placeholders string
	var vals []interface{}

	for i := 0; i < n; i++ {
		row := rowfunc(i)
		if placeholders == "" {
			placeholders = bulkPlaceholders(len(row), n)
		}
		vals = append(vals, row...)
	}

	sql := strings.Replace(query, "@values", placeholders, 1)

	return db.Exec(sql, vals...).Error
}
