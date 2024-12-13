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

// linterCmd represents the linter command
var linterCmd = &cobra.Command{
	Use:   "linter",
	Short: "Set up ESLint and Prettier together",
	Long: `This command sets up both ESLint and Prettier for your project.

It installs ESLint, Prettier, and related plugins, creates configuration files,
and adds lint and format scripts to your package.json. This combined setup
ensures both code quality and consistent formatting in your project.`,
	Run: func(cmd *cobra.Command, args []string) {

		setupLinter()
	},
}

func setupLinter() {
	const eslintFile = `eslint.config.mjs`
	const pnpmCommand = `pnpm install --save-dev eslint globals @eslint/js typescript-eslint && pnpm add --save-dev --save-exact prettier eslint-config-prettier eslint-plugin-prettier `
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
	const eslintConfig = `
import globals from 'globals'
import pluginJs from '@eslint/js'
import tseslint from 'typescript-eslint'
import prettierConfig from 'eslint-config-prettier'
import prettierPlugin from 'eslint-plugin-prettier'

export default [
    {
        files: ['**/*.{js,mjs,cjs,ts}'],
        languageOptions: {
            globals: {
                ...globals.browser,
            },
        },
        plugins: {
            prettier: prettierPlugin,
        },
        rules: {
            'no-unused-vars': 'warn',
            'no-undef': 'warn',
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
            ],
            'prettier/prettier': [
                'error',
                {
                    singleQuote: true,
                    semi: false,
                    tabWidth: 4,
                },
            ],
        },
    },
    pluginJs.configs.recommended,
    ...tseslint.configs.recommended,
    prettierConfig,
]
`
	fmt.Println("linter called")

	// Create eslint.config.mjs file
	err := os.WriteFile(eslintFile, []byte(eslintConfig), 0644)
	if err != nil {
		fmt.Printf("Error creating %s: %v\n", eslintFile, err)
	} else {
		fmt.Printf("%s created successfully\n", eslintFile)
	}

	// Create .prettierignore file
	err = os.WriteFile(ignoreFile, []byte(prettierIgnoreConfig), 0644)
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
	cmd2 := exec.Command("pnpm", "install", "--save-dev", "eslint", "globals", "@eslint/js", "typescript-eslint")
	cmd2.Stdout = os.Stdout
	cmd2.Stderr = os.Stderr
	err = cmd2.Run()
	if err != nil {
		fmt.Printf("Error running first pnpm command: %v\n", err)
	} else {
		fmt.Println("ESLint and dependencies installed successfully")
	}

	cmd2 = exec.Command("pnpm", "add", "--save-dev", "--save-exact", "prettier", "eslint-config-prettier", "eslint-plugin-prettier")
	cmd2.Stdout = os.Stdout
	cmd2.Stderr = os.Stderr
	err = cmd2.Run()
	if err != nil {
		fmt.Printf("Error running second pnpm command: %v\n", err)
	} else {
		fmt.Println("Prettier and related ESLint plugins installed successfully")
	}

	newEntries := `
		"lint": "eslint . --fix",`
	common.UpdatePackageJSON(newEntries)
}

func init() {
	rootCmd.AddCommand(linterCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// linterCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// linterCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
