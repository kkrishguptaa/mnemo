package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
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
		RunE:  write,
	}

	listCmd = &cobra.Command{
		Use:   "list",
		Short: "List all letters you have written.",
		Args:  cobra.NoArgs,
		RunE:  list,
	}

	readCmd = &cobra.Command{
		Use:   "read <date>",
		Short: "Read a letter you have written for a specific date.",
		Args:  cobra.ExactArgs(1),
		RunE:  read,
	}
)

func write(cmd *cobra.Command, args []string) error {
	var (
		title   string
		date    string
		content string
	)

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("What would be the title of your letter?").CharLimit(50).Value(&title),
			huh.NewInput().Title("What date would you like to write this letter for? YYYY-MM-DD").
				Validate(func(input string) error {
					if _, err := time.Parse("2006-01-02", input); err != nil {
						return fmt.Errorf("invalid date format, please use YYYY-MM-DD")
					}
					return nil
				}).
				CharLimit(10).
				Value(&date),
		),
		huh.NewGroup(
			huh.NewText().Title("What would you like to write?").Value(&content),
		),
	)

	if err := form.Run(); err != nil {
		return fmt.Errorf("form submission failed: %v", err)
	}

	home, err := os.UserHomeDir()

	if err != nil {
		return fmt.Errorf("could not get home directory: %v", err)
	}

	directory := home + "/.iris/letters"
	if err := os.MkdirAll(directory, 0755); err != nil {
		return fmt.Errorf("could not create letters directory: %v", err)
	}

	filename := fmt.Sprintf("%s/%s.md", directory, date)

	file, err := os.Create(filename)

	if err != nil {
		return fmt.Errorf("could not create letter file: %v", err)
	}

	file.Write([]byte(fmt.Sprintf("# %s\n\nDate: %s\n\n%s", title, date, content)))
	file.Close()

	fmt.Printf("Your letter has been saved")

	return nil
}

func list(cmd *cobra.Command, args []string) error {
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

func read(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("please provide a date in the format YYYY-MM-DD")
	}

	date := args[0]
	if _, err := time.Parse("2006-01-02", date); err != nil {
		return fmt.Errorf("invalid date format, please use YYYY-MM-DD")
	}
	if date > time.Now().Format("2006-01-02") {
		return fmt.Errorf("you cannot read a letter for a future date")
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("could not get home directory: %v", err)
	}

	if _, err := os.Stat(home + "/.iris/letters/" + date + ".md"); os.IsNotExist(err) {
		return fmt.Errorf("no letter found for the date %s", date)
	}

	directory := home + "/.iris/letters"
	filename := fmt.Sprintf("%s/%s.md", directory, date)

	content, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("could not read letter file: %v", err)
	}

	renderedContent, err := glamour.Render(string(content), "dark")
	if err != nil {
		return fmt.Errorf("could not render letter content: %v", err)
	}

	fmt.Println(renderedContent)

	return nil
}

func main() {
	rootCmd.AddCommand(writeCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(readCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}
