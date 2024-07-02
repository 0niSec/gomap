package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/urfave/cli/v2"
)

// bannerStyle is the style for the ASCII banner using [github.com/charmbracelet/lipgloss]
var bannerStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#00B9E8"))

// PrintBanner prints the ASCII banner for the program with other information
func PrintBanner() {
	fmt.Println(bannerStyle.Render(`
 ██████╗  ██████╗ ███╗   ███╗ █████╗ ██████╗ 
██╔════╝ ██╔═══██╗████╗ ████║██╔══██╗██╔══██╗
██║  ███╗██║   ██║██╔████╔██║███████║██████╔╝
██║   ██║██║   ██║██║╚██╔╝██║██╔══██║██╔═══╝ 
╚██████╔╝╚██████╔╝██║ ╚═╝ ██║██║  ██║██║     
 ╚═════╝  ╚═════╝ ╚═╝     ╚═╝╚═╝  ╚═╝╚═╝     

The lightweight, thin Go port scanning tool

💻 https://github.com/0niSec/gomap
🐳 docker pull 0nisec/gomap:latest
	`))
}

// ScanInfo prints the scan information based on the CLI flags used
func ScanInfo(c *cli.Context) {
	// TODO: Align the output
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("[+] Target: %s\n", c.String("target"))
	if c.String("ports") == "" {
		fmt.Println("[+] Ports: Top 1000")
	} else {
		fmt.Printf("[+] Ports: %s\n", c.String("ports"))
	}
	fmt.Printf("[+] Timeout: %s\n", c.Duration("timeout"))
	fmt.Println(strings.Repeat("=", 80))
}
