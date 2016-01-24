package client

import (
	"fmt"

	"github.com/spf13/cobra"
	"log"
	"net/http"
	"io/ioutil"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get --appname=[appname] --source=[source_url]",
	Short: "Get the values to be setup. From appname + label",
	Long: "Get the values to be setup. From appname + label",
	Run: func(cmd *cobra.Command, args []string) {
		if Source=="" {
			Source = "https://gophersiesta.herokuapp.com/"
		}
		if Source[len(Source)-1:]!="/" {
			Source += "/"
		}
		Url := Source + "conf/" + Appname + "/values"
		fmt.Println( "[api call] " + Url )
		res, err := http.Get( Url )
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)
		fmt.Println(string(body))
	},
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