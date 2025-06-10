package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/huh"
	"github.com/kkrishguptaa/iris/util"
	"github.com/spf13/cobra"
)

// listCommand represents the write command
var readCommand = &cobra.Command{
	Use:   "read",
	Short: "Read a letter you have written",
	Long:  `Read a letter you have written. You get a select menu of all the letters you have written, and you can select one to read.You can only read letters that are not set in the future.`,
	RunE:  read,
	Args:  cobra.NoArgs,
}

func init() {
	rootCmd.AddCommand(readCommand)
}

func read(cmd *cobra.Command, args []string) error {
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
		date := util.ErrorHandler(time.Parse("2006-01-02", strings.Split(file.Name(), ".")[0]))

		if date.After(time.Now()) {
			continue // Skip letters set in the future
		}

		text := util.ErrorHandler(os.ReadFile(util.LettersPath + "/" + file.Name()))

		title := strings.Split(string(text), "\n")[0]
		title = strings.TrimPrefix(title, "# ")
		letters = append(letters, date.Format("2006-01-02")+": "+title)
	}

	var letter string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Options(huh.NewOptions(letters...)...).
				Title("Select a letter to read").
				Value(&letter),
		),
	)

	util.ErrorPrinter(form.Run())

	if letter == "" {
		util.ErrorPrinter(fmt.Errorf("no letter selected"))
	}

	letterDate := strings.Split(letter, ":")[0]

	letterFile := util.LettersPath + "/" + letterDate + ".md"

	if os.IsNotExist(util.ErrorOnlyHandler(os.Stat(letterFile))) {
		util.ErrorPrinter(fmt.Errorf("the letter for %s does not exist", letterDate))
	}

	content := util.ErrorHandler(os.ReadFile(letterFile))
	out := util.ErrorHandler(glamour.Render(string(content), "dark"))

	fmt.Println(out)

	return nil
}
