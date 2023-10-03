package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/pablorc/go-short/internal/keygen"
	"github.com/pablorc/go-short/internal/redis"
)

func getForm(w http.ResponseWriter, r *http.Request, red redis.Connection) {
	username := os.Getenv("HTTP_USERNAME")
	password := os.Getenv("HTTP_PASSWORD")

	u, p, ok := r.BasicAuth()
	if !ok {
		fmt.Println("Error parsing basic auth")
		w.WriteHeader(401)
		return
	}
	if u != username || p != password {
		fmt.Println("Invalid credencials")
		w.WriteHeader(401)
		return
	}

	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}

	rawUrl := r.FormValue("url")
	if rawUrl == "" {
		fmt.Println("ERR: Saving without URL")
		return
	}

	url, uerr := url.Parse(rawUrl)
	if uerr != nil {
		panic(uerr)
	}

	key := keygen.NewKey()
	red.Set(key, *url)

	fmt.Printf("Saved new URL %s: %s\n", url.String(), key)
	fmt.Fprintf(w, "http://localhost:3333/%s\n", key)
}

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
		http.Redirect(w, r, url.String(), http.StatusSeeOther)
	}
}

func redirectionFor(path string, red redis.Connection) (url.URL, error) {
	url, err := red.Get(path)
	if err != nil {
		panic(err)
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

	http.HandleFunc("/save", func(w http.ResponseWriter, r *http.Request) {
		getForm(w, r, red)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		getRoot(w, r, red)
	})

	err := http.ListenAndServe(":3333", nil)

	if err != nil {
		fmt.Printf("ERR: %s\n", err)
	}
}
