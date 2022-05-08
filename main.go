package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"search-engine/app/midwares"
	"search-engine/config/database"
	"search-engine/config/router"
)

func main() {
	database.Init()

	r := gin.Default()
	r.Use(cors.Default())
	r.Use(midwares.ErrHandler())
	r.NoMethod(midwares.HandleNotFound)
	r.NoRoute(midwares.HandleNotFound)

	router.Init(r)

	err := r.Run()
	if err != nil {
		log.Fatal("ServerStartFailed", err)
	}
}
