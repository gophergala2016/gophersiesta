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
	Long: `Fetch configuration files for a given <label>.
Fetched from source url.`,
	Run: func(cmd *cobra.Command, args []string) {

		if source == "" {
			source = "https://gophersiesta.herokuapp.com/"
		}
		if source[len(source)-1:] != "/" {
			source += "/"
		}

		url := source + "conf/" + appName + "/render/yml"

		if label != "" {
			url = url + "?labels=" + label
		}

		fmt.Println(url)

		res, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}

		defer res.Body.Close()

		body, _ := ioutil.ReadAll(res.Body)

		fmt.Println(string(body))
	},
}

func init() {
	rootCmd.AddCommand(renderCmd)

	renderCmd.Flags().StringVarP(&appName, "appname", "a", "", "Application name")
	renderCmd.Flags().StringVarP(&source, "source", "s", "", "Source directory to read from")
	renderCmd.Flags().StringVarP(&label, "label", "l", "", "Select label to be fetch")
}
