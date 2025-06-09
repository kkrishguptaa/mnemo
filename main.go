package main

import (
	"fmt"
	"os"

	"github.com/kkrishguptaa/iris/commands/list"
	"github.com/kkrishguptaa/iris/commands/read"
	"github.com/kkrishguptaa/iris/commands/write"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "iris",
		Short: "Iris, write a letter to your future self.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	writeCmd = &cobra.Command{
		Use:   "write",
		Short: "Write a letter to your future self.",
		Args:  cobra.NoArgs,
		RunE:  write.Write,
	}

	listCmd = &cobra.Command{
		Use:   "list",
		Short: "List all letters you have written.",
		Args:  cobra.NoArgs,
		RunE:  list.List,
	}

	readCmd = &cobra.Command{
		Use:   "read <date>",
		Short: "Read a letter you have written for a specific date.",
		Args:  cobra.ExactArgs(1),
		RunE:  read.Read,
	}
)

func main() {
	rootCmd.AddCommand(writeCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(readCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}
