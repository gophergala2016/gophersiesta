// Copyright © 2016 GOPHERSIESTA
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package client

import (
	"fmt"

	"github.com/gophergala2016/gophersiesta/Godeps/_workspace/src/github.com/spf13/cobra"
	"net/http"

	"io/ioutil"
	"log"
)

// renderCmd represents the render command
var renderCmd = &cobra.Command{
	Use:   "render",
	Short: "Get the configuration file from config source",
	Long: `

.      .__                 __
  _____|__| ____   _______/  |______
 /  ___/  |/ __ \ /  ___/\   __\__  \
 \___ \|  \  ___/ \___ \  |  |  / __ \_
/____  >__|\___  >____  > |__| (____  /
     \/        \/     \/            \/

Fetch configuration files for a given <label>.
Fetched from source url.`,
	Run: func(cmd *cobra.Command, args []string) {

		if Source == "" {
			Source = "https://gophersiesta.herokuapp.com/"
		}
		if Source[len(Source)-1:] != "/" {
			Source += "/"
		}

		fmt.Println("Source " + Source)
		/*
			err := checkArgs()
			if err != nil {
				log.Fatal(err)
			}
		*/
		// https://gophersiesta.herokuapp.com/conf/app1/
		Url := Source + "/conf/" + Appname
		fmt.Println(Url)

		res, err := http.Get(Url)
		if err != nil {
			log.Fatal(err)
		}

		defer res.Body.Close()

		body, _ := ioutil.ReadAll(res.Body)

		fmt.Println(string(body))
	},
}

var Appname string
var Source string
var Label string

func init() {
	RootCmd.AddCommand(renderCmd)

	/*	aError := make([]error, 2)
		if err := renderCmd.MarkFlagRequired("appname") ; err != nil {
			aError[0] = err
			log.Printf(" -> %q\n", err)
		}

		if err := renderCmd.MarkFlagRequired("source") ; err != nil {
			aError[1] = err
			log.Printf(" -> %q\n", err)
		}*/
	renderCmd.Flags().StringVarP(&Appname, "appname", "a", "", "Application name")
	renderCmd.Flags().StringVarP(&Source, "source", "s", "", "Source directory to read from")
	renderCmd.Flags().StringVarP(&Label, "label", "l", "", "Select label to be fetch")

	/*if len(aError)>0 {
		os.Exit(-1)
	}*/

}
