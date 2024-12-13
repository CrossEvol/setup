/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/CrossEvol/setup/common"
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
	const pnpmCommand = `pnpm add --save-dev --save-exact prettier`
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

	// Run pnpm command
	cmd2 := exec.Command("pnpm", "add", "--save-dev", "--save-exact", "prettier")
	cmd2.Stdout = os.Stdout
	cmd2.Stderr = os.Stderr
	err = cmd2.Run()
	if err != nil {
		fmt.Printf("Error running pnpm command: %v\n", err)
	} else {
		fmt.Println("Prettier installed successfully")
	}

	newEntries := `
		"prettier": "npx prettier . --write",`
	common.UpdatePackageJSON(newEntries)
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
