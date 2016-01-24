package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gophergala2016/gophersiesta/Godeps/_workspace/src/github.com/BurntSushi/toml"
	"github.com/gophergala2016/gophersiesta/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"github.com/gophergala2016/gophersiesta/Godeps/_workspace/src/github.com/spf13/viper"
	"github.com/gophergala2016/gophersiesta/Godeps/_workspace/src/gopkg.in/yaml.v2"
	"github.com/gophergala2016/gophersiesta/server/placeholders"
	"github.com/gophergala2016/gophersiesta/server/storage"
	"io/ioutil"
	"net/http"
	"strings"
)

// GetPlaceHolders return the placeholders that are present in the appname config file
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

// GetValues return the placeholders values for a given appname and concrete namespace represented as labels.
// If no labels are provided the "default" label is used
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

		vs := placeholders.CreateValues(list)

		c.IndentedJSON(http.StatusOK, vs)
	}
}

// SetValues set the placeholders values for a given appname and concrete namespace represented as labels
// If no labels are provided the "default" label is used
func SetValues(s storage.Storage) func(c *gin.Context) {
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
			if len(data) > 0 { //it's a JSON
				for _, label := range lbls {
					for k, v := range data {
						s.SetOption(name, label, k, fmt.Sprint(v))
					}
				}
				c.String(http.StatusOK, "Ok")
			} else if strings.Contains(string(x), "=") {
				pairs := strings.Split(string(x), ",")
				for _, label := range lbls {
					for _, v := range pairs {
						vv := strings.Split(v, "=")
						s.SetOption(name, label, vv[0], strings.Join(vv[1:], "="))
					}
				}
				c.String(http.StatusOK, "Ok")
			} else {
				c.String(http.StatusBadRequest, "Properties not well-formed")
			}

		}

	}
}

// ReplacePlaceholders generated a config file using the base config file for appname by replacing the possible placeholders
// that the config fiel might have given a group of labels. These labels are used to retrieve the placeholders stored values
func ReplacePlaceholders(s storage.Storage) func(c *gin.Context) {
	return func(c *gin.Context) {
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
				list[v.PlaceHolder] = v
			}

			lbls := strings.Split(labels, ",")
			// MERGE values of different labels, last overrides current value
			for _, label := range lbls {
				l := s.GetOptions(name, label)
				for k, v := range l {
					if list[k] != nil {
						list[k].PropertyValue = v
					}
				}
			}

			template := replaceTemplatePlaceHolders(myViper, list)

			extension := getFileExtension(myViper)

			replacedViper := viper.New()
			replacedViper.SetConfigType(extension)
			replacedViper.ReadConfig(bytes.NewBuffer([]byte(template)))

			b, err := render(replacedViper, renderType)

			if err == nil {
				c.Data(http.StatusOK, "text/plain", b)
			}

			if err != nil {

				c.String(http.StatusInternalServerError, "Could not render %s", err)
			}
		}

	}
}

func render(v *viper.Viper, configOutputType string) ([]byte, error) {

	var b []byte
	var err error
	var conf = make(map[string]interface{})

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
