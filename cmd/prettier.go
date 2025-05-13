/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// prettierCmd represents the prettier command
var prettierCmd = &cobra.Command{
	Use:   "prettier",
	Short: "Set up Prettier for your project",
	Long: `This command sets up Prettier for your project.

It installs Prettier, creates configuration files (.prettierrc and .prettierignore),
and adds a format script to your package.json. Prettier helps maintain consistent
code formatting across your project.`,
	Run: func(cmd *cobra.Command, args []string) {
		setupPrettier()
	},
}

func setupPrettier() {
	const ignoreFile = `.prettierignore`
	const prettierFile = `.prettierrc`
	const prettierConfig = `
{
    "singleQuote": true,
    "semi": false,
    "tabWidth": 4,
    "plugins": []
}

`
	const prettierIgnoreConfig = `
# Ignore artifacts:
build
coverage
.next
node_modules

`
	fmt.Println("prettier called")

	// Create .prettierignore file
	err := os.WriteFile(ignoreFile, []byte(prettierIgnoreConfig), 0644)
	if err != nil {
		fmt.Printf("Error creating %s: %v\n", ignoreFile, err)
	} else {
		fmt.Printf("%s created successfully\n", ignoreFile)
	}

	// Create .prettierrc file
	err = os.WriteFile(prettierFile, []byte(prettierConfig), 0644)
	if err != nil {
		fmt.Printf("Error creating %s: %v\n", prettierFile, err)
	} else {
		fmt.Printf("%s created successfully\n", prettierFile)
	}

	// Define package managers and their commands
	packageManagers := []struct {
		name       string
		installCmd []string
	}{
		{"pnpm", []string{"pnpm", "add", "-D", "--save-exact", "prettier"}},
		{"npm", []string{"npm", "install", "--save-dev", "--save-exact", "prettier"}},
		{"yarn", []string{"yarn", "add", "--dev", "--exact", "prettier"}},
		{"bun", []string{"bun", "add", "--dev", "--exact", "prettier"}},
	}

	var foundPackageManager string
	var installCmdArgs []string

	// Check for package managers in order
	for _, pm := range packageManagers {
		_, err := exec.LookPath(pm.name)
		if err == nil {
			foundPackageManager = pm.name
			installCmdArgs = pm.installCmd
			fmt.Printf("Found package manager: %s\n", foundPackageManager)
			break // Use the first one found
		}
	}

	if foundPackageManager == "" {
		fmt.Println("Error: No supported package manager (pnpm, npm, yarn, bun) found.")
		fmt.Println("Please install one of these package managers and try again.")
		// return
	}

	// Run install command
	if len(installCmdArgs) > 0 {
		fmt.Printf("Running installation command: %s %v\n", installCmdArgs[0], installCmdArgs[1:])
		installCmd := exec.Command(installCmdArgs[0], installCmdArgs[1:]...)
		installCmd.Stdout = os.Stdout
		installCmd.Stderr = os.Stderr
		err = installCmd.Run()
		if err != nil {
			fmt.Printf("Error installing Prettier with %s: %v\n", foundPackageManager, err)
		} else {
			fmt.Println("Prettier installed successfully.")
		}
	} else if foundPackageManager != "" {
		fmt.Printf("Package manager %s found, but no install command configured for Prettier.\n", foundPackageManager)
	}

	// Update package.json
	packageJSONPath := "package.json"
	packageJSONData, err := os.ReadFile(packageJSONPath)
	if err != nil {
		fmt.Printf("Error reading package.json: %v\n", err)
	} else {
		var pkgJSON map[string]interface{}
		err = json.Unmarshal(packageJSONData, &pkgJSON)
		if err != nil {
			fmt.Printf("Error parsing package.json: %v\n", err)
		} else {
			scripts, ok := pkgJSON["scripts"].(map[string]interface{})
			if !ok {
				scripts = make(map[string]interface{})
				pkgJSON["scripts"] = scripts
			}
			scripts["prettier"] = "npx prettier . --write" // or use detected package manager's run command

			updatedData, err := json.MarshalIndent(pkgJSON, "", "  ")
			if err != nil {
				fmt.Printf("Error marshalling updated package.json: %v\n", err)
			} else {
				err = os.WriteFile(packageJSONPath, updatedData, 0644)
				if err != nil {
					fmt.Printf("Error writing updated package.json: %v\n", err)
				} else {
					fmt.Println("'prettier' script added/updated in package.json.")
				}
			}
		}
	}
}

func init() {
	rootCmd.AddCommand(prettierCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// prettierCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// prettierCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
