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

// eslintCmd represents the eslint command
var eslintCmd = &cobra.Command{
	Use:   "eslint",
	Short: "Set up ESLint for your project",
	Long: `This command sets up ESLint for your project.

It installs the necessary dependencies, creates an ESLint configuration file,
and adds a lint script to your package.json. This helps maintain code quality
and consistency in your JavaScript and TypeScript projects.`,
	Run: func(cmd *cobra.Command, args []string) {
		const pnpmCommand = `pnpm install --save-dev eslint globals @eslint/js typescript-eslint`
		const eslintFile = `eslint.config.mjs`
		const eslintConfig = `
import pluginJs from "@eslint/js";
import globals from "globals";
import tseslint from "typescript-eslint";

export default [
  {
    files: ["**/*.{js,mjs,cjs,ts}"], rules: {
      'no-unused-vars': 'error',
      'no-undef': 'error',
      '@typescript-eslint/no-unused-vars': [
        'error',
        {
          args: 'all',
          argsIgnorePattern: '^_',
          caughtErrors: 'all',
          caughtErrorsIgnorePattern: '^_',
          destructuredArrayIgnorePattern: '^_',
          varsIgnorePattern: '^_',
          ignoreRestSiblings: true,
        },
      ]
    },
  },
  { languageOptions: { globals: globals.node } },
  pluginJs.configs.recommended,
  ...tseslint.configs.recommended,
];
`
		fmt.Println("eslint called")

		// Create eslint.config.mjs file
		err := os.WriteFile(eslintFile, []byte(eslintConfig), 0644)
		if err != nil {
			fmt.Printf("Error creating %s: %v\n", eslintFile, err)
		} else {
			fmt.Printf("%s created successfully\n", eslintFile)
		}

		// Run pnpm command
		cmd2 := exec.Command("pnpm", "install", "--save-dev", "eslint", "globals", "@eslint/js", "typescript-eslint")
		cmd2.Stdout = os.Stdout
		cmd2.Stderr = os.Stderr
		err = cmd2.Run()
		if err != nil {
			fmt.Printf("Error running pnpm command: %v\n", err)
		} else {
			fmt.Println("ESLint and dependencies installed successfully")
		}

		newEntries := `
		"lint": "eslint . --fix",`
		common.UpdatePackageJSON(newEntries)
	},
}

func init() {
	rootCmd.AddCommand(eslintCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// eslintCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// eslintCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
