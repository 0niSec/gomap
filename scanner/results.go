package scanner

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	openStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00"))
	closedStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0800"))
	filteredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFEF00"))
)

func PrettyPrintScanResults(results map[uint16]string) {
	fmt.Println(lipgloss.JoinHorizontal(lipgloss.Left,
		lipgloss.NewStyle().Width(10).Render("PORT"),
		lipgloss.NewStyle().Width(10).Render("PROTOCOL"),
		lipgloss.NewStyle().Width(10).Render("STATE"),
		lipgloss.NewStyle().Width(10).Render("SERVICE"),
	))

	for port, status := range results {
		coloredStatus := getColoredStatus(status)
		fmt.Println(lipgloss.JoinHorizontal(lipgloss.Left,
			lipgloss.NewStyle().Width(10).Render(fmt.Sprintf("%d", port)),
			lipgloss.NewStyle().Width(10).Render("tcp"),
			lipgloss.NewStyle().Width(10).Render(coloredStatus),
			lipgloss.NewStyle().Width(10).Render(""),
		))
	}
	// Add a blank line to separate the results
	fmt.Println("")
}

func getColoredStatus(status string) string {
	switch status {
	case "open":
		return openStyle.Render(status)
	case "closed":
		return closedStyle.Render(status)
	case "filtered":
		return filteredStyle.Render(status)
	default:
		return status
	}
}

// func PrettyPrintScanResults(results map[uint16]string) {
// 	fmt.Printf("%-10s %-10s %-10s %-10s\n", "PORT", "PROTOCOL", "STATE", "SERVICE")
// 	for port, result := range results {
// 		coloredStatus := result
// 		switch result {
// 		case "open":
// 			coloredStatus = openStyle.Render(result)
// 		case "closed":
// 			coloredStatus = closedStyle.Render(result)
// 		case "filtered":
// 			coloredStatus = filteredStyle.Render(result)
// 		}
// 		fmt.Printf("%-10d %-10s %-10s %-10s\n", port, "tcp", coloredStatus, "http")
// 	}
// }
