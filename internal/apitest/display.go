package apitest

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	cellStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Padding(0, 1)

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF4365"))

	warningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFA400"))

	borderStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#3C3836"))

	totalStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1)
)

type TableRow struct {
	Endpoint  string
	Method    string
	Response  string
	Assertion string
}

func DisplayTable(rows []TableRow) {
	// Sort rows by Endpoint and Method
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].Endpoint == rows[j].Endpoint {
			return rows[i].Method < rows[j].Method
		}
		return rows[i].Endpoint < rows[j].Endpoint
	})

	var (
		maxEndpoint  = 60
		maxMethod    = 15
		maxResponse  = 30
		maxAssertion = 40
	)

	table := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240"))

	header := lipgloss.JoinHorizontal(lipgloss.Top,
		headerStyle.Width(maxEndpoint).Render("Endpoint"),
		borderStyle.Render("│"),
		headerStyle.Width(maxMethod).Render("Method"),
		borderStyle.Render("│"),
		headerStyle.Width(maxResponse).Render("Response"),
		borderStyle.Render("│"),
		headerStyle.Width(maxAssertion).Render("Assertion"),
	)

	var renderedRows []string
	for _, row := range rows {
		renderedRow := lipgloss.JoinHorizontal(lipgloss.Top,
			cellStyle.Width(maxEndpoint).Render(truncate(row.Endpoint, maxEndpoint)),
			borderStyle.Render("│"),
			cellStyle.Width(maxMethod).Render(row.Method),
			borderStyle.Render("│"),
			cellStyle.Width(maxResponse).Render(truncate(row.Response, maxResponse)),
			borderStyle.Render("│"),
			renderMultilineAssertion(row.Assertion, maxAssertion),
		)
		renderedRows = append(renderedRows, renderedRow)
	}

	renderedTable := lipgloss.JoinVertical(lipgloss.Left,
		header,
		strings.Join(renderedRows, "\n"),
	)

	totalEndpoints := fmt.Sprintf("Total Endpoints Tested: %d", len(rows))
	renderedTotal := totalStyle.Render(totalEndpoints)

	fmt.Println(table.Render(renderedTable))
	fmt.Println(renderedTotal)
}

func renderMultilineAssertion(assertion string, width int) string {
	style := cellStyle.Copy().Width(width)
	words := strings.Fields(assertion)
	lines := []string{}
	currentLine := ""

	for _, word := range words {
		if len(currentLine)+len(word)+1 <= width {
			if currentLine != "" {
				currentLine += " "
			}
			currentLine += word
		} else {
			lines = append(lines, currentLine)
			currentLine = word
		}
	}
	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	switch {
	case strings.Contains(assertion, "PASS"):
		return successStyle.Inherit(style).Render(strings.Join(lines, "\n"))
	case strings.Contains(assertion, "FAIL"):
		return errorStyle.Inherit(style).Render(strings.Join(lines, "\n"))
	default:
		return warningStyle.Inherit(style).Render(strings.Join(lines, "\n"))
	}
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}
