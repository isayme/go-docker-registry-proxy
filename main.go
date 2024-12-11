package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/isayme/go-docker-registry-proxy/src"
)

func main() {
	conf := src.GetConfig()

	r := gin.Default()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/__version__", src.TokenHandler)
	r.GET("/__token__", src.TokenHandler)

	r.Any("/v2/*path", src.V2Handler)

	address := conf.Server.Addr
	log.Printf("listen at '%s'", address)
	r.Run(address)
}
