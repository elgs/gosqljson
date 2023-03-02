# gosqljson

A Go library to work with SQL database using standard `database/sql` api. It provides a set of functions to convert query result to array, map and struct.

# Installation

`go get -u github.com/elgs/gosqljson`

# Sample code

Please note all the `err`s are ignored for brevity in the following code. You should always check the `err` returned.

```go
package main

import (
	"database/sql"
	"fmt"

	"github.com/elgs/gosqljson"
	_ "modernc.org/sqlite"
)

type User struct {
	Id   int    `db:"id"`
	Name string `db:"name"`
}

func main() {
	db, _ := sql.Open("sqlite", ":memory:")

	result, _ := gosqljson.Exec(db, "CREATE TABLE test (ID INTEGER PRIMARY KEY, NAME TEXT)")
	fmt.Printf("result: %+v\n", result)
	// result: map[last_insert_id:0 rows_affected:0]

	result, _ = gosqljson.Exec(db, "INSERT INTO test (ID, NAME) VALUES (?, ?)", 1, "Alpha")
	fmt.Printf("result: %+v\n", result)
	// result: map[last_insert_id:1 rows_affected:1]

	result, _ = gosqljson.Exec(db, "INSERT INTO test (ID, NAME) VALUES (?, ?)", 2, "Beta")
	fmt.Printf("result: %+v\n", result)
	// result: map[last_insert_id:2 rows_affected:1]

	result, _ = gosqljson.Exec(db, "INSERT INTO test (ID, NAME) VALUES (?, ?)", 3, "Gamma")
	fmt.Printf("result: %+v\n", result)
	// result: map[last_insert_id:3 rows_affected:1]

	cols, resultArray, _ := gosqljson.QueryToArray(db, gosqljson.AsIs, "SELECT * FROM test WHERE ID > ?", 1)
	fmt.Printf("cols: %+v\n", cols)         // cols: [ID NAME]
	fmt.Printf("array: %+v\n", resultArray) // array: [[2 Beta] [3 Gamma]]

	resultMap, _ := gosqljson.QueryToMap(db, gosqljson.AsIs, "SELECT * FROM test WHERE ID < ?", 3)
	fmt.Printf("map: %+v\n", resultMap)
	// map: [map[ID:1 NAME:Alpha] map[ID:2 NAME:Beta]]

	resultStructs := []User{}
	_ = gosqljson.QueryToStruct(db, &resultStructs, "SELECT  NAME,ID FROM test WHERE ID > ?", 0)
	fmt.Printf("structs: %+v\n", resultStructs)
	// structs: [{Id:1 Name:Alpha} {Id:2 Name:Beta} {Id:3 Name:Gamma}]
}
```
