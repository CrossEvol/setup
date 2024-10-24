/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/CrossEvol/setup/assets"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// gitignoreCmd represents the gitignore command
var gitignoreCmd = &cobra.Command{
	Use:   "ignore",
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

		if all {
			fmt.Println("Available keys:")
			fmt.Println("=====================================================>")
			for key := range gitignorePairs {
				fmt.Println(key)
			}
			return
		}

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

		if len(args) == 0 {
			log.Fatal("Please provide a language or framework name.")
		}

		targetKey := strings.ToLower(args[0])

		for key, value := range gitignorePairs {
			if strings.ToLower(key) == targetKey {
				url := fmt.Sprintf("https://raw.githubusercontent.com/github/gitignore/main/%s", value)
				resp, err := http.Get(url)
				if err != nil {
					log.Fatalf("Failed to get URL: %v", err)
				}
				defer resp.Body.Close()

				if resp.StatusCode == http.StatusOK {
					body, err := ioutil.ReadAll(resp.Body)
					if err != nil {
						log.Fatalf("Failed to read response body: %v", err)
					}

					// Write the content to a file
					if err := ioutil.WriteFile(value, body, 0644); err != nil {
						log.Fatalf("Failed to write file: %v", err)
					}
					fmt.Printf("Content saved to file: %s\n", value)
				} else {
					log.Fatalf("Failed to fetch content, status code: %d", resp.StatusCode)
				}
				return
			}
		}
		log.Printf("No matching key found for: %s", targetKey)
	},
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
