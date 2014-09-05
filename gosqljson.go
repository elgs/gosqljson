package gosqljson

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

func QueryDbToArrayJson(db *sql.DB, toLower bool, sqlStatement string, sqlParams ...interface{}) (string, error) {
	data, err := QueryDbToArray(db, toLower, sqlStatement, sqlParams...)
	jsonString, err := json.Marshal(data)
	return string(jsonString), err
}

func QueryDbToMapJson(db *sql.DB, toLower bool, sqlStatement string, sqlParams ...interface{}) (string, error) {
	data, err := QueryDbToMap(db, toLower, sqlStatement, sqlParams...)
	jsonString, err := json.Marshal(data)
	return string(jsonString), err
}

func QueryDbToArray(db *sql.DB, toLower bool, sqlStatement string, sqlParams ...interface{}) ([][]string, error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	SqlSafe(&sqlStatement)

	var results [][]string
	if strings.HasPrefix(strings.ToUpper(sqlStatement), "SELECT") {
		rows, err := db.Query(sqlStatement, sqlParams...)
		if err != nil {
			return results, err
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
					result[i] = ""
				} else {
					result[i] = string(raw)
				}
			}
			results = append(results, result)
		}
	}
	return results, nil
}

func QueryDbToMap(db *sql.DB, toLower bool, sqlStatement string, sqlParams ...interface{}) ([]map[string]string, error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	SqlSafe(&sqlStatement)

	var results []map[string]string
	if strings.HasPrefix(strings.ToUpper(sqlStatement), "SELECT ") {
		rows, err := db.Query(sqlStatement, sqlParams...)
		if err != nil {
			return results, err
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
					if toLower {
						result[colsLower[i]] = ""
					} else {
						result[cols[i]] = ""
					}
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
	return results, nil
}

func ExecDb(db *sql.DB, sqlStatement string, sqlParams ...interface{}) (int64, error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	SqlSafe(&sqlStatement)

	sqlUpper := strings.ToUpper(sqlStatement)
	if strings.HasPrefix(sqlUpper, "UPDATE ") ||
		strings.HasPrefix(sqlUpper, "INSERT ") ||
		strings.HasPrefix(sqlUpper, "DELETE FROM ") {
		result, err := db.Exec(sqlStatement, sqlParams...)
		if err != nil {
			fmt.Println(err)
			return 0, err
		}
		return result.RowsAffected()
	}
	return 0, errors.New(fmt.Sprint("Invalid SQL:", sqlStatement))
}

func SqlSafe(s *string) {
	*s = strings.Replace(*s, "'", "''", -1)
	*s = strings.Replace(*s, "--", "", -1)
}
