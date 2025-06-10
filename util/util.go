package util

import (
	"encoding/json"
	"os"

	"github.com/charmbracelet/lipgloss"
)

func ErrorPrinter(err error) {
	if err != nil {
		errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true).BorderBottom(true).BorderStyle(lipgloss.NormalBorder())
		errorMessage := errorStyle.Render("Error: " + err.Error())
		println(errorMessage)
		os.Exit(1)
	}
}

func ErrorHandler[T any](result T, err error) T {
	if err != nil {
		ErrorPrinter(err)
		return result
	}
	return result
}

func ErrorOnlyHandler[T any](result T, err error) error {
	if err == nil {
		ErrorPrinter(err)
		return err
	}

	return nil
}

func SuccessPrinter(message string) {
	successStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Bold(true).BorderBottom(true).BorderStyle(lipgloss.NormalBorder())
	successMessage := successStyle.Render("Success: " + message)
	println(successMessage)
}

var home = ErrorHandler(os.UserHomeDir())

func FetchData() []string {
	folder := home + "/.mnemo"

	ErrorPrinter(os.MkdirAll(folder, 0755))

	filepath := folder + "/data.json"

	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		file := ErrorHandler(os.Create(filepath))

		file.Write([]byte("[]"))

		defer file.Close()
	}

	file := ErrorHandler(os.ReadFile(filepath))

	if len(file) == 0 {
		os.WriteFile(filepath, []byte("[]"), 0644)
	}

	var data []string

	ErrorPrinter(json.Unmarshal(file, &data))

	return data
}

func SaveData(data []byte) {
	folder := home + "/.mnemo"

	ErrorPrinter(os.MkdirAll(folder, 0755))

	filepath := folder + "/data.json"

	os.WriteFile(filepath, data, 0644)
	SuccessPrinter("Data saved successfully.")
}
