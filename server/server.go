package server

import (
	"fmt"
	"github.com/gophergala2016/gophersiesta/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"github.com/gophergala2016/gophersiesta/Godeps/_workspace/src/github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"io/ioutil"
	"encoding/json"
	"strings"
	"github.com/gophergala2016/gophersiesta/server/storage"
)

var storage server.Storage

func StartServer() {

	storage = &server.Ethereal{}
	storage.Init()



	storage.SetOption("app1", "", "admin", "GOPHER")
	storage.SetOption("app1", "dev", "admin", "GOPHER-dev")
	storage.SetOption("app1", "prod", "admin", "GOPHER-prod")
	storage.SetOption("app1", "", "password", "FOOBAR")
	storage.SetOption("app1", "dev", "password", "LOREM")
	storage.SetOption("app1", "prod", "password", "IPSUM")

	storage.SetOption("app2", "", "password", "DOCKER-PASS")
	storage.SetOption("app2", "dev", "password", "DEV-PASS")

	router := gin.Default()

	// This handler will match /conf/appname but will not match neither /conf/ or /conf
	router.GET("/conf/:appname", func(c *gin.Context) {
		name := c.Param("appname")
		myViper, err := readTemplate(name)
		if err != nil {
			c.String(http.StatusNotFound, "Config file for %s not found\n", name)
		} else {
			filename := myViper.ConfigFileUsed()
			c.String(http.StatusOK, safeFileRead(filename) + "\n")

		}
	})

	// Return list of placeholders
	router.GET("/conf/:appname/placeholders", func(c *gin.Context) {
		name := c.Param("appname")
		myViper, err := readTemplate(name)
		if err != nil {
			c.String(http.StatusNotFound, "Config file for %s not found\n", name)
		} else {
			list := getPlaceHolders(myViper)
			list_json, _ := json.Marshal(list)
			c.String(http.StatusOK, string(list_json))
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

	// However, this one will match /conf/app1/ and also /conf/app1/send
	// If no other routers match /conf/app1, it will redirect to /conf/app1/
	/*router.GET("/conf/:appname/*action", func(c *gin.Context) {
		name := c.Param("appname")
		action := c.Param("action")
		message := name + " is " + action + "\n"
		c.String(http.StatusOK, message)
	})        */

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

func getPlaceHolders(conf *viper.Viper) map[string]string {
	list := parseMap(conf.AllSettings())
	return list
}

func parseMap(aMap map[string]interface{}) map[string]string {
	list := make(map[string]string)
	for key, value := range aMap {
		switch v := value.(type) {
		case map[interface{}]interface{}:
			l := parseMapInterface(v)
			for pkey, pvalue := range l {
				list[pkey] = pvalue
			}
		case string:
			if v[:2] == "${" {
				list[key] = v
			}
		default:
		}
	}
	return list
}

func parseMapInterface(aMap map[interface{}]interface{}) map[string]string {
	list := make(map[string]string)
	for key, value := range aMap {
		switch v := value.(type) {
		case map[interface{}]interface{}:
			l := parseMapInterface(v)
			for pkey, pvalue := range l {
				list[pkey] = pvalue
			}
		case string:
			if v[:2] == "${" {
				keystr := fmt.Sprint(key) // <-- HACK
				list[keystr] = v
			}
		default:
		}
	}
	return list
}

