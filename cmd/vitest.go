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

// vitestCmd represents the vitest command
var vitestCmd = &cobra.Command{
	Use:   "vitest",
	Short: "Set up Vitest for your project",
	Long: `This command sets up Vitest for your project.

It installs Vitest, creates configuration files (vitest.config.ts and vitest.setup.ts),
and adds a test script to your package.json. Vitest is a fast and lightweight testing
framework for Vite-based projects.`,
	Run: func(cmd *cobra.Command, args []string) {
		setupVitest()
	},
}

func setupVitest() {
	const vitestSetupFile = `vitest.setup.ts`
	const vitestConfigFile = `vitest.config.ts`
	const vitestConfig = `
import path from 'path'
import { defineConfig } from 'vitest/config'

export default defineConfig({
    resolve: {
        alias: {
            '@': path.join(__dirname, 'src'),
        },
    },
    test: {
        environment: 'node',
        setupFiles: ['./vitest.setup.ts'],
    },
})
`
	const vitestSetupConfig = `
import { afterEach } from 'vitest'

afterEach(() => {})

`
	fmt.Println("vitest called")
	// Create vitest.config.ts file
	err := os.WriteFile(vitestConfigFile, []byte(vitestConfig), 0644)
	if err != nil {
		fmt.Printf("Error creating %s: %v\n", vitestConfigFile, err)
	} else {
		fmt.Printf("%s created successfully\n", vitestConfigFile)
	}

	// Create vitest.setup.ts file
	err = os.WriteFile(vitestSetupFile, []byte(vitestSetupConfig), 0644)
	if err != nil {
		fmt.Printf("Error creating %s: %v\n", vitestSetupFile, err)
	} else {
		fmt.Printf("%s created successfully\n", vitestSetupFile)
	}

	// Define package managers and their commands
	packageManagers := []struct {
		name       string
		installCmd []string
	}{
		{"pnpm", []string{"pnpm", "add", "-D", "vitest"}},
		{"npm", []string{"npm", "install", "--save-dev", "vitest"}},
		{"yarn", []string{"yarn", "add", "--dev", "vitest"}},
		{"bun", []string{"bun", "add", "--dev", "vitest"}},
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
			fmt.Printf("Error installing Vitest with %s: %v\n", foundPackageManager, err)
		} else {
			fmt.Println("Vitest installed successfully.")
		}
	} else if foundPackageManager != "" {
		fmt.Printf("Package manager %s found, but no install command configured for Vitest.\n", foundPackageManager)
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
			scripts["test"] = "vitest . " // or use detected package manager's run command

			updatedData, err := json.MarshalIndent(pkgJSON, "", "  ")
			if err != nil {
				fmt.Printf("Error marshalling updated package.json: %v\n", err)
			} else {
				err = os.WriteFile(packageJSONPath, updatedData, 0644)
				if err != nil {
					fmt.Printf("Error writing updated package.json: %v\n", err)
				} else {
					fmt.Println("'test' script added/updated in package.json.")
				}
			}
		}
	}
}

func init() {
	rootCmd.AddCommand(vitestCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// vitestCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// vitestCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
