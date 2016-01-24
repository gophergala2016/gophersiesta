package cmd

import (
	"github.com/spf13/cobra"
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"bytes"
	"crypto/tls"
	"log"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set the values to configuration manager. Needed {appname + label}",
	Long: "Set the values to be setup. From appname + label",
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		var res *http.Response

		if Source=="" {
			Source = "https://gophersiesta.herokuapp.com/"
		}
		if Source[len(Source)-1:]!="/" {
			Source += "/"
		}
		Url := Source + "conf/" + Appname + "/values"
		fmt.Println( Url )

		// Data is already a json (or should be)
		/*jsonString, jsonError := json.Marshal(params)
		if jsonError!=nil {
			log.Fatal(jsonError)
		} */
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}

		fmt.Println(Properties)
		req, _ := http.NewRequest("POST", Url, bytes.NewBuffer([]byte(Properties)))
		res, err = client.Do(req)
		if err!=nil {
			log.Fatal(err)
		}
		body, err := ioutil.ReadAll(res.Body)
		if err!=nil {
			log.Fatal(err)
		}
		res.Body.Close()
		var data string//map[string]interface{}{}
		json.Unmarshal(body, &data)
		
		fmt.Println(data)
	},
}

var Properties string

func init() {
	RootCmd.AddCommand(setCmd)

	setCmd.Flags().StringVarP(&Appname, "appname", "a", "", "Application name")
	setCmd.Flags().StringVarP(&Source, "source", "s", "https://gophersiesta.herokuapp.com/", "Source directory to read from")
	setCmd.Flags().StringVarP(&Label, "label", "l", "", "Select label to be fetch")
	setCmd.Flags().StringVarP(&Properties, "properties", "p", "", "json encoded properties")

}
