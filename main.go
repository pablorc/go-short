package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

func getRoot(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path
	fmt.Printf("got %s request\n", url)
	if url == "/" {
		io.WriteString(w, "This is my website!\n")
	} else {
		redirect(w, r, url[1:])
	}
}

func redirect(w http.ResponseWriter, r *http.Request, path string) {
	url, err := redirectionFor(path)
	if err != nil {
		fmt.Printf("ERR: %s\n", err.Error())
	} else {
		http.Redirect(w, r, url, http.StatusSeeOther)
	}
}

func redirectionFor(path string) (string, error) {
	reds := make(map[string]string)
	reds["001"] = "https://google.com"
	reds["002"] = "https://wikipedia.org"

	if reds[path] == "" {
		return "", errors.New("No redirection")
	}

	return reds[path], nil
}

func main() {
	fmt.Println("Start")
	http.HandleFunc("/", getRoot)

	err := http.ListenAndServe(":3333", nil)

	if err != nil {
		fmt.Printf("ERR: %s\n", err)
	}
}
