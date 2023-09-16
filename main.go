package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/pablorc/shrtnr/internal/redis"
)

func getRoot(w http.ResponseWriter, r *http.Request, red redis.Connection) {
	url := r.URL.Path
	fmt.Printf("got %s request\n", url)
	if url == "/" {
		io.WriteString(w, "This is my website!\n")
	} else {
		redirect(w, r, url[1:], red)
	}
}

func redirect(w http.ResponseWriter, r *http.Request, path string, red redis.Connection) {
	url, err := redirectionFor(path, red)
	if err != nil {
		fmt.Printf("ERR: %s\n", err.Error())
	} else {
		http.Redirect(w, r, url, http.StatusSeeOther)
	}
}

func redirectionFor(path string, red redis.Connection) (string, error) {
	url, err := red.Get(path)
	if err != nil {
		panic(err)
	}

	if url == "" {
		return "", errors.New("No redirection")
	}

	return url, nil
}

func main() {
	fmt.Println("Start")

	red, rerr := redis.Connect()
	if rerr != nil {
		fmt.Printf("ERR: %s", rerr.Error())
		os.Exit(1)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		getRoot(w, r, red)
	})

	err := http.ListenAndServe(":3333", nil)

	if err != nil {
		fmt.Printf("ERR: %s\n", err)
	}
}
