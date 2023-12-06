package main

import (
	"cache/pkg/cache"
	"cache/pkg/http"
)

func main() {
	c := cache.New("inmemory")
	http.New(c).Listen()
}
