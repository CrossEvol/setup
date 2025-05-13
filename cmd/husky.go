/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// huskyCmd represents the husky command
var huskyCmd = &cobra.Command{
	Use:   "husky",
	Short: "Set up Husky for Git hooks",
	Long: `This command sets up Husky for managing Git hooks in your project.

It detects your package manager (pnpm, npm, yarn, or bun), installs Husky,
and runs the initialization command to set up the .husky directory and the
prepare script in package.json.`,
	Run: func(cmd *cobra.Command, args []string) {
		setupHusky()
	},
}

func setupHusky() {
	fmt.Println("Setting up Husky...")

	// Define package managers and their commands
	// Note: Yarn init is different, skipping for now based on docs
	packageManagers := []struct {
		name       string
		installCmd []string
		initCmd    []string
	}{
		{"pnpm", []string{"pnpm", "add", "--save-dev", "husky"}, []string{"pnpm", "exec", "husky", "init"}},
		{"npm", []string{"npm", "install", "--save-dev", "husky"}, []string{"npx", "husky", "init"}},
		// {"yarn", []string{"yarn", "add", "--dev", "husky"}, []string{"yarn", "husky", "init"}}, // Yarn init requires manual steps
		{"bun", []string{"bun", "add", "--dev", "husky"}, []string{"bunx", "husky", "init"}},
	}

	var foundPackageManager string
	var installCmdArgs []string
	var initCmdArgs []string

	// Check for package managers in order
	for _, pm := range packageManagers {
		// Check if the package manager command exists
		_, err := exec.LookPath(pm.name)
		if err == nil {
			foundPackageManager = pm.name
			installCmdArgs = pm.installCmd
			initCmdArgs = pm.initCmd
			fmt.Printf("Found package manager: %s\n", foundPackageManager)
			break // Use the first one found
		}
	}

	if foundPackageManager == "" {
		fmt.Println("Error: No supported package manager (pnpm, npm, yarn, bun) found.")
		fmt.Println("Please install one of these package managers and try again.")
		return
	}

	// 1. Install Husky
	fmt.Printf("Running installation command: %s %v\n", installCmdArgs[0], installCmdArgs[1:])
	installCmd := exec.Command(installCmdArgs[0], installCmdArgs[1:]...)
	installCmd.Stdout = os.Stdout
	installCmd.Stderr = os.Stderr
	err := installCmd.Run()
	if err != nil {
		fmt.Printf("Error installing Husky with %s: %v\n", foundPackageManager, err)
		return
	}
	fmt.Println("Husky installed successfully.")

	// 2. Run husky init
	// Special handling for yarn init if we decide to add it later, but for now, skip.
	if foundPackageManager == "yarn" {
		fmt.Println("Skipping husky init for yarn. Please refer to Husky documentation for manual setup.")
		return
	}

	fmt.Printf("Running init command: %s %v\n", initCmdArgs[0], initCmdArgs[1:])
	initCmd := exec.Command(initCmdArgs[0], initCmdArgs[1:]...)
	initCmd.Stdout = os.Stdout
	initCmd.Stderr = os.Stderr
	err = initCmd.Run()
	if err != nil {
		fmt.Printf("Error running husky init with %s: %v\n", foundPackageManager, err)
		return
	}
	fmt.Println("Husky initialized successfully.")

	fmt.Println("Husky setup complete.")
}

func init() {
	rootCmd.AddCommand(huskyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// huskyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags, which will only run when this command
	// is called directly, e.g.:
	// huskyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
