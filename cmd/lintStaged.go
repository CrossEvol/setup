/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

// lintStagedCmd represents the lintStaged command
var lintStagedCmd = &cobra.Command{
	Use:   "lintStaged",
	Short: "Set up lint-staged for your project",
	Long: `This command sets up lint-staged for your project.

It detects your package manager, installs lint-staged, creates a configuration file,
adds a 'pre-commit' script to package.json, and integrates with Husky if it's set up.
lint-staged runs linters and formatters on staged files before committing.`,
	Run: func(cmd *cobra.Command, args []string) {
		setupLintStaged()
	},
}

func setupLintStaged() {
	fmt.Println("Setting up lint-staged...")

	// Define package managers and their commands
	packageManagers := []struct {
		name       string
		installCmd []string
		runCmd     string // Command prefix for running scripts (e.g., "pnpm run", "npm run")
	}{
		{"pnpm", []string{"pnpm", "add", "--save-dev", "lint-staged"}, "pnpm run"},
		{"npm", []string{"npm", "install", "--save-dev", "lint-staged"}, "npm run"},
		{"yarn", []string{"yarn", "add", "--dev", "lint-staged"}, "yarn run"},
		{"bun", []string{"bun", "add", "--dev", "lint-staged"}, "bun run"},
	}

	var foundPackageManager string
	var installCmdArgs []string
	var runCmdPrefix string

	// Check for package managers in order
	for _, pm := range packageManagers {
		_, err := exec.LookPath(pm.name)
		if err == nil {
			foundPackageManager = pm.name
			installCmdArgs = pm.installCmd
			runCmdPrefix = pm.runCmd
			fmt.Printf("Found package manager: %s\n", foundPackageManager)
			break // Use the first one found
		}
	}

	if foundPackageManager == "" {
		fmt.Println("Error: No supported package manager (pnpm, npm, yarn, bun) found.")
		fmt.Println("Please install one of these package managers and try again.")
		return
	}

	// 1. Install lint-staged
	fmt.Printf("Running installation command: %s %v\n", installCmdArgs[0], installCmdArgs[1:])
	installCmd := exec.Command(installCmdArgs[0], installCmdArgs[1:]...)
	installCmd.Stdout = os.Stdout
	installCmd.Stderr = os.Stderr
	err := installCmd.Run()
	if err != nil {
		fmt.Printf("Warning: Error installing lint-staged with %s: %v\n", foundPackageManager, err)
		fmt.Println("Attempting to continue assuming lint-staged is already installed.")
		// Continue setup even if installation fails, maybe it's already installed
	} else {
		fmt.Println("lint-staged installed successfully.")
	}

	// 2. Create lint-staged.config.js file
	const lintStagedConfigFile = `lint-staged.config.js`
	// Use the rules provided by the user, formatted for the .js config file
	lintStagedConfigContentTemplate := `
/** @type {import('./lib/types').Configuration} */
export default {
  'src/**/*.{js,jsx,ts,tsx,json}': [
    '%s lint', // Use the determined run command prefix
    '%s format' // Assuming 'format' script exists (e.g., prettier)
  ]
}
`
	// Format the content with the actual run command prefix
	configContent := fmt.Sprintf(lintStagedConfigContentTemplate, runCmdPrefix, runCmdPrefix)

	fmt.Printf("Creating %s...\n", lintStagedConfigFile)
	err = os.WriteFile(lintStagedConfigFile, []byte(configContent), 0644)
	if err != nil {
		fmt.Printf("Error creating %s: %v\n", lintStagedConfigFile, err)
		// Continue setup
	} else {
		fmt.Printf("%s created successfully.\n", lintStagedConfigFile)
	}

	// 3. Add "pre-commit": "lint-staged" script to package.json
	packageJSONPath := "package.json"
	packageJSONData, err := os.ReadFile(packageJSONPath)
	if err != nil {
		fmt.Printf("Error reading package.json: %v\n", err)
		// Cannot update package.json, stop here for this part
	} else {
		var pkgJSON map[string]interface{}
		err = json.Unmarshal(packageJSONData, &pkgJSON)
		if err != nil {
			fmt.Printf("Error parsing package.json: %v\n", err)
			// Cannot update package.json, stop here for this part
		} else {
			scripts, ok := pkgJSON["scripts"].(map[string]interface{})
			if !ok {
				// Scripts key doesn't exist or isn't a map, create it
				scripts = make(map[string]interface{})
				pkgJSON["scripts"] = scripts
			}
			scripts["pre-commit"] = "lint-staged"

			updatedData, err := json.MarshalIndent(pkgJSON, "", "  ")
			if err != nil {
				fmt.Printf("Error marshalling updated package.json: %v\n", err)
				// Cannot update package.json, stop here for this part
			} else {
				err = os.WriteFile(packageJSONPath, updatedData, 0644)
				if err != nil {
					fmt.Printf("Error writing updated package.json: %v\n", err)
				} else {
					fmt.Println("'pre-commit' script added/updated in package.json.")
				}
			}
		}
	}

	// 4. Integrate with Husky if .husky/pre-commit exists
	huskyPreCommitPath := filepath.Join(".husky", "pre-commit")
	_, err = os.Stat(huskyPreCommitPath)
	if err == nil { // File exists
		fmt.Printf("Husky pre-commit hook found at %s. Appending lint-staged command...\n", huskyPreCommitPath)

		// Read existing content
		existingContent, readErr := os.ReadFile(huskyPreCommitPath)
		if readErr != nil {
			fmt.Printf("Error reading %s: %v\n", huskyPreCommitPath, readErr)
			// Continue setup
		} else {
			// Append the command, ensuring it's on a new line and executable
			// Add a newline if the file doesn't end with one
			contentToAppend := fmt.Sprintf("%s pre-commit\n", runCmdPrefix)
			if len(existingContent) > 0 && existingContent[len(existingContent)-1] != '\n' {
				contentToAppend = "\n" + contentToAppend
			}

			newContent := string(existingContent) + contentToAppend

			// Write back the updated content
			// Use 0755 permissions to ensure the hook is executable
			writeErr := os.WriteFile(huskyPreCommitPath, []byte(newContent), 0755)
			if writeErr != nil {
				fmt.Printf("Error writing to %s: %v\n", huskyPreCommitPath, writeErr)
			} else {
				fmt.Printf("Command '%s pre-commit' appended to %s.\n", runCmdPrefix, huskyPreCommitPath)
			}
		}
	} else if !os.IsNotExist(err) {
		// Some other error occurred checking the file
		fmt.Printf("Error checking for %s: %v\n", huskyPreCommitPath, err)
	} else {
		// File does not exist, Husky is likely not set up via this tool or manually
		fmt.Println("Husky pre-commit hook not found. Skipping integration.")
	}

	fmt.Println("lint-staged setup complete.")
}

func init() {
	rootCmd.AddCommand(lintStagedCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// lintStagedCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// lintStagedCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
