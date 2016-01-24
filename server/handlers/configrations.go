package handlers

import (
	"github.com/gophergala2016/gophersiesta/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"net/http"
	"github.com/gophergala2016/gophersiesta/Godeps/_workspace/src/github.com/spf13/viper"
	"github.com/gophergala2016/gophersiesta/server/storage"
	"os"
	"io"
)


type Labels struct {
	Labels []string `json:"labels"`
}

func GetConfig(c *gin.Context) {
	name := c.Param("appname")
	myViper, err := readTemplate(name)

	if err != nil {
		c.String(http.StatusNotFound, "Config file for %s not found\n", name)
	} else {
		filename := myViper.ConfigFileUsed()
		content, err := fileReader(filename)

		if (err != nil){
			c.String(http.StatusNotFound, "Config file for %s not found\n", name)
		}
		w := c.Writer

		io.Copy(w, content)

	}
}

func GetLabels(s storage.Storage) func (c *gin.Context){
	return func (c *gin.Context) {

		name := c.Param("appname")

		lbls := s.GetLabels(name)

		labels := &Labels{lbls}

		c.IndentedJSON(http.StatusOK, labels)
	}
}

func readTemplate(appname string) (*viper.Viper, error) {

	aux := viper.New()
	aux.SetConfigName("config")
	aux.AddConfigPath("apps/" + appname + "/")

	err := aux.ReadInConfig()
	return aux, err

}

func fileReader(filename string) (io.Reader, error) {
	file, errFile := os.Open(filename)

	if errFile != nil {
		return nil, errFile
	}

	return file, errFile
}