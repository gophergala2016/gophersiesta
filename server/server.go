package server

import (
	"encoding/json"
	"fmt"
	"github.com/gophergala2016/gophersiesta/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"github.com/gophergala2016/gophersiesta/Godeps/_workspace/src/github.com/spf13/viper"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

var options map[string]map[string]map[string]string

type Property struct {
	PropertyName string `json:"property_name"` // the full path to the property datasource.url
	PropertyValue string `json:"property_value"` // ${DATASOURCE_URL:jdbc:mysql://localhost:3306/shcema?profileSQL=true}
	PlaceHolder string `json:"placeholder"`// DATASOURCE_URL
}

type Properties struct {
	Properties []*Property `json:"properties"`
}

func StartServer() {

	options = make(map[string]map[string]map[string]string)

	setOption("app1", "prod", "datasource.url", "jdbc:mysql://proddatabaseserver:3306/shcema?profileSQL=true")
	setOption("app1", "", "datasource.username", "GOPHER")
	setOption("app1", "dev", "datasource.username", "GOPHER-dev")
	setOption("app1", "prod", "datasource.username", "GOPHER-prod")
	setOption("app1", "", "datasource.password", "FOOBAR")
	setOption("app1", "dev", "datasource.password", "LOREM")
	setOption("app1", "prod", "datasource.password", "IPSUM")

	setOption("app2", "", "datasource.password", "DOCKER-PASS")
	setOption("app2", "dev", "datasource.password", "DEV-PASS")

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
			properties := getPlaceHolders(myViper)
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
				l := getOptions(name, label)
				for k, v := range l {
					list[k] = v
				}
			}
		} else {
			list = getOptions(name, labels)
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

func getPlaceHolders(conf *viper.Viper) Properties {
	list := parseMap(conf.AllSettings())

	properties := createProperties(list)

	return Properties{properties}
}

func createProperties(propsMap map[string]string) []*Property{
	count := len(propsMap)

	ps := make([]*Property, count)
	i := 0
	for k, v := range propsMap {
		p, err := extractPlaceholder(v)
		if (err == nil){
			p := &Property{k, v, p}
			ps[i] = p
		}

		i++
	}

	return ps
}

func extractPlaceholder(s string) (string, error){
	if s[:2] != "${" {
		return "", fmt.Errorf("%s does not contain any placeholder with format ${PLACEHOLER_VARIABLE[:defaultvalue]}", s)
	}

	if s[len(s)-1:len(s)] != "}" {
		return "", fmt.Errorf("%s does not contain any placeholder with format ${PLACEHOLER_VARIABLE[:defaultvalue]}", s)
	}

	s = s[2:]
	s = s[0:len(s)-1]

	return strings.Split(s, ":")[0], nil
}

func parseMap(props map[string]interface{}) map[string]string {
	list := make(map[string]string)
	for key, value := range props {
		switch v := value.(type) {
		case map[interface{}]interface{}:
			l := parseMapInterface(v, key, list)
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

func parseMapInterface(props map[interface{}]interface{}, key string, list map[string]string) map[string]string {
	for k, value := range props {
		actKey := key + "." + fmt.Sprint(k)

		switch v := value.(type) {
		case map[interface{}]interface{}:
			list = parseMapInterface(v, actKey, list)
		case string:
			if v[:2] == "${" {
				keystr := fmt.Sprint(actKey) // <-- HACK
				list[keystr] = v
			}
		default:
		}
	}
	return list
}


func setOption(appname, label, variable, value string) {

	if label=="" {
		label = "default"
	}

	if options[appname]==nil {
		options[appname] = make(map[string]map[string]string)
	}
	if options[appname][label]==nil {
		options[appname][label] = make(map[string]string)
	}

	options[appname][label][variable] = value

}

func getOption(appname, label, variable string) string {

	if label=="" {
		label = "default"
	}

	if options[appname]==nil {
		options[appname] = make(map[string]map[string]string)
	}
	if options[appname][label]==nil {
		options[appname][label] = make(map[string]string)
	}

	return options[appname][label][variable]

}

func getOptions(appname, label string) map[string]string {

	if label=="" {
		label = "default"
	}

	if options[appname]==nil {
		options[appname] = make(map[string]map[string]string)
	}
	if options[appname][label]==nil {
		options[appname][label] = make(map[string]string)
	}

	return options[appname][label]

}
