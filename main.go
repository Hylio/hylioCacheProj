package main

import (
	"fmt"
	"log"
	"net/http"
)

//func serve(c *gin.Context) {
//	log.Println(c.Request.URL.Path)
//	c.String(http.StatusOK, "hello, world")
//}
//
//func main() {
//	r := gin.Default()
//	r.GET("/", serve)
//
//	err := r.Run("localhost:8080")
//	if err != nil {
//		log.Fatal("Failed to start server")
//	}
//}

var db = map[string]string{
	"zhanghao":    "hylio",
	"wangrui":     "civet",
	"zhouruqiang": "dio",
}

func main() {
	NewGroup("scores", 2<<10, GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[DB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))

	addr := "localhost:9999"
	peers := NewHTTPPool(addr)
	log.Println("geecache is running at", addr)
	log.Fatal(http.ListenAndServe(addr, peers))
}
