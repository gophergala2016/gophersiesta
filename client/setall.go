package client

import (
	"encoding/json"
	"fmt"
	"github.com/gophergala2016/gophersiesta/Godeps/_workspace/src/github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"net/http"
	"github.com/gophergala2016/gophersiesta/server/placeholders"
	"bufio"
	"os"
	"bytes"
)

// setCmd represents the set command
var setAllCmd = &cobra.Command{
	Use:   "setall",
	Short: "Set the values to configuration manager. Needed {appname + label}",
	Long:  "Set the values to be setup. From appname + label",
	Run: func(cmd *cobra.Command, args []string) {
		getPlaceholders()
	},
}

func getPlaceholders(){

	 pls := &placeholders.Placeholders{}

	if Source == "" {
		Source = "https://gophersiesta.herokuapp.com/"
	}
	if Source[len(Source)-1:] != "/" {
		Source += "/"
	}
	url := Source + "conf/" + Appname + "/placeholders"

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

	err = json.Unmarshal(body, pls)

	if (err != nil){
		log.Fatal("Could not get Placehodlers")
	}

	fmt.Printf("There must be %v placeholders set. Listing: \n", len(pls.Placeholders))

	for _, p := range pls.Placeholders {
		fmt.Printf("%s[%s]\n", p.PropertyName, p.PlaceHolder)
	}

	fmt.Printf("Type value for each placeholder: \n")

	var buff bytes.Buffer
	for i, p := range pls.Placeholders {
		setPropertyderValue(p)
		buff.WriteString(p.PlaceHolder)
		buff.WriteString("=")
		buff.WriteString(p.PropertyValue)

		if i != len(pls.Placeholders) -1 {
			buff.WriteString(",")
		}
	}

	SendProp(string(buff.Bytes()))
}

func setPropertyderValue(p *placeholders.Placeholder){

	fmt.Printf("%s[%s]:", p.PropertyName, p.PlaceHolder)

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		p.PropertyValue = scanner.Text()
		break;
	}
}

func init() {
	RootCmd.AddCommand(setAllCmd)

	setAllCmd.Flags().StringVarP(&Appname, "appname", "a", "", "Application name")
	setAllCmd.Flags().StringVarP(&Source, "source", "s", "https://gophersiesta.herokuapp.com/", "Source directory to read from")
	setAllCmd.Flags().StringVarP(&Label, "label", "l", "", "Select label to be fetch")

}
