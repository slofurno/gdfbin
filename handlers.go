package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func createAccount(res http.ResponseWriter, req *http.Request) {
	var err error

	var acc = &struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(acc)

	if err != nil {
		res.Write([]byte(err.Error()))
		return
	}

	account, err := NewAccount(acc.Email, acc.Password)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = store.Accounts.Insert(account)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	login := NewLogin(account)

	err = store.Logins.Insert(login)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	res.Write([]byte(login.Token + "\n"))
}

func createBookmark(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	paste := vars["paste"]
	name := vars["name"]

	token := req.Header.Get("Auth")

	account, err := store.Logins.GetAccount(token)

	if err != nil {
		fmt.Println(err)
		return
	}

	if account == nil {
		res.Write([]byte("no account found"))
		return
	}

	bookmark := &Bookmark{
		Account: account.Id,
		Name:    name,
		Paste:   paste,
	}

	err = store.Bookmarks.Insert(bookmark)

	if err != nil {
		res.Write([]byte(err.Error()))
		return
	}

	res.WriteHeader(200)
}

func getBookmark(res http.ResponseWriter, req *http.Request) {

	token := req.Header.Get("Auth")
	account, err := store.Logins.GetAccount(token)

	if err != nil {
		return
	}

	vars := mux.Vars(req)
	name := vars["name"]

	bookmark := &Bookmark{
		Name:    name,
		Account: account.Id,
	}

	paste := store.Bookmarks.GetPaste(bookmark)

	if paste == nil {
		return
	}

	res.Write(paste.Content)
}

func getBookmarks(res http.ResponseWriter, req *http.Request) {

	token := req.Header.Get("Auth")
	account, err := store.Logins.GetAccount(token)

	if err != nil {
		fmt.Println(err)
		return
	}

	bookmarks, err := store.Bookmarks.Get(account)

	if err != nil {
		fmt.Println(err)
		return
	}

	for _, bookmark := range bookmarks {
		dt := time.Now().Sub(time.Unix(bookmark.Time/1000, 0))
		var modified string
		h := int(dt.Hours())

		if h >= 24 {
			d := int(h / 24)
			modified = strconv.Itoa(d) + " days ago"
		} else {
			modified = strconv.Itoa(h) + " hours ago"
		}

		res.Write([]byte(bookmark.Name + "\t" + bookmark.Paste + "\t" + modified + "\n"))
	}
}

func pasteHandler(res http.ResponseWriter, req *http.Request) {

	switch req.Method {
	case "GET":
		path := req.URL.Path[1:]

		if len(path) <= 1 {
			res.Write([]byte("<!DOCTYPE html><meta charset=\"utf-8\"><pre>cat main.go | curl --data-binary @- https://gdf3.com</pre>"))
			return
		}

		paste, err := store.Pastes.Get(path)

		if err != nil {
			fmt.Println(err)
			res.WriteHeader(500)
			return
		}

		res.Header().Add("Content-Type", "text/plain; charset=utf-8")
		res.Write(paste.Content)

		break
	case "POST":

		buf := bytes.NewBuffer(nil)
		_, err := io.Copy(buf, req.Body)

		if err != nil {
			fmt.Println(err)
			res.WriteHeader(500)
			return
		}

		paste := NewPaste()
		paste.Content = buf.Bytes()
		err = store.Pastes.Insert(paste)

		if err != nil {
			fmt.Println(err)
			res.WriteHeader(500)
			return
		}

		res.Write([]byte("https://gdf3.com/" + paste.Id + "\n"))
		break
	}

}
