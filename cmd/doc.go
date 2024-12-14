/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/CrossEvol/setup/assets"
	"github.com/charmbracelet/huh"
	"log"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strings"

	_ "github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

// docCmd represents the doc command
var docCmd = &cobra.Command{
	Use:   "doc [tool]?",
	Short: "Open documentation for a specific tool",
	Long: `Open the documentation website for a specific development tool.
Available tools: eslint, tslint, prettier, vitest and so on.

Example usage:
  yourapp doc eslint or yourapp doc , if you pass the name, it will open the corresponding doc. If not, it will provide you multi selections.`,
	Run: func(cmd *cobra.Command, args []string) {
		var docPairs map[string]string
		data := assets.DocPairs
		if err := json.Unmarshal(data, &docPairs); err != nil {
			log.Fatalf("Failed to load Doc options: %v", err)
		}

		if len(args) == 1 {
			tool := args[0]
			url, ok := docPairs[strings.ToLower(tool)]
			if !ok {
				fmt.Printf("Unknown tool: %s\n", tool)
				return
			}
			openBrowser(url)
			return
		}

		var options []huh.Option[string]
		for key, value := range docPairs {
			option := huh.NewOption(strings.ToUpper(key[:1])+key[1:], fmt.Sprintf("%s ---> %s", key, value))
			options = append(options, option)
		}

		sort.Slice(options, func(i, j int) bool {
			return strings.Compare(options[i].Key, options[j].Key) < 0
		})

		var selections []string
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewMultiSelect[string]().Title("Docs for Tool-chain").
					Options(
						options...,
					).
					Description("Open the doc for your chosen tools: ").
					Value(&selections).
					Limit(10),
			),
		)

		err := form.Run()

		if err != nil {
			fmt.Println("Uh oh:", err)
			os.Exit(1)
		}

		for _, selection := range selections {
			split := strings.Split(selection, "--->")
			openBrowser(strings.TrimSpace(split[1]))
		}

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
