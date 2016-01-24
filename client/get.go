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


func GetValues() []byte {
	source := Source
	if source == "" {
		source = "https://gophersiesta.herokuapp.com/"
	}
	if source[len(source)-1:] != "/" {
		source += "/"
	}
	url := source + "conf/" + Appname + "/values"

	if (Label != ""){
		url = url + "?labels=" + Label
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
	RootCmd.AddCommand(getCmd)

	//aError := make([]error, 2)
	/*if err := renderCmd.MarkFlagRequired(Appname) ; err != nil {
		aError[0] = err
		log.Printf(" -> %q\n", err)
	}*/

	getCmd.Flags().StringVarP(&Appname, "appname", "a", "", "Application name")
	getCmd.Flags().StringVarP(&Source, "source", "s", "https://gophersiesta.herokuapp.com/", "Source directory to read from")
	getCmd.Flags().StringVarP(&Label, "label", "l", "", "Select label to be fetch")

	/*if len(aError)>0 {
		os.Exit(-1)
	}*/

}
