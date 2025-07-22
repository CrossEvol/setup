/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

// copyCmd represents the copy command
var copyCmd = &cobra.Command{
	Use:   "copy",
	Short: "Copy .env to .env.example",
	Long: `Copy .env to .env.example.
If .env.example exists, you can choose to override or combine.`,
	Run: func(cmd *cobra.Command, args []string) {
		// check if .env exists
		if _, err := os.Stat(".env"); os.IsNotExist(err) {
			fmt.Println("Error: .env file not found")
			return
		}

		// check if .env.example exists
		if _, err := os.Stat(".env.example"); err == nil {
			var choice string
			form := huh.NewSelect[string]().
				Title("`.env.example` already exists. What do you want to do?").
				Options(
					huh.NewOption("Override", "override"),
					huh.NewOption("Combine", "combine"),
				).
				Value(&choice)
			err := form.Run()
			if err != nil {
				return
			}

			if choice == "override" {
				err := copyEnv(true)
				if err != nil {
					fmt.Println("Error:", err)
				}
			} else if choice == "combine" {
				err := copyEnv(false)
				if err != nil {
					fmt.Println("Error:", err)
				}
			}
		} else {
			err := copyEnv(true)
			if err != nil {
				fmt.Println("Error:", err)
			}
		}
	},
}

func copyEnv(override bool) error {
	envFile, err := os.Open(".env")
	if err != nil {
		return err
	}
	defer envFile.Close()

	var newLines []string
	if override {
		// Override: read .env and strip values
		scanner := bufio.NewScanner(envFile)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, "=") && !strings.HasPrefix(strings.TrimSpace(line), "#") {
				parts := strings.SplitN(line, "=", 2)
				newLines = append(newLines, parts[0]+"=")
			} else {
				newLines = append(newLines, line)
			}
		}
	} else {
		// Combine: read .env.example, then add missing keys from .env
		exampleFile, err := os.Open(".env.example")
		if err != nil {
			return err
		}
		defer exampleFile.Close()

		exampleKeys := make(map[string]bool)
		scanner := bufio.NewScanner(exampleFile)
		for scanner.Scan() {
			line := scanner.Text()
			newLines = append(newLines, line)
			if strings.Contains(line, "=") && !strings.HasPrefix(strings.TrimSpace(line), "#") {
				parts := strings.SplitN(line, "=", 2)
				exampleKeys[strings.TrimSpace(parts[0])] = true
			}
		}

		var additionalLines []string
		scanner = bufio.NewScanner(envFile)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, "=") && !strings.HasPrefix(strings.TrimSpace(line), "#") {
				parts := strings.SplitN(line, "=", 2)
				key := strings.TrimSpace(parts[0])
				if !exampleKeys[key] {
					additionalLines = append(additionalLines, key+"=")
					exampleKeys[key] = true // Prevent duplicates from .env itself
				}
			}
		}

		if len(additionalLines) > 0 {
			if len(newLines) > 0 && newLines[len(newLines)-1] != "" {
				newLines = append(newLines, "")
			}
			newLines = append(newLines, "")
			newLines = append(newLines, additionalLines...)
		}
	}

	output := strings.Join(newLines, "\n")
	err = os.WriteFile(".env.example", []byte(output), 0644)
	if err != nil {
		return err
	}

	fmt.Println(".env.example updated successfully")
	return nil
}

func init() {
	envCmd.AddCommand(copyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// copyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// copyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
