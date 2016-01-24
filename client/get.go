package client

import (
	"fmt"

	"github.com/gophergala2016/gophersiesta/Godeps/_workspace/src/github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"net/http"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get --appname=[appname] --source=[source_url]",
	Short: "Get the values to be setup. From appname + label",
	Long:  "Get the values to be setup. From appname + label",
	Run: func(cmd *cobra.Command, args []string) {
		body := GetValues()
		fmt.Println(string(body))
	},
}
// GetValues return the raw response from the server calling http://url/conf/:appname/values
func GetValues() []byte {

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

	fmt.Println("[api call] " + url)
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	return body
}

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.Flags().StringVarP(&appName, "appname", "a", "", "Application name")
	getCmd.Flags().StringVarP(&source, "source", "s", "https://gophersiesta.herokuapp.com/", "Source directory to read from")
	getCmd.Flags().StringVarP(&label, "label", "l", "", "Select label to be fetch")

}
