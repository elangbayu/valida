package apitest

import (
	"fmt"
	"os"
	"strconv"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

func displayResultsAsTable(results []TestResult) {
	const (
		green       = lipgloss.Color("#009900")
		red         = lipgloss.Color("#f61901")
		mutedGreen  = lipgloss.Color("#729762")
		brightGreen = lipgloss.Color("#E7F0DC")
	)

	re := lipgloss.NewRenderer(os.Stdout)

	var (
		HeaderStyle       = re.NewStyle().Foreground(brightGreen).Bold(true).Align(lipgloss.Center)
		CellStyle         = re.NewStyle().Padding(0, 1).Width(14)
		PassedStatusStyle = re.NewStyle().Foreground(green).Width(14)
		FailedStatusStyle = re.NewStyle().Foreground(red).Width(14)
		BorderStyle       = lipgloss.NewStyle().Foreground(mutedGreen)
	)

	t := table.New().
		Border(lipgloss.ThickBorder()).
		BorderStyle(BorderStyle).
		Headers("ENDPOINT", "METHOD", "STATUS").
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == 0 {
				return HeaderStyle
			}
			if row-1 < len(results) && col == 2 {
				switch results[row-1].Status {
				case "PASSED":
					return PassedStatusStyle
				case "FAILED":
					return FailedStatusStyle
				}
			}
			return CellStyle
		})

	var passedCount, failedCount int

	for _, result := range results {
		t.Row(result.Endpoint, result.Method, result.Status)
		if result.Status == "PASSED" {
			passedCount++
		} else if result.Status == "FAILED" {
			failedCount++
		}
	}

	fmt.Println(t)

	passedCountStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(green)).
		MarginTop(1).
		Width(22)

	failedCountStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(red)).
		MarginTop(1).
		Width(22)

	fmt.Println(passedCountStyle.Render("PASSED: " + strconv.Itoa(passedCount)))
	fmt.Println(failedCountStyle.Render("FAILED: " + strconv.Itoa(failedCount)))
}
