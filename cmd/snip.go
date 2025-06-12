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
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/kkrishguptaa/mnemo/lib"
	"github.com/kkrishguptaa/mnemo/util"
	"github.com/spf13/cobra"
)

// snipCmd represents the snip command
var snipCmd = &cobra.Command{
	Use:   "snip",
	Short: "Manage your snippets",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func listSnippets(cmd *cobra.Command, args []string) {
	// This function will list all snippets in the specified store.
	store, _ := cmd.Flags().GetString("store")

	if store == "" {
		store = "default"
	}

	Store := lib.FetchStore(store)

	for _, snippet := range Store.Data {
		if snippet.Encrypted {
			println(
				lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205")).Render(snippet.Id),
				":",
				lipgloss.NewStyle().Foreground(lipgloss.Color("8")).
					Render("encrypted, use 'mnemo snip read "+snippet.Id+"' to view the value"))
			continue
		}

		println(lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205")).Render(snippet.Id), ":", snippet.Value)
	}
}

func CreateSnip(cmd *cobra.Command, args []string) {
	store, _ := cmd.Flags().GetString("store")
	if store == "" {
		store = "default"
	}

	encrypted := false

	password, _ := cmd.Flags().GetString("password")

	if password != "" {
		encrypted = true
	}

	key := args[0]
	value := strings.Join(args[1:], " ")

	if key == "" || value == "" {
		cmd.Help()
		return
	}

	Store := lib.FetchStore(store)

	if encrypted {
		value = lib.Encrypt(value, password)
	}

	snippet := lib.Snippet{
		Id:        key,
		Value:     value,
		Encrypted: encrypted,
	}

	Store.Data = append(Store.Data, snippet)

	lib.WriteStore(store, Store.Data)

	util.SuccessPrinter("Snippet created successfully in store: " + store)
}

func ReadSnip(cmd *cobra.Command, args []string) {
	store, _ := cmd.Flags().GetString("store")
	if store == "" {
		store = "default"
	}

	if len(args) != 1 {
		cmd.Help()
		return
	}

	key := args[0]

	Store := lib.FetchStore(store)

	for _, snippet := range Store.Data {
		if snippet.Id == key {
			if snippet.Encrypted {
				password, _ := cmd.Flags().GetString("password")
				if password == "" {
					util.ErrorPrinter(fmt.Errorf("snippet with key '%s' is encrypted, please provide a password using the '-p' flag", key))
					return
				}

				value := lib.Decrypt(snippet.Value, password)

				println(value)
			} else {
				println(snippet.Value)
			}
			return
		}
	}

	util.ErrorPrinter(fmt.Errorf("snippet with key '%s' not found in store '%s'", key, store))
}

func init() {
	var listSnipsCmd = &cobra.Command{
		Use:     "list",
		Short:   "List all snippets in the specified store",
		Aliases: []string{"ls"},
		Run:     listSnippets,
	}
	ApplySnipFlags(listSnipsCmd)
	snipCmd.AddCommand(listSnipsCmd)

	var createSnipCmd = &cobra.Command{
		Use:   "create [key] [value]",
		Short: "Create a new snippet in the specified store",
		Long:  `Create a new snippet with the specified key and value in the specified store. If the store does not exist, it will be created.`,
		Args:  cobra.ExactArgs(2),
		Run:   CreateSnip,
	}
	ApplySnipFlags(createSnipCmd)
	createSnipCmd.Flags().StringP("password", "p", "", "Password for the store (if you want it encrypted)")
	snipCmd.AddCommand(createSnipCmd)

	var readSnipCmd = &cobra.Command{
		Use:   "read [key]",
		Short: "Read a snippet from the specified store",
		Long:  `Read a snippet with the specified key from the specified store. If the snippet is encrypted, you will need to provide the password using the '-p' flag.`,
		Args:  cobra.ExactArgs(1),
		Run:   ReadSnip,
	}
	ApplySnipFlags(readSnipCmd)
	readSnipCmd.Flags().StringP("password", "p", "", "Password for the snippet (if it is encrypted)")
	snipCmd.AddCommand(readSnipCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// snipCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// snipCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.AddCommand(snipCmd)

}

func ApplySnipFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("store", "s", "", "Specify the store to use for snippets")
}
