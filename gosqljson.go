package gosqljson

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
)

const (
	AsIs = iota
	Lower
	Upper
	Camel
)

type DB interface {
	Query(query string, args ...any) (*sql.Rows, error)
	Exec(query string, args ...any) (sql.Result, error)
}

// QueryToArrayJSON - run the sql and return a a JSON string of array
func QueryToArrayJSON[T DB](db T, theCase int, sqlStatement string, sqlParams ...any) (string, error) {
	headers, data, err := QueryToArray(db, theCase, sqlStatement, sqlParams...)
	result := map[string]any{
		"headers": headers,
		"data":    data,
	}
	jsonString, err := json.Marshal(result)
	return string(jsonString), err
}

// QueryToMapJSON - run the sql and return a JSON string of array of objects.
func QueryToMapJSON[T DB](db T, theCase int, sqlStatement string, sqlParams ...any) (string, error) {
	data, err := QueryToMap(db, theCase, sqlStatement, sqlParams...)
	jsonString, err := json.Marshal(data)
	return string(jsonString), err
}

// QueryToArray - headers, data, error
func QueryToArray[T DB](db T, theCase int, sqlStatement string, sqlParams ...any) ([]string, [][]string, error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	data := [][]string{}
	rows, err := db.Query(sqlStatement, sqlParams...)
	if err != nil {
		fmt.Println("Error executing: ", sqlStatement)
		return []string{}, data, err
	}
	cols, _ := rows.Columns()
	lenCols := len(cols)
	for i, v := range cols {
		if theCase == Lower {
			cols[i] = strings.ToLower(v)
		} else if theCase == Upper {
			cols[i] = strings.ToUpper(v)
		} else if theCase == Camel {
			cols[i] = toCamel(v)
		}
	}

	rawResult := make([][]byte, lenCols)

	dest := make([]any, lenCols) // A temporary any slice
	for i := range rawResult {
		dest[i] = &rawResult[i] // Put pointers to each string in the interface slice
	}

	for rows.Next() {
		result := make([]string, lenCols)
		rows.Scan(dest...)
		for i, raw := range rawResult {
			if raw == nil {
				result[i] = ""
			} else {
				result[i] = string(raw)
			}
		}
		data = append(data, result)
	}
	return cols, data, nil
}

// QueryToMap - run sql and return an array of maps
func QueryToMap[T DB](db T, theCase int, sqlStatement string, sqlParams ...any) ([]map[string]string, error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	results := []map[string]string{}
	rows, err := db.Query(sqlStatement, sqlParams...)
	if err != nil {
		fmt.Println("Error executing: ", sqlStatement)
		return results, err
	}
	cols, _ := rows.Columns()
	lenCols := len(cols)

	for i, v := range cols {
		if theCase == Lower {
			cols[i] = strings.ToLower(v)
		} else if theCase == Upper {
			cols[i] = strings.ToUpper(v)
		} else if theCase == Camel {
			cols[i] = toCamel(v)
		}
	}

	rawResult := make([][]byte, lenCols)

	dest := make([]any, lenCols) // A temporary any slice
	for i := range rawResult {
		dest[i] = &rawResult[i] // Put pointers to each string in the interface slice
	}

	for rows.Next() {
		result := make(map[string]string, lenCols)
		rows.Scan(dest...)
		for i, raw := range rawResult {
			if raw == nil {
				result[cols[i]] = ""
			} else {
				result[cols[i]] = string(raw)
			}
		}
		results = append(results, result)
	}
	return results, nil
}

// Exec - run the sql and returns rows affected.
func Exec[T DB](db T, sqlStatement string, sqlParams ...any) (int64, error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	result, err := db.Exec(sqlStatement, sqlParams...)
	if err != nil {
		fmt.Println("Error executing: ", sqlStatement)
		fmt.Println(err)
		return 0, err
	}
	return result.RowsAffected()
}

func toCamel(s string) (ret string) {
	s = strings.ToLower(s)
	a := strings.Split(s, "_")
	for i, v := range a {
		if i == 0 {
			ret += v
		} else {
			f := strings.ToUpper(string(v[0]))
			n := string(v[1:])
			ret += fmt.Sprint(f, n)
		}
	}
	return
}
