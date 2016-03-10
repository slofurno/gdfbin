package main

import (
	"testing"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var testStore *DataStore

var content []byte
var keys []string

func init() {
	content = []byte(`
		hey this is some content
		<!DOCTYPE html5>
		<html>
		<head>
		<style>
		<body>
		<br>,br><br>br<rbrbdfgfdgdfg
		DSFDSFDSFDSFDS
		:D
		/`)

	var db *sql.DB
	db, _ = sql.Open("sqlite3", "./pastes.db")

	db.Exec("PRAGMA synchronous = OFF; PRAGMA journal_mode = MEMORY;")

	testStore = NewDataStore(db)

	keys = make([]string, 1024)

	for i, _ := range keys {

		paste := NewPaste()
		paste.Content = content

		pastes := testStore.Pastes2

		for i := 0; i < len(pastes); i++ {
			pastes[i].Insert(paste)
		}

		keys[i] = paste.Id
	}

}

func BenchmarkMemcachedGet(b *testing.B) {
	for _, key := range keys {
		_, err := testStore.Pastes2[0].Get(key)

		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkSqliteGet(b *testing.B) {
	for _, key := range keys {
		_, err := testStore.Pastes2[1].Get(key)

		if err != nil {
			b.Error(err)
		}
	}
}
