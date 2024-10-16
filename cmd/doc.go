/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"
)

var docMap = map[string]string{
	"eslint":   "https://eslint.org/",
	"tslint":   "https://typescript-eslint.io/",
	"prettier": "https://prettier.io/",
	"vitest":   "https://vitest.dev/",
}

// docCmd represents the doc command
var docCmd = &cobra.Command{
	Use:   "doc [tool]",
	Short: "Open documentation for a specific tool",
	Long: `Open the documentation website for a specific development tool.
Available tools: eslint, tslint, prettier, vitest.

Example usage:
  yourapp doc eslint`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tool := args[0]
		url, ok := docMap[tool]
		if !ok {
			fmt.Printf("Unknown tool: %s\n", tool)
			return
		}
		openBrowser(url)
	},
}

func init() {
	rootCmd.AddCommand(docCmd)
}

func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		fmt.Printf("Unsupported platform. Please open %s manually.\n", url)
		return
	}

	if err != nil {
		fmt.Printf("Error opening browser: %v\n", err)
	}
}
