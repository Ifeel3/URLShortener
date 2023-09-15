package main

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math"
	"net/http"

	"github.com/redis/go-redis/v9"
)

type RedisWithCache struct {
	cache     map[string]string
	redirects map[string]uint
	db        *redis.Client
}

func NewCache(db *redis.Client) RedisWithCache {
	return RedisWithCache{make(map[string]string, 100), make(map[string]uint, 100), db}
}

func (rwc *RedisWithCache) deleteMinRedirects() {
	var min uint = math.MaxUint
	for _, val := range rwc.redirects {
		if val < min {
			min = val
		}
	}
	for key, val := range rwc.redirects {
		if val == min {
			delete(rwc.cache, key)
			delete(rwc.redirects, key)
		}
	}
}

func (rwc *RedisWithCache) Set(key string, value string) {
	if key == "" || value == "" {
		return
	}
	if len(rwc.redirects) > 100 {
		rwc.deleteMinRedirects()
	}
	rwc.cache[key] = value
	rwc.redirects[key] = 0
	rwc.db.Set(context.TODO(), key, value, 0)
}

func (rwc *RedisWithCache) Get(key string) string {
	res, ok := rwc.cache[key]
	if !ok {
		tmp, err := rwc.db.Get(context.TODO(), key).Result()
		if err != nil {
			return ""
		}
		rwc.cache[key] = tmp
		rwc.redirects[key] = 0
		res = tmp
	}
	rwc.redirects[key]++
	return res
}

func Mapping(w http.ResponseWriter, r *http.Request) {
	tmp := r.URL.Query()
	url := tmp["url"][0]
	if url == "" {
		w.Write([]byte("Bad Request\n"))
	} else {
		hash := fmt.Sprintf("%x", sha256.Sum256([]byte(url)))[:8]
		cache.Set(hash, url)
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"hash": hash})
	}
}

func Redirect(w http.ResponseWriter, r *http.Request) {
	hash := r.RequestURI[3:len(r.RequestURI)]
	url := cache.Get(hash)
	if url == "" {
		w.Write([]byte("Not Found"))
		return
	}
	http.Redirect(w, r, url, http.StatusMovedPermanently)
}

var cache RedisWithCache

func main() {
	cache = NewCache(redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	}))
	http.HandleFunc("/a/", Mapping)
	http.HandleFunc("/s/", Redirect)
	http.ListenAndServe(":8080", nil)
}
