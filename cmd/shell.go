/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
)

// shellCmd represents the shell command
var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Add shell header",
	Long:  `for tab completion in terminal, provide the shell header is necessary, should provide "#!/bin/bash" for suffix of .sh , "#!/usr/bin/env pwsh" for suffix of .ps1`,
	Run: func(cmd *cobra.Command, args []string) {
		var nothingTodo = true

		// Get current directory
		currentDir, err := os.Getwd()
		if err != nil {
			fmt.Printf("Error getting current directory: %v\n", err)
			return
		}

		// Walk through directory
		err = filepath.Walk(currentDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Skip if not a file
			if info.IsDir() {
				return nil
			}

			// Check file extensions
			var shebang string
			switch {
			case strings.HasSuffix(info.Name(), ".sh"):
				shebang = "#!/bin/bash"
			case strings.HasSuffix(info.Name(), ".ps1"):
				shebang = "#!/usr/bin/env pwsh"
			default:
				return nil
			}

			// Read file content
			content, err := os.ReadFile(path)
			if err != nil {
				fmt.Printf("Error reading file %s: %v\n", path, err)
				return nil
			}

			// Get the first line
			lines := strings.Split(string(content), "\n")
			if len(lines) == 0 || lines[0] != shebang {
				// Create new content with shebang
				newContent := shebang
				if len(lines) > 0 {
					// Add original content after shebang
					newContent += "\n" + strings.Join(lines, "\n")
				}

				// Write back to file
				err = os.WriteFile(path, []byte(newContent), info.Mode())
				if err != nil {
					fmt.Printf("Error writing to file %s: %v\n", path, err)
					return nil
				}
				fmt.Printf("Added shell header to %s\n", path)
				nothingTodo = false
			}

			return nil
		})

		if err != nil {
			fmt.Printf("Error walking directory: %v\n", err)
			return
		}

		if nothingTodo {
			fmt.Printf("Nothing to do!\n")
		}
	},
}

func init() {
	rootCmd.AddCommand(shellCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// shellCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// shellCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
