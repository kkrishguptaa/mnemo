package write

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

func Write(cmd *cobra.Command, args []string) error {
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
