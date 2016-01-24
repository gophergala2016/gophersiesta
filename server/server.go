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
	"github.com/gophergala2016/gophersiesta/server/properties"
)

var storage server.Storage

func StartServer() {

	storage = &server.Ethereal{} // RAM
	//storage = &server.LevelDB{"db/options", nil, nil, nil} // LevelDB

	storage.Init()

	storage.SetOption("app1", "prod", "datasource.url", "jdbc:mysql://proddatabaseserver:3306/shcema?profileSQL=true")
	storage.SetOption("app1", "", "datasource.username", "GOPHER")
	storage.SetOption("app1", "dev", "datasource.username", "GOPHER-dev")
	storage.SetOption("app1", "prod", "datasource.username", "GOPHER-prod")
	storage.SetOption("app1", "", "datasource.password", "FOOBAR")
	storage.SetOption("app1", "dev", "datasource.password", "LOREM")
	storage.SetOption("app1", "prod", "datasource.password", "IPSUM")

	storage.SetOption("app2", "", "datasource.password", "DOCKER-PASS")
	storage.SetOption("app2", "dev", "datasource.password", "DEV-PASS")

	router := gin.Default()

	// This handler will match /conf/appname but will not match neither /conf/ or /conf
	router.GET("/conf/:appname", func(c *gin.Context) {
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
	router.GET("/conf/:appname/placeholders", func(c *gin.Context) {
		name := c.Param("appname")
		myViper, err := readTemplate(name)
		if err != nil {
			c.String(http.StatusNotFound, "Config file for %s not found\n", name)
		} else {
			properties := properties.GetPlaceHolders(myViper)
			propsJson, _ := json.Marshal(properties)
			c.String(http.StatusOK, string(propsJson))
		}
	})

	// Return list of set values
	router.GET("/conf/:appname/values", func(c *gin.Context) {
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
		list_json, _ := json.Marshal(list)
		c.String(http.StatusOK, string(list_json))
	})

	router.Run(getPort())
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

