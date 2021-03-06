package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/gophergala2016/gophersiesta/Godeps/_workspace/src/github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"net/http"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set all the values to configuration manager. Needed {appname + label}",
	Long:  "Set all the values to be setup. From appname + label",
	Run: func(cmd *cobra.Command, args []string) {

		SendProp(properties, label)
	},
}

// SendProp send the prop to the config service for the namespace label
func SendProp(prop string, label string) {
	var err error
	var res *http.Response

	if source == "" {
		source = "https://gophersiesta.herokuapp.com/"
	}
	if source[len(source)-1:] != "/" {
		source += "/"
	}
	url := source + "conf/" + appName + "/values"

	if label != "" {
		url = url + "?labels=" + label
	}

	fmt.Println(url)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	fmt.Println(prop)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer([]byte(prop)))
	res, err = client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	res.Body.Close()
	var data string //map[string]interface{}{}
	json.Unmarshal(body, &data)

	fmt.Println(data)
}


func init() {
	rootCmd.AddCommand(setCmd)

	setCmd.Flags().StringVarP(&appName, "appname", "a", "", "Application name")
	setCmd.Flags().StringVarP(&source, "source", "s", "https://gophersiesta.herokuapp.com/", "Source directory to read from")
	setCmd.Flags().StringVarP(&label, "label", "l", "", "Select label to be fetch")
	setCmd.Flags().StringVarP(&properties, "properties", "p", "", "json encoded properties")

}
