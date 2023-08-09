package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

var db = map[string]string{
	"zhanghao":    "hylio",
	"wangrui":     "civet",
	"zhouruqiang": "dio",
}

func main() {
	NewGroup("aka", 2<<10, GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[DB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))

	addr := "localhost:9999"
	peers := NewHTTPPool(addr)
	log.Println("hyliocache is running at", addr)
	//log.Fatal(http.ListenAndServe(addr, peers))
	r := gin.Default()
	r.GET("_hyliocache/:groupName/:key", peers.Serve)
	log.Fatal(r.Run(addr))
}
