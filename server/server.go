package server

import (
	"encoding/json"
	"github.com/gophergala2016/gophersiesta/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"github.com/gophergala2016/gophersiesta/Godeps/_workspace/src/github.com/spf13/viper"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"github.com/gophergala2016/gophersiesta/server/storage"
	"github.com/gophergala2016/gophersiesta/server/placeholders"
	"fmt"
)

var storage server.Storage

type Server struct {
	Storage server.Storage
	*gin.Engine
}

func StartServer() *Server {

	storage = &server.BoltDb{} // RAM

	storage.Init()

	server.CreateSampleData(storage)

	router := gin.Default()

	server := &Server{storage, router}

	// This handler will match /conf/appname but will not match neither /conf/ or /conf
	server.GET("/conf/:appname", func(c *gin.Context) {
		name := c.Param("appname")
		myViper, err := readTemplate(name)
		if err != nil {
			c.String(http.StatusNotFound, "Config file for %s not found\n", name)
		} else {
			filename := myViper.ConfigFileUsed()
			c.String(http.StatusOK, safeFileRead(filename)+"\n")

		}
	})

	// Return list of placeholders
	server.GET("/conf/:appname/placeholders", func(c *gin.Context) {
		name := c.Param("appname")
		myViper, err := readTemplate(name)
		if err != nil {
			c.String(http.StatusNotFound, "Config file for %s not found\n", name)
		} else {
			properties := placeholders.GetPlaceHolders(myViper)
			propsJson, _ := json.Marshal(properties)
			c.String(http.StatusOK, string(propsJson))
		}
	})

	// Return list of set values
	server.GET("/conf/:appname/values", func(c *gin.Context) {
		name := c.Param("appname")
		labels := c.DefaultQuery("labels", "default")

		list := make(map[string]string)
		if strings.Contains(labels, ",") {
			lbls := strings.Split(labels, ",")
			// MERGE values of different labels, last overrides current value
			for _, label := range lbls {
				l := storage.GetOptions(name, label)
				for k, v := range l {
					list[k] = v
				}
			}
		} else {
			list = storage.GetOptions(name, labels)
		}

		vs := placeholders.CreateValues(list);

		vs_json, _ := json.Marshal(vs)
		c.String(http.StatusOK, string(vs_json))
	})

	// Return list of set values
	server.POST("/conf/:appname/values", func(c *gin.Context) {
		name := c.Param("appname")
		labels := c.DefaultQuery("labels", "default")

		body := c.Request.Body
		x, err := ioutil.ReadAll(body)

		if err!=nil {
			c.String(http.StatusBadRequest, "Bad request")
		} else {

			data := map[string]interface{}{}
			json.Unmarshal(x, &data)

			lbls := strings.Split(labels, ",")
			for _, label := range lbls {
				for k, v := range data {
					storage.SetOption(name, label, k, fmt.Sprint(v))
				}
			}


			c.String(http.StatusOK, "Ok")
		}

	})

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

func readTemplate(appname string) (*viper.Viper, error) {

	aux := viper.New()
	aux.SetConfigName("config")
	aux.AddConfigPath("apps/" + appname + "/")

	err := aux.ReadInConfig()
	return aux, err

}

func safeFileRead(filename string) string {
	fileContent, errFile := ioutil.ReadFile(filename)
	if errFile != nil {
		fileContent = []byte("")
	}
	return string(fileContent)
}


/*func readJSON(str string) map[string]string {
	data := map[string]interface{}{}
	json.Unmarshal(dataStream, &data)
	return
}*/
