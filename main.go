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
		fmt.Println(path)

		rows, err := pastebin.Query("select content from pastes WHERE id = ?", path)

		if err != nil {
			fmt.Println(err)
		}

		for rows.Next() {
			var content string
			rows.Scan(&content)
			res.Write([]byte(content))
		}

		break
	case "POST":

		buf := bytes.NewBuffer(nil)
		_, err := io.Copy(buf, req.Body)

		if err != nil {
			return
		}

		hash := randomHash()

		_, err = pastebin.Exec("insert into pastes values (?,?)", hash, buf.Bytes())

		if err != nil {
			fmt.Println(err)
			return
		}

		res.Write([]byte(hash))

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
	http.ListenAndServe(":666", nil)

}
