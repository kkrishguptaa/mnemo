package read

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/glamour"
	"github.com/spf13/cobra"
)

func Read(cmd *cobra.Command, args []string) error {
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
