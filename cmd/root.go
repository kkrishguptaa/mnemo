/*
Copyright © 2025 Krish Gupta <m.krishggupta@icloud.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/kkrishguptaa/iris/util"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mnemo",
	Short: "Save small snippets of data and access them when needed.",
	Run:   mnemo,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func mnemo(cmd *cobra.Command, args []string) {
	// If there are 0 arguments, print the list of saved data
	// if there are 1 or more arguments, save the data
	if len(args) == 0 {
		data := util.FetchData()

		if len(data) == 0 {
			util.ErrorPrinter(fmt.Errorf("no data found. please save some data first using `mnemo <data>` command"))
			return
		}

		for i, d := range data {
			fmt.Printf("%d: %s\n", i+1, d)
		}
		return
	}

	// If there are arguments, save the data
	data := util.FetchData()

	data = append(data, args...)

	newData := util.ErrorHandler(json.Marshal(data))

	util.SaveData(newData)
}
