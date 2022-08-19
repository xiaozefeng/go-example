package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
)
import "github.com/gin-gonic/gin"

func main() {
	r := gin.New()
	go func() {
		log.Fatal(http.ListenAndServe(":6060", nil))
	}()
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello Gin\n")
	})
	log.Fatal(r.Run(":8080"))
}
