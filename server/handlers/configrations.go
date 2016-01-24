package handlers

import (
	"github.com/gophergala2016/gophersiesta/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"github.com/gophergala2016/gophersiesta/Godeps/_workspace/src/github.com/spf13/viper"
	"github.com/gophergala2016/gophersiesta/server/placeholders"
	"github.com/gophergala2016/gophersiesta/server/storage"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Labels is a collection of labels that can be associated with the values of the placeholders
type Labels struct {
	Labels []string `json:"labels"`
}

// GetConfig return the configuration file for a given appname
func GetConfig(c *gin.Context) {
	name := c.Param("appname")
	myViper, err := readTemplate(name)

	if err != nil {
		c.String(http.StatusNotFound, "Config file for %s not found\n", name)
	} else {
		filename := myViper.ConfigFileUsed()
		content, err := fileReader(filename)

		if err != nil {
			c.String(http.StatusNotFound, "Config file for %s not found\n", name)
		}
		w := c.Writer

		io.Copy(w, content)

	}
}


// GetLabels return the stored labels for a given appname
func GetLabels(s storage.Storage) func(c *gin.Context) {
	return func(c *gin.Context) {

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

func getFileExtension(v *viper.Viper) string {
	filename := v.ConfigFileUsed()
	extension := filepath.Ext(filename)

	extension = strings.Replace(extension, ".", "", 1)

	return extension
}

func readConfigFile(v *viper.Viper) string {

	filename := v.ConfigFileUsed()

	r, err := fileReader(filename)

	if err != nil {
		return ""
	}

	configFile, err := ioutil.ReadAll(r)

	if err != nil {
		return ""
	}

	return string(configFile)
}

func replaceTemplatePlaceHolders(v *viper.Viper, list map[string]*placeholders.Placeholder) string {

	template := readConfigFile(v)

	for _, v := range list {
		re := regexp.MustCompile("\\${" + v.PlaceHolder + ":?([^}]*)}")
		template = re.ReplaceAllString(template, v.PropertyValue)
	}

	return template
}
