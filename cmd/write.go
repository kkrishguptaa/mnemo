package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/kkrishguptaa/iris/util"
	"github.com/spf13/cobra"
)

// writeCmd represents the write command
var writeCmd = &cobra.Command{
	Use:   "write",
	Short: "Write a letter to your future self",
	Long:  `Write a letter to your future self. You can specify the title, date, and content of the letter.`,
	RunE:  write,
	Args:  cobra.NoArgs,
}

func init() {
	rootCmd.AddCommand(writeCmd)
}

func write(cmd *cobra.Command, args []string) error {
	var (
		date    string
		month   uint8
		year    string
		Title   string
		Content string
	)

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("What would be the title of your letter?").CharLimit(50).Value(&Title),
			huh.NewInput().Title("What date would you like to write this letter for? (1-31)").Validate(func(s string) error {
				return util.ErrorOnlyHandler(time.Parse("2006-01-02", fmt.Sprintf("2020-12-%s", s)))
			}).Value(&date),
			huh.NewSelect[uint8]().Options(
				huh.NewOption("January", uint8(1)),
				huh.NewOption("February", uint8(2)),
				huh.NewOption("March", uint8(3)),
				huh.NewOption("April", uint8(4)),
				huh.NewOption("May", uint8(5)),
				huh.NewOption("June", uint8(6)),
				huh.NewOption("July", uint8(7)),
				huh.NewOption("August", uint8(8)),
				huh.NewOption("September", uint8(9)),
				huh.NewOption("October", uint8(10)),
				huh.NewOption("November", uint8(11)),
				huh.NewOption("December", uint8(12)),
			).Title("What month would you like to write this letter for?").Value(&month),
			huh.NewInput().Title("What year would you like to write this letter for?").Validate(func(s string) error {
				return util.ErrorOnlyHandler(time.Parse("2006-01-02", fmt.Sprintf("%s-01-01", s)))
			}).Value(&year),
		),
		huh.NewGroup(
			huh.NewText().Title("What would you like to write?").Value(&Content),
		),
	)

	util.ErrorPrinter(form.Run())

	Date := util.ErrorHandler(time.Parse("2006-01-02", fmt.Sprintf("%s-%02d-%s", year, month, date)))

	if Date.Before(time.Now()) {
		util.ErrorPrinter(fmt.Errorf("you cannot write a letter for a past date: %s", Date.Format("2006-01-02")))
	}

	home := util.ErrorHandler(os.UserHomeDir())
	directory := home + "/.iris/letters"

	util.ErrorPrinter((os.MkdirAll(directory, 0755)))

	filename := fmt.Sprintf("%s/%s.md", directory, Date.Format("2006-01-02"))

	file := util.ErrorHandler(os.Create(filename))

	file.WriteString(fmt.Sprintf("# %s\n\n%s\n\n---\n\nWritten for: %s", Title, Content, Date.Format("2006-01-02")))

	defer util.ErrorPrinter(file.Close())

	util.SuccessPrinter("Letter saved successfully!")

	return nil
}
