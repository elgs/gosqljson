gosqljson
=========

A Go SQL to JSON library.

#Installation
`go get -u github.com/elgs/gosqljson`

# Sample code
Data in the table:
```
ID	NAME

0	Alicia
1	Brian
2	Cloe
4	Bianca
5	Leo
6	Joy
7	Sam
8	Elgs
```
```go
package main

import (
	"database/sql"
	"fmt"
	"github.com/elgs/gosqljson"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	ds := "username:password@#@tcp(host:3306)/db"
	db, err := sql.Open("mysql", ds)
	defer db.Close()

	if err != nil {
		fmt.Println("sql.Open:", err)
	}

	a := gosqljson.QueryDbToArrayJson(db, "SELECT ID,NAME FROM t LIMIT ?,?", 0, 7)
	fmt.Println(a)
	// [["0","Alicia"],["1","Brian"],["2","Cloe"],["4","Bianca"],["5","Leo"],["6","Joy"],["7","Sam"]]

	m := gosqljson.QueryDbToMapJson(db, "SELECT ID,NAME FROM t LIMIT ?,?", 0, 7)
	fmt.Println(m)
	// [{"ID":"0","NAME":"Alicia"},{"ID":"1","NAME":"Brian"},{"ID":"2","NAME":"Cloe"},{"ID":"4","NAME":"Bianca"},{"ID":"5","NAME":"Leo"},{"ID":"6","NAME":"Joy"},{"ID":"7","NAME":"Sam"}]

}
```