package client

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gophergala2016/gophersiesta/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/gophergala2016/gophersiesta/server/placeholders"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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

func getPlaceholders() {

	pls, err := readPlaceHolders()

	if err != nil {
		log.Fatal("Could not get Placehodlers")
	}

	mValues, err := readValues()

	if err != nil {
		log.Fatal("Could not get Placehodlers current stored values")
	}

	fmt.Printf("\nThere are %v placeholders. Listing: \n", len(pls.Placeholders))

	for _, p := range pls.Placeholders {
		fmt.Printf("%s [$%s]\n", p.PropertyName, p.PlaceHolder)
	}

	fmt.Printf("\n\n")
	fmt.Printf("Type value for each placeholder and press ENTER, or ENTER to skip or left as before: \n")
	fmt.Printf("	explanation: property.path [$PLACE_HOLDER] --> curentvalue \n")

	for _, p := range pls.Placeholders {

		pv := mValues[p.PlaceHolder]

		setPropertyHolderValue(p, pv)
	}

}

func readPlaceHolders() (*placeholders.Placeholders, error) {
	pls := &placeholders.Placeholders{}

	source := Source
	if source == "" {
		source = "https://gophersiesta.herokuapp.com/"
	}
	if source[len(source)-1:] != "/" {
		source += "/"
	}
	url := source + "conf/" + Appname + "/placeholders"

	label := Label
	if label != "" {
		url = url + "?labels=" + label
	}

	fmt.Println("[api call] " + url)
	res, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	err = json.Unmarshal(body, pls)

	if err != nil {
		return nil, err
	}

	return pls, err
}

func readValues() (map[string]string, error) {

	vs := &placeholders.Values{}
	mValues := make(map[string]string)

	body := GetValues()

	err := json.Unmarshal(body, vs)

	if err != nil {
		return nil, err
	}

	for _, v := range vs.Values {
		mValues[v.Name] = v.Value
	}

	return mValues, err
}

func setPropertyHolderValue(p *placeholders.Placeholder, currentVal string) {
	var value string
	fmt.Printf("%s [$%s] --> %s:", p.PropertyName, p.PlaceHolder, currentVal)

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		value = scanner.Text()
		break
	}

	// send Property
	if value != "" {

		var buff bytes.Buffer

		buff.WriteString(p.PlaceHolder)
		buff.WriteString("=")
		buff.WriteString(value)

		SendProp(string(buff.Bytes()), Label)
	}
}

func init() {
	RootCmd.AddCommand(setAllCmd)

	setAllCmd.Flags().StringVarP(&Appname, "appname", "a", "", "Application name")
	setAllCmd.Flags().StringVarP(&Source, "source", "s", "https://gophersiesta.herokuapp.com/", "Source directory to read from")
	setAllCmd.Flags().StringVarP(&Label, "label", "l", "", "Select label to be fetch")

}
