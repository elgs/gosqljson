package gosqljson

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var Version = "3"

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

// QueryToArrays - run sql and return an array of arrays
func QueryToArrays[T DB](db T, theCase int, sqlStatement string, sqlParams ...any) ([]string, [][]any, error) {
	data := [][]any{}
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

	rawResult := make([]any, lenCols)
	colTypes, _ := rows.ColumnTypes()
	dest := make([]any, lenCols) // A temporary any slice
	for i := range rawResult {
		dest[i] = &rawResult[i] // Put pointers to each string in the interface slice
	}
	for rows.Next() {
		result := make([]any, lenCols)
		rows.Scan(dest...)
		for i, raw := range rawResult {
			// faulty mysql driver workaround https://github.com/go-sql-driver/mysql/issues/1401
			if v, ok := raw.([]byte); ok {
				value := string(v)
				switch colTypes[i].DatabaseTypeName() {
				case "SMALLINT", "MEDIUMINT", "INT", "INTEGER", "BIGINT", "YEAR":
					raw, _ = strconv.Atoi(value)
				case "TINYINT", "BOOL", "BOOLEAN", "BIT":
					raw, _ = strconv.ParseBool(value)
				case "FLOAT", "DOUBLE", "DECIMAL":
					raw, _ = strconv.ParseFloat(value, 64)
				case "DATETIME", "TIMESTAMP":
					raw, _ = time.Parse("2006-01-02 15:04:05", value)
				case "DATE":
					raw, _ = time.Parse("2006-01-02", value)
				case "TIME":
					raw, _ = time.Parse("15:04:05", value)
				case "NULL":
					raw = nil
				}
				raw = value
			}
			result[i] = raw
		}
		data = append(data, result)
	}
	return cols, data, nil
}

// QueryToMaps - run sql and return an array of maps
func QueryToMaps[T DB](db T, theCase int, sqlStatement string, sqlParams ...any) ([]map[string]any, error) {
	results := []map[string]any{}
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

	rawResult := make([]any, lenCols)
	colTypes, _ := rows.ColumnTypes()
	dest := make([]any, lenCols) // A temporary any slice
	for i := range rawResult {
		dest[i] = &rawResult[i] // Put pointers to each string in the interface slice
	}
	for rows.Next() {
		result := make(map[string]any, lenCols)
		rows.Scan(dest...)
		for i, raw := range rawResult {
			// faulty mysql driver workaround https://github.com/go-sql-driver/mysql/issues/1401
			if v, ok := raw.([]byte); ok {
				value := string(v)
				switch colTypes[i].DatabaseTypeName() {
				case "SMALLINT", "MEDIUMINT", "INT", "INTEGER", "BIGINT", "YEAR":
					raw, _ = strconv.Atoi(value)
				case "TINYINT", "BOOL", "BOOLEAN", "BIT":
					raw, _ = strconv.ParseBool(value)
				case "FLOAT", "DOUBLE", "DECIMAL":
					raw, _ = strconv.ParseFloat(value, 64)
				case "DATETIME", "TIMESTAMP":
					raw, _ = time.Parse("2006-01-02 15:04:05", value)
				case "DATE":
					raw, _ = time.Parse("2006-01-02", value)
				case "TIME":
					raw, _ = time.Parse("15:04:05", value)
				case "NULL":
					raw = nil
				}
				raw = value
			}
			result[cols[i]] = raw
		}
		results = append(results, result)
	}
	return results, nil
}

func QueryToStructs[T DB, S any](db T, results *[]S, sqlStatement string, sqlParams ...any) error {
	rows, err := db.Query(sqlStatement, sqlParams...)
	if err != nil {
		fmt.Println("Error executing: ", sqlStatement)
		return err
	}
	cols, _ := rows.Columns()
	lenCols := len(cols)

	for rows.Next() { // iterate through rows
		colValues := make([]any, lenCols)
		var result S
		structValue := reflect.ValueOf(&result).Elem()
		for colIndex, colName := range cols { // iterate through columns
			found := false
			for fieldIndex := 0; fieldIndex < structValue.NumField(); fieldIndex++ { // iterate through struct fields
				field := structValue.Type().Field(fieldIndex)
				fieldTag := field.Tag.Get("db")
				if strings.EqualFold(colName, fieldTag) {
					colValues[colIndex] = structValue.Field(fieldIndex).Addr().Interface()
					found = true
					break
				}
			}
			if !found {
				colValues[colIndex] = new(any)
			}
		}
		rows.Scan(colValues...)
		*results = append(*results, result)
	}

	return nil
}

// Exec - run sql and return the number of rows affected
func Exec[T DB](db T, sqlStatement string, sqlParams ...any) (map[string]int64, error) {
	result, err := db.Exec(sqlStatement, sqlParams...)
	if err != nil {
		fmt.Println("Error executing: ", sqlStatement)
		fmt.Println(err)
		return nil, err
	}
	rowsffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	ret := map[string]int64{
		"rows_affected": rowsffected,
	}
	lastInsertId, err := result.LastInsertId()
	if err == nil {
		ret["last_insert_id"] = lastInsertId
	}
	return ret, nil
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
