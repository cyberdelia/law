package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", "postgres:///law_test?replication=database&sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	rows, err := db.Query("BASE_BACKUP")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var pos, tli string
		if err := rows.Scan(&pos, &tli); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("pos %s, tli %s\n", pos, tli)
	}
	if !rows.NextResultSet() {
		log.Println("no more result set (1)")
	}
	for rows.Next() {
		var oid, loc, size sql.NullString
		if err := rows.Scan(&oid, &loc, &size); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("oid %s, loc %s, size %s\n", oid.String, loc.String, size.String)
	}
	if !rows.NextResultSet() {
		log.Println("no more result set (2)")
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}
