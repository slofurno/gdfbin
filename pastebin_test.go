package main

import (
	"bytes"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"os"
	"testing"
)

var testStore *DataStore
var testAccount *Account

var testBookmarks = []*Bookmark{
	&Bookmark{
		Account: "1",
		Paste:   "paste1",
		Name:    "asdf",
	},
	&Bookmark{
		Account: "1",
		Paste:   "paste2",
		Name:    "asdf",
	},
	&Bookmark{
		Account: "1",
		Paste:   "paste3",
		Name:    "asdf",
	},
}

var testPastes = []*Paste{
	&Paste{
		Id:      "paste1",
		Content: []byte("paste1"),
		Time:    102,
	},
	&Paste{
		Id:      "paste2",
		Content: []byte("paste2"),
		Time:    103,
	},
	&Paste{
		Id:      "paste3",
		Content: []byte("paste3"),
		Time:    104,
	},
}

func die(err error) {
	fmt.Println(err.Error())
	os.Exit(1)
}

func init() {
	var schema []byte
	var err error

	if schema, err = ioutil.ReadFile("./schema.sql"); err != nil {
		die(err)
	}

	db, err := sql.Open("sqlite3", ":memory:")

	if _, err = db.Exec(string(schema)); err != nil {
		die(err)
	}

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	testStore = NewDataStore(db)
	testAccount, err := NewAccount("test_email", "test_password")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	testStore.Accounts.Insert(testAccount)
}

func TestIntegrations(t *testing.T) {
	InsertPastes(t)
	InsertBookmarks(t)
	GetPaste(t)
	GetHistory(t)
	GetLatestBookmarks(t)
}

func GetLatestBookmarks(t *testing.T) {
	bookmarks, err := testStore.Bookmarks.Get(&Account{
		Id: "1",
	})

	if err != nil {
		t.Error(err.Error())
	}

	if len(bookmarks) != 1 {
		t.Error("expected only most recent bookmark")
	}

	if bookmarks[0].Time != 104 {
		t.Error("expected only most recent bookmark")
	}
}

func InsertPastes(t *testing.T) {
	for _, paste := range testPastes {
		testStore.Pastes.Insert(paste)
	}
}

func InsertBookmarks(t *testing.T) {
	for _, bookmark := range testBookmarks {
		testStore.Bookmarks.Insert(bookmark)
	}
}

func GetPaste(t *testing.T) {
	paste := testStore.Bookmarks.GetPaste(&Bookmark{
		Name:    "asdf",
		Account: "1",
	})

	if paste == nil {
		t.Error("expected paste")
	}

	if bytes.Compare(paste.Content, []byte("paste3")) != 0 {
		t.Error("expected most recent paste")
	}
}

func GetHistory(t *testing.T) {
	pastes := testStore.Bookmarks.GetHistory(&Bookmark{
		Name:    "asdf",
		Account: "1",
	})

	if pastes == nil {
		t.Error("expected paste")
	}

	if len(pastes) != 3 {
		t.Error("expected 3 pastes in history")
	}

	if pastes[0].Time != 104 || pastes[1].Time != 103 || pastes[2].Time != 102 {
		t.Error("expected history to be sorted in reverse chronological")
	}
}
