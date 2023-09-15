package main

import (
	"context"
	"crypto/sha256"
	"fmt"
	"net/http"

	"github.com/redis/go-redis/v9"
)

type Node struct {
	url       string
	redirects int
}

func NewNode(url string) Node {
	return Node{url, 0}
}

func Mapping(w http.ResponseWriter, r *http.Request) {
	tmp := r.URL.Query()
	url := tmp["url"][0]
	if url == "" {
		w.Write([]byte("Bad Request\n"))
	} else {
		hash := fmt.Sprintf("%x", sha256.Sum256([]byte(url)))[:8]
		db.Set(context.TODO(), hash, url, 0)
		cache[hash] = NewNode(url)
	}
}

func Redirect(w http.ResponseWriter, r *http.Request) {
	hash := r.RequestURI[3:len(r.RequestURI)]
	w.Write([]byte(cache[hash].url))
	cache[hash].redirects++
}

var db *redis.Client
var cache map[string]Node

func main() {
	db = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	cache = make(map[string]Node, 100)
	http.HandleFunc("/a/", Mapping)
	http.HandleFunc("/s/", Redirect)
	http.ListenAndServe(":8080", nil)
}
