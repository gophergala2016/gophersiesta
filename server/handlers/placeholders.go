package handlers

import (
	"github.com/gophergala2016/gophersiesta/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"net/http"
	"github.com/gophergala2016/gophersiesta/server/placeholders"
	"encoding/json"
	"strings"
	"github.com/gophergala2016/gophersiesta/server/storage"
	"fmt"
	"io/ioutil"
)

func GetPlaceHolders(c *gin.Context) {
	name := c.Param("appname")
	myViper, err := readTemplate(name)
	if err != nil {
	c.String(http.StatusNotFound, "Config file for %s not found\n", name)
	} else {
	properties := placeholders.GetPlaceHolders(myViper)
	propsJson, _ := json.Marshal(properties)
	c.String(http.StatusOK, string(propsJson))
	}
}

func GetValues(s storage.Storage) func (c *gin.Context){
	return func (c *gin.Context) {
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

		vs_json, _ := json.Marshal(vs)
		c.String(http.StatusOK, string(vs_json))
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