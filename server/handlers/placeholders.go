package handlers

import (
	"github.com/gophergala2016/gophersiesta/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"net/http"
	"github.com/gophergala2016/gophersiesta/server/placeholders"
	"strings"
	"github.com/gophergala2016/gophersiesta/server/storage"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"regexp"
	"bytes"
	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v2"
	"github.com/gophergala2016/gophersiesta/Godeps/_workspace/src/github.com/spf13/viper"
	"path/filepath"
)

func GetPlaceHolders(c *gin.Context) {
	name := c.Param("appname")
	myViper, err := readTemplate(name)
	if err != nil {
		c.String(http.StatusNotFound, "Config file for %s not found\n", name)
	} else {
		properties := placeholders.GetPlaceHolders(myViper)
		c.IndentedJSON(http.StatusOK, properties)
	}
}

func GetValues(s storage.Storage) func(c *gin.Context) {
	return func(c *gin.Context) {
		name := c.Param("appname")
		labels := c.DefaultQuery("labels", "default")

		list := make(map[string]string)
		if strings.Contains(labels, ",") {
			lbls := strings.Split(labels, ",")
			// MERGE values of different labels, last overrides current value
			for _, label := range lbls {
				l := s.GetOptions(name, label)
				for k, v := range l {
					list[k] = v
				}
			}
		} else {
			list = s.GetOptions(name, labels)
		}

		vs := placeholders.CreateValues(list);

		c.IndentedJSON(http.StatusOK, vs)
	}
}


func SetValues(s storage.Storage) func (c *gin.Context) {
	return func(c *gin.Context) {
		name := c.Param("appname")
		labels := c.DefaultQuery("labels", "default")

		body := c.Request.Body
		x, err := ioutil.ReadAll(body)

		if err != nil {
			c.String(http.StatusBadRequest, "Bad request")
		} else {

			data := map[string]interface{}{}
			json.Unmarshal(x, &data)

			lbls := strings.Split(labels, ",")
			for _, label := range lbls {
				for k, v := range data {
					s.SetOption(name, label, k, fmt.Sprint(v))
				}
			}


			c.String(http.StatusOK, "Ok")
		}

	}
}

func ReplacePlaceholders(s storage.Storage) func (c *gin.Context){
	return func (c *gin.Context) {
		name := c.Param("appname")
		labels := c.DefaultQuery("labels", "default")
		renderType := c.Param("format")

		list := make(map[string]*placeholders.Placeholder)

		myViper, err := readTemplate(name)
		if err != nil {
			c.String(http.StatusNotFound, "")
		} else {
			properties := placeholders.GetPlaceHolders(myViper)
			for _, v := range properties.Placeholders {
				list[v.PropertyName] = v
			}

			lbls := strings.Split(labels, ",")
			// MERGE values of different labels, last overrides current value
			for _, label := range lbls {
				l := s.GetOptions(name, label)
				for k, v := range l {
					if list[k]!=nil {
						list[k].PropertyValue = v
					}
				}
			}

			filename := myViper.ConfigFileUsed()
			template := safeFileRead(filename)

			for _, v := range list {
				re := regexp.MustCompile("\\${" + v.PlaceHolder + ":?([^}]*)}")
				template = re.ReplaceAllString(template, v.PropertyValue)
			}

			replacedViper :=  viper.New()
			extension := filepath.Ext(filename)

			extension = strings.Replace(extension, ".", "", 1)

			replacedViper.SetConfigType(extension)
			replacedViper.ReadConfig(bytes.NewBuffer([]byte(template)))
			b , err := render(replacedViper, renderType)

			if err == nil {
				c.Data(http.StatusOK, "text/plain", b)
			}else{
				c.String(http.StatusInternalServerError, "Could not render %s", err)
			}
		}

	}
}


func render(v *viper.Viper, configOutputType string) ([]byte, error) {

	var b []byte
	var err error
	var conf  = make(map[string]interface{})

	m := v.AllSettings()

	b, err = yaml.Marshal(m)

	err = yaml.Unmarshal(b, conf)

	if err == nil {
		switch configOutputType {
		case "json":

			b, err = json.MarshalIndent(conf, "", "    ")

		case "toml":
			var buff bytes.Buffer

			err = toml.NewEncoder(&buff).Encode(conf)
			b = buff.Bytes()

		case "yaml", "yml":

			b, err = yaml.Marshal(conf)

		}
	}

	if err != nil {
		return nil, err
	}

	return b, nil

}
