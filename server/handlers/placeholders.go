package handlers

import (
	"github.com/gophergala2016/gophersiesta/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"net/http"
	"github.com/gophergala2016/gophersiesta/server/placeholders"
	"strings"
	"github.com/gophergala2016/gophersiesta/server/storage"
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

