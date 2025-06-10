package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/kkrishguptaa/iris/util"
	"github.com/spf13/cobra"
)

// listCommand represents the write command
var listCommand = &cobra.Command{
	Use:   "list",
	Short: "List all the letters you have written",
	Long:  `List all the letters you have written. This command will display the titles and dates of all your letters.`,
	RunE:  list,
	Args:  cobra.NoArgs,
}

func init() {
	rootCmd.AddCommand(listCommand)
}

func list(cmd *cobra.Command, args []string) error {
	var files = util.ErrorHandler(os.ReadDir(util.LettersPath))

	if len(files) == 0 {
		util.ErrorPrinter(fmt.Errorf("you have not written any letters yet"))
		return nil
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	letters := []string{}

	for _, file := range files {
		date := strings.Split(file.Name(), ".")[0]

		text := util.ErrorHandler(os.ReadFile(util.LettersPath + "/" + file.Name()))

		title := strings.Split(string(text), "\n")[0]
		title = strings.TrimPrefix(title, "# ")
		letters = append(letters, date+": "+title)
	}

	heading := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true).
		Padding(1, 0, 1, 0).
		Render("Your Letters")

	fmt.Println(heading)

	for _, letter := range letters {

		var date = util.ErrorHandler(time.Parse("2006-01-02", strings.Split(letter, ":")[0]))

		var todayDate = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC)

		if date.Equal(todayDate) {
			fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("3")).Bold(true).Render(letter))
		} else if date.Before(todayDate) {
			fmt.Println(lipgloss.NewStyle().Strikethrough(true).Render(letter))
		} else {
			fmt.Println(lipgloss.NewStyle().Render(letter))
		}
	}

	footer := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render("\nUse 'iris read' to read a letter.")

	fmt.Println(footer)

	return nil
}
