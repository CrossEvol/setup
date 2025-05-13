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
const COMMITLINT = "commitlint"
const LINTSTAGED = "lintStaged"
const RELEASEIT = "releaseIt"

// nodeCmd represents the node command
var nodeCmd = &cobra.Command{
	Use:   "node",
	Short: "Set up Node.js project tool-chains",
	Long: `Set up Node.js project tool-chains, include eslint, prettier, vitest, husky and so on.
It will not only install the needed packages, but also initialize the configuration files and add corresponding scripts.`,
	Run: func(cmd *cobra.Command, args []string) {
		var tools []string

		fmt.Println("Node.js project initializing....")
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewMultiSelect[string]().Title("Tool Chains").
					Options(
						huh.NewOption(ESLINT, ESLINT),
						huh.NewOption(PRETTIER, PRETTIER),
						huh.NewOption(VITEST, VITEST),
						huh.NewOption(HUSKY, HUSKY),
						huh.NewOption(COMMITLINT, COMMITLINT),
						huh.NewOption(LINTSTAGED, LINTSTAGED),
						huh.NewOption(RELEASEIT, RELEASEIT),
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

		// Handle ESLint and Prettier combined setup
		configureEslintWithPrettier := false
		if slices.Contains(tools, ESLINT) && slices.Contains(tools, PRETTIER) {
			// Check if ONLY ESLINT and PRETTIER are selected for the combined setup
			// If other tools are selected, we might want to run them individually
			// The original logic checked if len(tools) == 2. Let's stick to that for now.
			if len(tools) == 2 {
				configureEslintWithPrettier = true
			} else {
				// If other tools are selected along with ESLint and Prettier,
				// we should probably run them individually.
				configureEslintWithPrettier = false
			}
		}

		if configureEslintWithPrettier {
			fmt.Println()
			fmt.Println("=============== Setup Eslint with Prettier BEGIN  =====================")
			setupLinter() // Assuming setupLinter handles both
			fmt.Println("=============== Setup Eslint with Prettier END  =====================")
			fmt.Println()
		}

		// Handle individual tool setups, skipping if combined was handled
		if slices.Contains(tools, ESLINT) && !configureEslintWithPrettier {
			fmt.Println()
			fmt.Println("=============== Setup Eslint BEGIN  =====================")
			setupEslint() // Call the specific ESLint setup
			fmt.Println("=============== Setup Eslint END  =====================")
			fmt.Println()
		}

		if slices.Contains(tools, PRETTIER) && !configureEslintWithPrettier {
			fmt.Println()
			fmt.Println("=============== Setup Prettier BEGIN  =====================")
			setupPrettier() // Call the specific Prettier setup
			fmt.Println("=============== Setup Prettier END  =====================")
			fmt.Println()
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

		if slices.Contains(tools, COMMITLINT) {
			fmt.Println()
			fmt.Println("=============== Setup Commitlint BEGIN  =====================")
			setupCommitlint()
			fmt.Println()
			fmt.Println("=============== Setup Commitlint END  =====================")
			fmt.Println()
		}

		if slices.Contains(tools, LINTSTAGED) {
			fmt.Println()
			fmt.Println("=============== Setup Lint-Staged BEGIN  =====================")
			setupLintStaged()
			fmt.Println()
			fmt.Println("=============== Setup Lint-Staged END  =====================")
			fmt.Println()
		}

		if slices.Contains(tools, RELEASEIT) {
			fmt.Println()
			fmt.Println("=============== Setup Release-It BEGIN  =====================")
			setupReleaseIt()
			fmt.Println()
			fmt.Println("=============== Setup Release-It END  =====================")
			fmt.Println()
		}
	},
}

func init() {
	rootCmd.AddCommand(nodeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// nodeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// nodeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// Note: The actual setup functions (setupEslint, setupPrettier, setupVitest,
// setupHusky, setupLinter, setupCommitlint, setupLintStaged, setupReleaseIt)
// are assumed to be defined in other files within the 'cmd' package
// (e.g., cmd/eslint.go, cmd/prettier.go, etc.) and are accessible here.
// The skip() function is also assumed to be defined elsewhere, likely in root.go.
