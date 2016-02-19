package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func userHandler(res http.ResponseWriter, req *http.Request) {

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
		res.Write(paste.content)

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
		paste.content = buf.Bytes()
		err = store.Pastes.Insert(paste)

		if err != nil {
			fmt.Println(err)
			res.WriteHeader(500)
			return
		}

		res.Write([]byte("https://gdf3.com/" + paste.id + "\n"))
		break
	}

}
