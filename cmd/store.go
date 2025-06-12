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

	"slices"

	"github.com/kkrishguptaa/mnemo/lib"
	"github.com/kkrishguptaa/mnemo/util"
	"github.com/spf13/cobra"
)

var storeCmd = &cobra.Command{
	Use:     "store",
	Short:   "Manage your snippet stores",
	Long:    `The store command allows you to manage your snippet stores. You can create, list, and delete stores, as well as perform other operations related to snippet storage.`,
	GroupID: "mnemo",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func listStores(cmd *cobra.Command, args []string) {
	stores := lib.ListStores()

	if len(stores) == 0 {
		fmt.Printf("No stores found. Use 'mnemo store create [name]' to create a new store.\n")
		return
	}

	for _, store := range stores {
		fmt.Println(store)
	}
}

func createStore(cmd *cobra.Command, args []string) {
	name := args[0]
	stores := lib.ListStores()

	if slices.Contains(stores, name) {
		util.ErrorPrinter(fmt.Errorf("store already exists, if you wish to clear it, use 'mnemo store clear %s' command", name))
		return
	}

	lib.CreateStore(name)
	util.SuccessPrinter(fmt.Sprintf("Store '%s' created successfully.", name))
}

func clearStore(cmd *cobra.Command, args []string) {
	name := args[0]
	if len(args) != 1 {
		util.ErrorPrinter(fmt.Errorf("please provide a store name to clear"))
		return
	}
	store := lib.FetchStore(name)

	lib.WriteStore(store.Name, []lib.Snippet{})

	util.SuccessPrinter(fmt.Sprintf("Store '%s' cleared successfully.", name))
}

func deleteStore(cmd *cobra.Command, args []string) {
	name := args[0]
	if len(args) != 1 {
		util.ErrorPrinter(fmt.Errorf("please provide a store name to delete"))
		return
	}

	lib.DeleteStore(name)
	util.SuccessPrinter(fmt.Sprintf("Store '%s' deleted successfully.", name))
}

func init() {
	var listStoresCmd = &cobra.Command{
		Use:   "list",
		Short: "List all snippet stores",
		Long:  `List all snippet stores available in the system. This command will display the names of all stores that have been created.`,
		Args:  cobra.NoArgs,
		Run:   listStores,
	}
	storeCmd.AddCommand(listStoresCmd)

	var createStoreCmd = &cobra.Command{
		Use:   "create [name]",
		Short: "Create a new snippet store",
		Long:  `Create a new snippet store with the specified name. If a store with the same name already exists, an error will be returned.`,
		Args:  cobra.ExactArgs(1),
		Run:   createStore,
	}
	storeCmd.AddCommand(createStoreCmd)

	var clearStoreCmd = &cobra.Command{
		Use:   "clear [name]",
		Short: "Clear a snippet store",
		Long:  `Clear all snippets from the specified store. This will remove all snippets from the store, but the store itself will remain.`,
		Args:  cobra.ExactArgs(1),
		Run:   clearStore,
	}
	storeCmd.AddCommand(clearStoreCmd)

	var deleteStoreCmd = &cobra.Command{
		Use:   "delete [name]",
		Short: "Delete a snippet store",
		Long:  `Delete the specified snippet store. This will remove the store and all its contents permanently.`,
		Args:  cobra.ExactArgs(1),
		Run:   deleteStore,
	}
	storeCmd.AddCommand(deleteStoreCmd)

	rootCmd.AddCommand(&cobra.Command{
		Use:   "ls",
		Short: "Alias for 'mnemo store list'",
		Args:  cobra.NoArgs,
		Run:   listStores,
	})

	rootCmd.AddCommand(storeCmd)
}
