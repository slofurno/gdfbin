package main

import (
	"bytes"
	"crypto/rand"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"net/http"
)

var pastebin *sql.DB

func randomHash() string {

	var b = make([]byte, 5)
	rand.Read(b)
	return crock32(b)
}

func pasteHandler(res http.ResponseWriter, req *http.Request) {

	switch req.Method {
	case "GET":
		path := req.URL.Path[1:]

		if len(path) <= 1 {
			res.Write([]byte("<!DOCTYPE html><meta charset=\"utf-8\"><pre>cat main.go | curl --data-binary @- https://gdf3.com</pre>"))
			return
		}

		rows, err := pastebin.Query("select content from pastes WHERE id = ?", path)

		if err != nil {
			fmt.Println(err)
			res.WriteHeader(500)
			return
		}

		for rows.Next() {
			var content []byte
			rows.Scan(&content)
			res.Header().Add("Content-Type", "text/plain; charset=utf-8")
			res.Write(content)
		}

		break
	case "POST":

		buf := bytes.NewBuffer(nil)
		_, err := io.Copy(buf, req.Body)

		if err != nil {
			fmt.Println(err)
			res.WriteHeader(500)
			return
		}

		hash := randomHash()

		_, err = pastebin.Exec("insert into pastes values (?,?)", hash, buf.Bytes())

		if err != nil {
			fmt.Println(err)
			res.WriteHeader(500)
			return
		}

		res.Write([]byte("https://gdf3.com/" + hash + "\n"))
		break
	}

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
	pastebin, err = sql.Open("sqlite3", "./pastes.db")

	if err != nil {
		fmt.Println(err)
	}

	http.HandleFunc("/", pasteHandler)
	err = http.ListenAndServe(":666", nil)

	if err != nil {
		fmt.Println(err)
	}

}
