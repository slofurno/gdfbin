package main

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

var store *DataStore

func randomHash(n int) string {

	var b = make([]byte, n)
	rand.Read(b)
	return crock32(b)
}

var symbols = "0123456789ABCDEFGHJKMNPQRSTVWXYZ"

func crock32(bytes []byte) string {

	count := len(bytes) * 8
	out := make([]byte, int(count/5))

	var n uint = 0
	var j uint = 0
	var k uint = 0
	var m uint = 0

	for i := 0; i < count; i++ {
		j = uint(i / 8)
		k = uint(i / 5)
		n = uint(i % 5)
		m = uint(i % 8)

		out[k] |= ((bytes[j] >> m) & 1) << n
	}

	for i, _ := range out {
		out[i] = symbols[out[i]]
	}

	return string(out)
}

func main() {

	var err error
	var db *sql.DB
	db, err = sql.Open("sqlite3", "./pastes.db")
	store = NewDataStore(db)

	if err != nil {
		fmt.Println(err)
	}

	r := mux.NewRouter()

	r.HandleFunc("/bookmarks/{paste}/{name}", createBookmark).Methods("POST")
	r.HandleFunc("/bookmarks", getBookmarks).Methods("GET")
	r.HandleFunc("/bookmarks/{name}", getBookmark).Methods("GET")

	r.HandleFunc("/user", createAccount)
	r.HandleFunc("/", pasteHandler).Methods("POST", "GET")
	err = http.ListenAndServe(":666", r)

	if err != nil {
		fmt.Println(err)
	}

}
