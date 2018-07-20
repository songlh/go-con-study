package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"testing"
	"go.uber.org/goleak"
	"os"
)


func _checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func TestNotCloseDBTransaction(t *testing.T) {
	defer goleak.VerifyNoLeaks(t)

	// delete db file if exist
	if _, err := os.Stat("./foo.db"); err == nil {
		os.Remove("./foo.db")
	}

	db, err := sql.Open("sqlite3", "./foo.db")
	_checkErr(err)
	defer db.Close()

	txn, err := db.Begin()
	_checkErr(err)

	if _, err := txn.Query("SELECT * FROM i_do.not_exist"); err == nil {
		t.Fatal("Expected an error but didn't get one")
	}

	//forgot closing the transaction

}

func TestCloseDBTransaction(t *testing.T) {
	defer goleak.VerifyNoLeaks(t)

	// delete db file if exist
	if _, err := os.Stat("./foo.db"); err == nil {
		os.Remove("./foo.db")
	}

	db, err := sql.Open("sqlite3", "./foo.db")
	_checkErr(err)
	defer db.Close()

	txn, err := db.Begin()
	_checkErr(err)

	if _, err := txn.Query("SELECT * FROM i_do.not_exist"); err == nil {
		t.Fatal("Expected an error but didn't get one")
	}

	//close the transaction
	if err := txn.Rollback(); err != nil {
		t.Fatal(err)
	}

}