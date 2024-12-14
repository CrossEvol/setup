/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/CrossEvol/setup/assets"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	"io"
	"log"
	"net/http"
	"os"
	"slices"
	"sort"
	"strings"
)

// gitignoreCmd represents the gitignore command
var gitignoreCmd = &cobra.Command{
	Use:   "ignore [lang]? [--all]? [--char]?",
	Short: "Generate a .gitignore file for a specified language or framework",
	Long: `This command is used to generate a .gitignore file for a specified language or framework.
It fetches the .gitignore template from GitHub's collection of .gitignore files and saves it locally.
For example, you can use this command to quickly set up a .gitignore file for your Go projects.`,
	Run: func(cmd *cobra.Command, args []string) {
		data := assets.GitignorePairs

		var gitignorePairs map[string]string
		if err := json.Unmarshal(data, &gitignorePairs); err != nil {
			log.Fatalf("Failed to parse JSON: %v", err)
		}

		// Check for flags
		all, _ := cmd.Flags().GetBool("all")
		char, _ := cmd.Flags().GetString("char")

		// list all options
		if all {
			fmt.Println("Available keys:")
			fmt.Println("=====================================================>")
			for key := range gitignorePairs {
				fmt.Println(key)
			}
			return
		}

		// list all options with prefix letter
		if char != "" {
			fmt.Printf("Available Keys starting with '%s':\n", char)
			fmt.Println("=====================================================>")
			for key := range gitignorePairs {
				if strings.HasPrefix(strings.ToLower(key), strings.ToLower(char)) {
					fmt.Println(key)
				}
			}
			return
		}

		// choose the language with the prefix letter
		if len(args) == 0 {
			fmt.Println("=====================================================>")

			var letters []string
			for i := 65; i < 91; i++ {
				ch := string(uint8(i))
				letters = append(letters, ch)
				letters = append(letters, strings.ToLower(ch))
			}

			var ch string

			form := huh.NewForm(
				huh.NewGroup(
					huh.NewInput().
						Value(&ch).
						Title("Input the first letter: [A-Za-z]").
						Placeholder("A").
						Validate(func(s string) error {
							if !slices.Contains(letters, ch) {
								return errors.New("unknown Character")
							}
							return nil
						}).
						Description("Then it will list all options start with it."),
				),
			)

			err := form.Run()

			if err != nil {
				fmt.Println("Uh oh:", err)
				os.Exit(1)
			}

			var options []huh.Option[string]
			for key := range gitignorePairs {
				if strings.ToLower(key)[:1] == strings.ToLower(ch) {
					option := huh.NewOption(key, fmt.Sprintf("%s", key))
					options = append(options, option)
				}
			}

			sort.Slice(options, func(i, j int) bool {
				return strings.Compare(options[i].Key, options[j].Key) < 0
			})

			var selection string
			form = huh.NewForm(
				huh.NewGroup(
					huh.NewSelect[string]().Title("Choose the Lang:").
						Options(
							options...,
						).
						Description("Download the .gitignore for the language").
						Value(&selection),
				),
			)

			err = form.Run()

			if err != nil {
				fmt.Println("Uh oh:", err)
				os.Exit(1)
			}

			download(gitignorePairs[selection])
			return
		}

		// pass the target language
		targetKey := strings.ToLower(args[0])
		for key, value := range gitignorePairs {
			if strings.ToLower(key) == targetKey {
				download(value)
				return
			}
		}

		log.Printf("No matching key found for: %s", targetKey)
	},
}

func download(value string) {
	url := fmt.Sprintf("https://raw.githubusercontent.com/github/gitignore/main/%s", value)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Failed to get URL: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Failed to read response body: %v", err)
		}

		// Write the content to a file
		if err := os.WriteFile(value, body, 0644); err != nil {
			log.Fatalf("Failed to write file: %v", err)
		}
		fmt.Printf("Content saved to file: %s\n", value)
	} else {
		log.Fatalf("Failed to fetch content, status code: %d", resp.StatusCode)
	}
}

func init() {
	rootCmd.AddCommand(gitignoreCmd)

	// Define flags
	gitignoreCmd.Flags().BoolP("all", "a", false, "Output all available keys")
	gitignoreCmd.Flags().StringP("char", "c", "", "Output keys starting with a specific character")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// gitignoreCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// gitignoreCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
