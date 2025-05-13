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

// eslintCmd represents the eslint command
var eslintCmd = &cobra.Command{
	Use:   "eslint",
	Short: "Set up ESLint for your project",
	Long: `This command sets up ESLint for your project.

It installs the necessary dependencies, creates an ESLint configuration file,
and adds a lint script to your package.json. This helps maintain code quality
and consistency in your JavaScript and TypeScript projects.`,
	Run: func(cmd *cobra.Command, args []string) {
		setupEslint()
	},
}

func setupEslint() {
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

	// Define package managers and their commands
	packageManagers := []struct {
		name       string
		installCmd []string
	}{
		{"pnpm", []string{"pnpm", "add", "-D", "eslint", "globals", "@eslint/js", "typescript-eslint"}},
		{"npm", []string{"npm", "install", "--save-dev", "eslint", "globals", "@eslint/js", "typescript-eslint"}},
		{"yarn", []string{"yarn", "add", "--dev", "eslint", "globals", "@eslint/js", "typescript-eslint"}},
		{"bun", []string{"bun", "add", "--dev", "eslint", "globals", "@eslint/js", "typescript-eslint"}},
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
		// Decide if you want to return or try a default
		// For now, let's try to proceed with a default or let it fail if no command is set
		// Or, more safely, return:
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
			fmt.Printf("Error installing ESLint with %s: %v\n", foundPackageManager, err)
		} else {
			fmt.Println("ESLint and dependencies installed successfully.")
		}
	} else if foundPackageManager != "" {
		// This case should ideally not be reached if packageManagers struct is well-defined
		fmt.Printf("Package manager %s found, but no install command configured for ESLint.\n", foundPackageManager)
	}
	// If no package manager was found and we didn't return, the original pnpm command would have run.
	// The logic above replaces that.

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
			scripts["lint"] = "eslint . --fix"

			updatedData, err := json.MarshalIndent(pkgJSON, "", "  ")
			if err != nil {
				fmt.Printf("Error marshalling updated package.json: %v\n", err)
			} else {
				err = os.WriteFile(packageJSONPath, updatedData, 0644)
				if err != nil {
					fmt.Printf("Error writing updated package.json: %v\n", err)
				} else {
					fmt.Println("'lint' script added/updated in package.json.")
				}
			}
		}
	}
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
