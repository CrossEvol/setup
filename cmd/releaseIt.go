/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// releaseItCmd represents the releaseIt command
var releaseItCmd = &cobra.Command{
	Use:   "releaseIt",
	Short: "Set up release-it for automated releases",
	Long: `This command sets up release-it for your project.

It detects your package manager, installs release-it and the conventional changelog plugin,
creates a configuration file (.release-it.json), and adds a 'release' script to package.json.
release-it helps automate version bumping, changelog generation, and publishing.`,
	Run: func(cmd *cobra.Command, args []string) {
		setupReleaseIt()
	},
}

func setupReleaseIt() {
	fmt.Println("Setting up release-it...")

	// Define package managers and their commands
	packageManagers := []struct {
		name       string
		installCmd []string
	}{
		{"pnpm", []string{"pnpm", "add", "-D", "release-it", "@release-it/conventional-changelog"}},
		{"npm", []string{"npm", "install", "-D", "release-it", "@release-it/conventional-changelog"}},
		{"yarn", []string{"yarn", "add", "-D", "release-it", "@release-it/conventional-changelog"}},
		{"bun", []string{"bun", "add", "-D", "release-it", "@release-it/conventional-changelog"}},
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
		return
	}

	// 1. Install release-it packages
	fmt.Printf("Running installation command: %s %v\n", installCmdArgs[0], installCmdArgs[1:])
	installCmd := exec.Command(installCmdArgs[0], installCmdArgs[1:]...)
	installCmd.Stdout = os.Stdout
	installCmd.Stderr = os.Stderr
	err := installCmd.Run()
	if err != nil {
		fmt.Printf("Warning: Error installing release-it packages with %s: %v\n", foundPackageManager, err)
		fmt.Println("Attempting to continue assuming packages are already installed.")
		// Continue setup even if installation fails, maybe they are already installed
	} else {
		fmt.Println("release-it packages installed successfully.")
	}

	// 2. Create .release-it.json file
	const releaseItConfigFile = `.release-it.json`
	const releaseItConfigContent = `{
  "plugins": {
    "@release-it/conventional-changelog": {
      "preset": {
        "name": "conventionalcommits",
        "types": [
          { "type": "feat", "section": "✨ Features | 新功能" },
          { "type": "fix", "section": "🐛 Bug Fixes | Bug 修复" },
          { "type": "chore", "section": "🎫 Chores | 其他更新" },
          { "type": "docs", "section": "📝 Documentation | 文档" },
          { "type": "style", "section": "💄 Styles | 风格" },
          { "type": "refactor", "section": "♻ Code Refactoring | 代码重构" },
          { "type": "perf", "section": "⚡ Performance Improvements | 性能优化" },
          { "type": "test", "section": "✅ Tests | 测试" },
          { "type": "revert", "section": "⏪ Reverts | 回退" },
          { "type": "build", "section": "👷‍ Build System | 构建" },
          { "type": "ci", "section": "🔧 Continuous Integration | CI 配置" },
          { "type": "config", "section": "🔨 CONFIG | 配置" }
        ]
      },
      "infile": "CHANGELOG.md",
      "ignoreRecommendedBump": true,
      "strictSemVer": true
    }
  },
  "git": {
    "commitMessage": "chore: Release v${version}"
  },
  "github": {
    "release": true,
    "draft": false
  }
}
`

	fmt.Printf("Creating %s...\n", releaseItConfigFile)
	err = os.WriteFile(releaseItConfigFile, []byte(releaseItConfigContent), 0644)
	if err != nil {
		fmt.Printf("Error creating %s: %v\n", releaseItConfigFile, err)
		// Continue setup
	} else {
		fmt.Printf("%s created successfully.\n", releaseItConfigFile)
	}

	// 3. Add "release": "release-it" script to package.json
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
			scripts["release"] = "release-it"

			updatedData, err := json.MarshalIndent(pkgJSON, "", "  ")
			if err != nil {
				fmt.Printf("Error marshalling updated package.json: %v\n", err)
				// Cannot update package.json, stop here for this part
			} else {
				err = os.WriteFile(packageJSONPath, updatedData, 0644)
				if err != nil {
					fmt.Printf("Error writing updated package.json: %v\n", err)
				} else {
					fmt.Println("'release' script added/updated in package.json.")
				}
			}
		}
	}

	fmt.Println("release-it setup complete.")
}

func init() {
	rootCmd.AddCommand(releaseItCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// releaseItCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// releaseItCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
