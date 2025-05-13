/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	"os"
	"slices"
)

const ESLINT = "eslint"
const PRETTIER = "prettier"
const VITEST = "vitest"
const HUSKY = "husky"

// esCmd represents the lint command
var esCmd = &cobra.Command{
	Use:   "es",
	Short: "Set up Javascript project tool-chains",
	Long: `Set up Javascript project tool-chains, include eslint, prettier, vitest, husky and so on.
It will not only install the needed packages, but also initialize the configuration files and add corresponding scripts.`,
	Run: func(cmd *cobra.Command, args []string) {
		var tools []string

		fmt.Println("ECMA script project initializing....")
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewMultiSelect[string]().Title("Tool Chains").
					Options(
						huh.NewOption(ESLINT, ESLINT),
						huh.NewOption(PRETTIER, PRETTIER),
						huh.NewOption(VITEST, VITEST),
						huh.NewOption(HUSKY, HUSKY),
					).
					Description("Choose your Tools").
					Value(&tools),
			),
		)

		err := form.Run()

		if err != nil {
			fmt.Println("Uh oh:", err)
			os.Exit(1)
		}

		fmt.Printf("Choose tools: %s\n", tools)

		configureEslintWithPrettier := false
		if len(tools) == 2 && slices.Contains(tools, ESLINT) && slices.Contains(tools, PRETTIER) {
			configureEslintWithPrettier = true
		}
		if configureEslintWithPrettier {
			fmt.Println()
			fmt.Println("=============== Setup Eslint with Prettier BEGIN  =====================")
			setupLinter()
			fmt.Println("=============== Setup Eslint with Prettier END  =====================")
			fmt.Println()
		}
		if slices.Contains(tools, ESLINT) {
			if configureEslintWithPrettier {
				skip()
			} else {
				fmt.Println()
				fmt.Println("=============== Setup Eslint BEGIN  =====================")
				setupLinter()
				fmt.Println("=============== Setup Eslint END  =====================")
				fmt.Println()
			}

		}
		if slices.Contains(tools, PRETTIER) {
			if configureEslintWithPrettier {
				skip()
			} else {
				fmt.Println()
				fmt.Println("=============== Setup Prettier BEGIN  =====================")
				setupPrettier()
				fmt.Println("=============== Setup Prettier END  =====================")
				fmt.Println()
			}
		}
		if slices.Contains(tools, VITEST) {
			fmt.Println()
			fmt.Println("=============== Setup Vitest BEGIN  =====================")
			setupVitest()
			fmt.Println()
			fmt.Println("=============== Setup Vitest END  =====================")
			fmt.Println()
		}
		if slices.Contains(tools, HUSKY) {
			fmt.Println()
			fmt.Println("=============== Setup Husky BEGIN  =====================")
			setupHusky()
			fmt.Println()
			fmt.Println("=============== Setup Husky END  =====================")
			fmt.Println()
		}
	},
}

func init() {
	rootCmd.AddCommand(esCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// esCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// esCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
