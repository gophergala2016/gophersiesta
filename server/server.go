package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/spf13/viper"
	"fmt"
)

func main() {
	router := gin.Default()

	// This handler will match /conf/appname but will not match neither /conf/ or /conf
	router.GET("/conf/:appname", func(c *gin.Context) {
		name := c.Param("appname")
		myViper, err := readConfig(name)
		if err!=nil {
			c.String(http.StatusNotFound, "Config file for %s not found\n", name)
		} else {
			fmt.Println(myViper)
			c.String(http.StatusOK, "Config file %s \n", name)
		}
	})

	// However, this one will match /conf/app1/ and also /conf/app1/send
	// If no other routers match /conf/app1, it will redirect to /conf/app1/
	router.GET("/conf/:appname/*action", func(c *gin.Context) {
		name := c.Param("appname")
		action := c.Param("action")
		message := name + " is " + action + "\n"
		c.String(http.StatusOK, message)
	})

	router.Run(":8080")
}

func readConfig(appname string) (*viper.Viper, error){

	aux := viper.New()
	aux.SetConfigName("config")
	aux.AddConfigPath("apps/" + appname + "/")

	err := aux.ReadInConfig()
	/*if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}*/

	return aux, err

}
