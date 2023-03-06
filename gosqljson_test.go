package gosqljson

import (
	"database/sql"
	"fmt"
	"testing"

	_ "modernc.org/sqlite"
)

type Test struct {
	Id   int    `db:"id"`
	Name string `db:"name"`
}

func TestAll(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	result, err := Exec(db, "CREATE TABLE test (ID INTEGER PRIMARY KEY, NAME TEXT)")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("result: %+v\n", result)

	tx, _ := db.Begin()
	result, err = Exec(tx, "INSERT INTO test (ID, NAME) VALUES (?, ?)", 1, "Alpha")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("result: %+v\n", result)

	result, err = Exec(tx, "INSERT INTO test (ID, NAME) VALUES (?, ?)", 2, "Beta")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("result: %+v\n", result)

	result, err = Exec(tx, "INSERT INTO test (ID, NAME) VALUES (?, ?)", 3, "Gamma")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("result: %+v\n", result)
	tx.Commit()

	cols, resultArray, err := QueryToArrays(db, AsIs, "SELECT * FROM test WHERE ID > ?", 1)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("cols: %+v\n", cols)
	fmt.Printf("arrays: %+v\n", resultArray)

	resultMaps, err := QueryToMaps(db, AsIs, "SELECT * FROM test WHERE ID < ?", 3)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("maps: %+v\n", resultMaps)

	resultStructs := []Test{}
	err = QueryToStructs(db, &resultStructs, "SELECT  NAME,ID FROM test WHERE ID > ?", 0)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("structs: %+v\n", resultStructs)
}
