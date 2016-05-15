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
	res := 0
	var j uint = 0
	out := []byte{}

	for i := 0; i < len(bytes); i++ {
		j += 8
		res <<= 8
		res |= int(bytes[i])

		for j >= 5 {
			j -= 5
			b := (res >> j) & 31
			out = append(out, byte(b))
		}
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

	r.HandleFunc("/bookmarks/{name}/history", getHistory).Methods("GET")
	r.HandleFunc("/bookmarks/{name}", removeBookmark).Methods("DELETE")

	r.HandleFunc("/user", createAccount)
	r.HandleFunc("/", getHome).Methods("GET")

	r.HandleFunc("/login", getToken).Methods("POST")

	r.HandleFunc("/", postPaste).Methods("POST")
	r.HandleFunc("/{paste}", getPaste).Methods("GET")

	err = http.ListenAndServe(":666", r)

	if err != nil {
		fmt.Println(err)
	}

}
