package gosqljson

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

var QueryDbToArrayJson = func(db *sql.DB, toLower bool, sqlStatement string, sqlParams ...interface{}) string {
	data := QueryDbToArray(db, toLower, sqlStatement, sqlParams...)
	jsonString, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	}
	return string(jsonString)
}

var QueryDbToMapJson = func(db *sql.DB, toLower bool, sqlStatement string, sqlParams ...interface{}) string {
	data := QueryDbToMap(db, toLower, sqlStatement, sqlParams...)
	jsonString, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	}
	return string(jsonString)
}

var QueryDbToArray = func(db *sql.DB, toLower bool, sqlStatement string, sqlParams ...interface{}) (results [][]string) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	if strings.HasPrefix(strings.ToUpper(sqlStatement), "SELECT") {
		rows, err := db.Query(sqlStatement, sqlParams...)
		if err != nil {
			fmt.Println("db.Query:", err)
		}
		cols, _ := rows.Columns()
		if toLower {
			colsLower := make([]string, len(cols))
			for i, v := range cols {
				colsLower[i] = strings.ToLower(v)
			}
			results = append(results, colsLower)
		} else {
			results = append(results, cols)
		}

		rawResult := make([][]byte, len(cols))

		dest := make([]interface{}, len(cols)) // A temporary interface{} slice
		for i, _ := range rawResult {
			dest[i] = &rawResult[i] // Put pointers to each string in the interface slice
		}

		for rows.Next() {
			result := make([]string, len(cols))
			rows.Scan(dest...)
			for i, raw := range rawResult {
				if raw == nil {
					result[i] = "\\N"
				} else {
					result[i] = string(raw)
				}
			}
			results = append(results, result)
		}
	}
	return
}

var QueryDbToMap = func(db *sql.DB, toLower bool, sqlStatement string, sqlParams ...interface{}) (results []map[string]string) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	if strings.HasPrefix(strings.ToUpper(sqlStatement), "SELECT") {
		rows, err := db.Query(sqlStatement, sqlParams...)
		if err != nil {
			fmt.Println("db.Query:", err)
		}
		cols, _ := rows.Columns()
		colsLower := make([]string, len(cols))
		if toLower {
			for i, v := range cols {
				colsLower[i] = strings.ToLower(v)
			}
		}
		rawResult := make([][]byte, len(cols))

		dest := make([]interface{}, len(cols)) // A temporary interface{} slice
		for i, _ := range rawResult {
			dest[i] = &rawResult[i] // Put pointers to each string in the interface slice
		}

		for rows.Next() {
			result := make(map[string]string, len(cols))
			rows.Scan(dest...)
			for i, raw := range rawResult {
				if raw == nil {
					result[cols[i]] = "\\N"
				} else {
					if toLower {
						result[colsLower[i]] = string(raw)
					} else {
						result[cols[i]] = string(raw)
					}
				}
			}
			results = append(results, result)
		}
	}
	return
}

var ExecDb = func(db *sql.DB, toLower bool, sqlStatement string, sqlParams ...interface{}) (int64, error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	sqlUpper := strings.ToUpper(sqlStatement)
	if strings.HasPrefix(sqlUpper, "UPDATE") ||
		strings.HasPrefix(sqlUpper, "INSERT INTO") ||
		strings.HasPrefix(sqlUpper, "DELETE FROM") {
		result, err := db.Exec(sqlStatement, sqlParams...)
		if err != nil {
			fmt.Println(err)
			return 0, err
		}
		return result.RowsAffected()
	}
	return 0, errors.New(fmt.Sprint("Invalid SQL:", sqlStatement))
}
