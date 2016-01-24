package cmd

import (
	"github.com/spf13/cobra"
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"bytes"
	"crypto/tls"
	"strings"
	"log"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set the values to configuration manager. Needed {appname + label}",
	Long: "Get the values to be setup. From appname + label",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Print: " + strings.Join(args, " "))

		var err error
		var res *http.Response
		// create post url
		// https://gophersiesta.herokuapp.com/conf/app1/
		// Url := Source + "/conf/" + Appname + "/get/?labels=" + Label
		Url := "http://garciademarina.com/"
		fmt.Println( Url )

		// post data to service
		jsonString, jsonError := json.Marshal(params)
		if jsonError!=nil {
			log.Fatal(jsonError)
		}
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}
		req, _ := http.NewRequest("POST", Url, bytes.NewBuffer([]byte(jsonString)))
		res, err = client.Do(req)
		if err!=nil {
			log.Fatal(err)
		}
		body, err := ioutil.ReadAll(res.Body)
		if err!=nil {
			log.Fatal(err)
		}
		res.Body.Close()
		data := map[string]interface{}{}
		json.Unmarshal(body, &data)
		
		fmt.Println(data)
	},
}

var params map[string]string

func init() {
	RootCmd.AddCommand(setCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
