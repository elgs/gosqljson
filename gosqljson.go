package gosqljson

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
)

var QueryDbToArrayJson = func(db *sql.DB, sqlStatement string, sqlParams ...interface{}) string {
	data := QueryDbToArray(db, sqlStatement, sqlParams...)
	jsonString, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	}
	return string(jsonString)
}

var QueryDbToMapJson = func(db *sql.DB, sqlStatement string, sqlParams ...interface{}) string {
	data := QueryDbToMap(db, sqlStatement, sqlParams...)
	jsonString, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	}
	return string(jsonString)
}

var QueryDbToArray = func(db *sql.DB, sqlStatement string, sqlParams ...interface{}) (results [][]string) {
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
	} else {
		_, err := db.Exec(sqlStatement, sqlParams...)
		if err != nil {
			fmt.Println(err)
		}
	}
	return
}

var QueryDbToMap = func(db *sql.DB, sqlStatement string, sqlParams ...interface{}) (results []map[string]string) {
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
					result[cols[i]] = string(raw)
				}
			}
			results = append(results, result)
		}
	} else {
		_, err := db.Exec(sqlStatement, sqlParams...)
		if err != nil {
			fmt.Println(err)
		}
	}
	return
}
