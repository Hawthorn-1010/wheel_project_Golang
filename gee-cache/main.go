package main

import (
	"fmt"
	"gee-cache/geecache"
	"log"
	"net/http"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func main() {
	geecache.NewGroup("scores", 2<<10, geecache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))

	addr := "localhost:4101"
	peers := geecache.NewHTTPPool(addr)
	log.Println("geecache is running at", addr)
	// 任何实现了ServeHTTP的方法都可以作为handler
	log.Fatal(http.ListenAndServe(addr, peers))
}
