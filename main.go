package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func serve(c *gin.Context) {
	log.Println(c.Request.URL.Path)
	c.String(http.StatusOK, "hello, world")
}

func main() {
	r := gin.Default()
	r.GET("/", serve)

	err := r.Run("localhost:8080")
	if err != nil {
		log.Fatal("Failed to start server")
	}
}
