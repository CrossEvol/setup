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

// vitestCmd represents the vitest command
var vitestCmd = &cobra.Command{
	Use:   "vitest",
	Short: "Set up Vitest for your project",
	Long: `This command sets up Vitest for your project.

It installs Vitest, creates configuration files (vitest.config.ts and vitest.setup.ts),
and adds a test script to your package.json. Vitest is a fast and lightweight testing
framework for Vite-based projects.`,
	Run: func(cmd *cobra.Command, args []string) {
		const pnpmCommand = `pnpm install --save-dev vitest`
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

		// Run pnpm command
		cmd2 := exec.Command("pnpm", "install", "--save-dev", "vitest")
		cmd2.Stdout = os.Stdout
		cmd2.Stderr = os.Stderr
		err = cmd2.Run()
		if err != nil {
			fmt.Printf("Error running pnpm command: %v\n", err)
		} else {
			fmt.Println("Vitest installed successfully")
		}

		newEntries := `
		"test": "vitest . ",`
		common.UpdatePackageJSON(newEntries)
	},
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
