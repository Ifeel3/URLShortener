package main

import (
	"context"
	"net/http"
	"strings"

	"github.com/redis/go-redis/v9"
)

func HandlerSet(w http.ResponseWriter, r *http.Request) {
	splitted := strings.Split(r.RequestURI, "/")
	if len(splitted) != 4 {
		w.Write([]byte("Bad Request"))
	} else {
		err := db.Set(context.TODO(), splitted[2], splitted[3], 0).Err()
		if err != nil {
			w.Write([]byte("Bad Request"))
			return
		}
		w.Write([]byte("OK"))
	}
}

func HandlerGet(w http.ResponseWriter, r *http.Request) {
	splitted := strings.Split(r.RequestURI, "/")
	if len(splitted) != 3 {
		w.Write([]byte("Bad Request"))
	} else {
		response, _ := db.Get(context.TODO(), splitted[2]).Result()
		w.Write([]byte(response))
	}
}

var db *redis.Client

func main() {
	db = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	http.HandleFunc("/set/", HandlerSet)
	http.HandleFunc("/get/", HandlerGet)
	http.ListenAndServe(":8080", nil)
}
