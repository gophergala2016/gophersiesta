package main

import (
	"github.com/gophergala2016/gophersiesta/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"net/http"
	"os"
	"log"
)

func main() {

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.Default()

	// This handler will match /conf/appname but will not match neither /conf/ or /conf
	router.GET("/conf/:appname", func(c *gin.Context) {
		name := c.Param("appname")
		c.String(http.StatusOK, "Config file  %s \n", name)
	})

	// However, this one will match /conf/app1/ and also /conf/app1/send
	// If no other routers match /conf/app1, it will redirect to /conf/app1/
	router.GET("/conf/:appname/*action", func(c *gin.Context) {
		name := c.Param("appname")
		action := c.Param("action")
		message := name + " is " + action + "\n"
		c.String(http.StatusOK, message)
	})

	router.Run(":" + port)
}
