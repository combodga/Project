package main

import (
	"crypto"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var pairs = make(map[string]string)

func main() {
	const host string = ""
	const port string = "8080"

	http.HandleFunc("/", mainHandler)

	fmt.Println("The service works on", host, ":", port)

	log.Fatal(http.ListenAndServe(host+":"+port, nil))
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if path == "" {
		path = "/"
	}

	id := path[1:]

	url, ok := getUrl(id)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte("error, there is no such link"))
		return
	}

	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	id := string(body)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	if len(id) > 2048 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error, the link cannot be longer than 2048 characters"))
		return
	} else if len(id) < 7 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error, the link cannot be shortened"))
		return
	}

	url, ok := getUrl(id)
	var err error
	if !ok {
		url, err = shortener(id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("error, failed to create a shortened URL"))
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(url))
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodGet:
		getHandler(w, r)

	case http.MethodPost:
		postHandler(w, r)

	default:
		http.Error(w, "error, you can only use the get and post methods", http.StatusMethodNotAllowed)

	}
}

func getUrl(id string) (string, bool) {
	if len(id) <= 0 {
		return "", false
	}
	url, ok := pairs[id]
	if !ok {
		return "", false
	}

	return url, true
}

func shortener(s string) (string, error) {
	h := crypto.MD5.New()
	if _, err := h.Write([]byte(s)); err != nil {
		return "", fmt.Errorf("abbreviation error URL: %v", err)
	}
	hash := string(h.Sum([]byte{}))
	hash = hash[len(hash)-5:]
	id := base64.StdEncoding.EncodeToString([]byte(hash))
	id = strings.ToLower(id)[:len(id)-1]
	id = strings.ReplaceAll(id, "/", "")
	id = strings.ReplaceAll(id, "=", "")

	pairs[id] = s

	return id, nil
}
