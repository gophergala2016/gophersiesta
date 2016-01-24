package server

import (
	"github.com/gophergala2016/gophersiesta/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"log"
	"os"
	"github.com/gophergala2016/gophersiesta/server/storage"
	"github.com/gophergala2016/gophersiesta/server/handlers"
)

var db storage.Storage

type Server struct {
	Storage storage.Storage
	*gin.Engine
}

func StartServer() *Server {

	db = &storage.BoltDb{} // RAM

	db.Init()

	storage.CreateSampleData(db)

	router := gin.Default()

	server := &Server{db, router}

	// This handler will match /conf/appname but will not match neither /conf/ or /conf
	server.GET("/conf/:appname", handlers.GetConfig)

	// Return list of placeholders
	server.GET("/conf/:appname/placeholders", handlers.GetPlaceHolders)

	// Return list of set values
	server.GET("/conf/:appname/values", handlers.GetValues(db))

	// Set values
	server.POST("/conf/:appname/values", handlers.SetValues(db))

	// Return the rendered template
	server.GET("/conf/:appname/render/:format", handlers.ReplacePlaceholders(db))

	// Return list of set labels
	server.GET("/conf/:appname/labels", handlers.GetLabels(db))

	server.Run(getPort())

	return server
}

func getPort() string {
	var port = os.Getenv("PORT")
	// Set a default port if there is nothing in the environment
	if port == "" {
		port = "4747"
		log.Println("INFO: No PORT environment variable detected, defaulting to " + port)
	}

	return ":" + port
}

