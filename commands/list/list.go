package list

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

func List(cmd *cobra.Command, args []string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("could not get home directory: %v", err)
	}

	directory := home + "/.iris/letters"
	files, err := os.ReadDir(directory)
	if err != nil {
		return fmt.Errorf("could not read letters directory: %v", err)
	}

	if len(files) == 0 {
		fmt.Println("No letters found.")
		return nil
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	letters := []string{}

	for _, file := range files {
		date := strings.Split(file.Name(), ".")[0]

		text, err := os.ReadFile(directory + "/" + file.Name())

		if err != nil {
			return fmt.Errorf("could not read letter file: %v", err)
		}

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
		if strings.Contains(letter, time.Now().Format("2006-01-02")) {
			fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("3")).Bold(true).Render(letter))
		} else if strings.Compare(strings.Split(letter, ":")[0], time.Now().Format("2006-01-02")) < 0 {
			fmt.Println(lipgloss.NewStyle().Strikethrough(true).Render(letter))
		} else {
			// gray
			fmt.Println(lipgloss.NewStyle().Render(letter))
		}
	}

	footer := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render("\nUse 'iris read <date>' to read a letter. date format: YYYY-MM-DD")

	fmt.Println(footer)

	return nil
}
