package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"testing"
	"go.uber.org/goleak"
	"os"
)


func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func TestNotCloseDBQuery(t *testing.T) {
	defer goleak.VerifyNoLeaks(t)

	// delete db file if exist
	if _, err := os.Stat("./foo.db"); err == nil {
		os.Remove("./foo.db")
	}

	db, err := sql.Open("sqlite3", "./foo.db")
	checkErr(err)
	defer db.Close()

	_, err = db.Exec("CREATE TABLE userinfo( id integer, username varchar(32) )")
	checkErr(err)

	// insert
	stmt, err := db.Prepare("INSERT INTO userinfo(id, username) values(?,?)")
	checkErr(err)

	res, err := stmt.Exec(1, "aaaa")
	checkErr(err)

	id, err := res.LastInsertId()
	checkErr(err)

	fmt.Println(id)

	// update
	stmt, err = db.Prepare("update userinfo set username=? where id=?")
	checkErr(err)

	res, err = stmt.Exec("astaxieupdate", id)
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)

	fmt.Println(affect)

	// query
	rows, err := db.Query("SELECT * FROM userinfo")
	checkErr(err)
	cols, err := rows.Columns()
	fmt.Println(len(cols))
	// forget to close

}


func TestCloseDBQuery(t *testing.T) {
	defer goleak.VerifyNoLeaks(t)

	// delete db file if exist
	if _, err := os.Stat("./foo.db"); err == nil {
		os.Remove("./foo.db")
	}

	db, err := sql.Open("sqlite3", "./foo.db")
	checkErr(err)
	defer db.Close()

	_, err = db.Exec("CREATE TABLE userinfo( id integer, username varchar(32) )")
	checkErr(err)

	// insert
	stmt, err := db.Prepare("INSERT INTO userinfo(id, username) values(?,?)")
	checkErr(err)

	res, err := stmt.Exec(1, "aaaa")
	checkErr(err)

	id, err := res.LastInsertId()
	checkErr(err)

	fmt.Println(id)

	// update
	stmt, err = db.Prepare("update userinfo set username=? where id=?")
	checkErr(err)

	res, err = stmt.Exec("astaxieupdate", id)
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)

	fmt.Println(affect)

	// query
	rows, err := db.Query("SELECT * FROM userinfo")
	checkErr(err)
	cols, err := rows.Columns()
	fmt.Println(len(cols))
	// forget to close
	rows.Close()

}