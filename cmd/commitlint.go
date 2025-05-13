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

// commitlintCmd represents the commitlint command
var commitlintCmd = &cobra.Command{
	Use:   "commitlint",
	Short: "Set up commitlint for conventional commits",
	Long: `This command sets up commitlint for your project.

It detects your package manager, installs commitlint CLI and conventional config,
creates a configuration file (commitlint.config.cjs), adds a 'commitlint' script
to package.json, and integrates with Husky if it's set up by adding a hook to
.husky/commit-msg.`,
	Run: func(cmd *cobra.Command, args []string) {
		setupCommitlint()
	},
}

func setupCommitlint() {
	fmt.Println("Setting up commitlint...")

	// Define package managers and their commands
	packageManagers := []struct {
		name       string
		installCmd []string
		runCmd     string // Command prefix for running scripts (e.g., "pnpm run", "npm run")
	}{
		{"pnpm", []string{"pnpm", "add", "-D", "@commitlint/cli", "@commitlint/config-conventional"}, "pnpm run"},
		{"npm", []string{"npm", "install", "-D", "@commitlint/cli", "@commitlint/config-conventional"}, "npm run"},
		{"yarn", []string{"yarn", "add", "-D", "@commitlint/cli", "@commitlint/config-conventional"}, "yarn run"},
		{"bun", []string{"bun", "add", "-D", "@commitlint/cli", "@commitlint/config-conventional"}, "bun run"},
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

	// 1. Install commitlint packages
	fmt.Printf("Running installation command: %s %v\n", installCmdArgs[0], installCmdArgs[1:])
	installCmd := exec.Command(installCmdArgs[0], installCmdArgs[1:]...)
	installCmd.Stdout = os.Stdout
	installCmd.Stderr = os.Stderr
	err := installCmd.Run()
	if err != nil {
		fmt.Printf("Warning: Error installing commitlint packages with %s: %v\n", foundPackageManager, err)
		fmt.Println("Attempting to continue assuming packages are already installed.")
		// Continue setup even if installation fails, maybe they are already installed
	} else {
		fmt.Println("commitlint packages installed successfully.")
	}

	// 2. Create commitlint.config.cjs file
	const commitlintConfigFile = `commitlint.config.cjs`
	const commitlintConfigContent = `module.exports = {
  extends: ['@commitlint/config-conventional'],
  rules: {
    'type-enum': [
      // type枚举
      2,
      'always',
      [
        'build', // 编译相关的修改，例如发布版本、对项目构建或者依赖的改动
        'feat', // 新功能
        'fix', // 修补bug
        'docs', // 文档修改
        'style', // 代码格式修改, 注意不是 css 修改
        'refactor', // 重构
        'perf', // 优化相关，比如提升性能、体验
        'test', // 测试用例修改
        'revert', // 代码回滚
        'ci', // 持续集成修改
        'config', // 配置修改
        'chore', // 其他改动
      ],
    ],
    'type-empty': [2, 'never'], // never: type不能为空; always: type必须为空
    'type-case': [0, 'always', 'lower-case'], // type必须小写，upper-case大写，camel-case小驼峰，kebab-case短横线，pascal-case大驼峰，等等
    'scope-empty': [0],
    'scope-case': [0],
    'subject-empty': [2, 'never'], // subject不能为空
    'subject-case': [0],
    'subject-full-stop': [0, 'never', '.'], // subject以.为结束标记
    'header-max-length': [2, 'always', 72], // header最长72
    'body-leading-blank': [0], // body换行
    'footer-leading-blank': [0, 'always'], // footer以空行开头
  },
};
`

	fmt.Printf("Creating %s...\n", commitlintConfigFile)
	err = os.WriteFile(commitlintConfigFile, []byte(commitlintConfigContent), 0644)
	if err != nil {
		fmt.Printf("Error creating %s: %v\n", commitlintConfigFile, err)
		// Continue setup
	} else {
		fmt.Printf("%s created successfully.\n", commitlintConfigFile)
	}

	// 3. Add "commitlint" script to package.json
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
			scripts["commitlint"] = "commitlint --config commitlint.config.cjs -e -V"

			updatedData, err := json.MarshalIndent(pkgJSON, "", "  ")
			if err != nil {
				fmt.Printf("Error marshalling updated package.json: %v\n", err)
				// Cannot update package.json, stop here for this part
			} else {
				err = os.WriteFile(packageJSONPath, updatedData, 0644)
				if err != nil {
					fmt.Printf("Error writing updated package.json: %v\n", err)
				} else {
					fmt.Println("'commitlint' script added/updated in package.json.")
				}
			}
		}
	}

	// 4. Integrate with Husky if .husky/commit-msg exists
	huskyCommitMsgPath := filepath.Join(".husky", "commit-msg")
	_, err = os.Stat(huskyCommitMsgPath)
	if err == nil { // File exists
		fmt.Printf("Husky commit-msg hook found at %s. Appending commitlint command...\n", huskyCommitMsgPath)

		// Read existing content
		existingContent, readErr := os.ReadFile(huskyCommitMsgPath)
		if readErr != nil {
			fmt.Printf("Error reading %s: %v\n", huskyCommitMsgPath, readErr)
			// Continue setup
		} else {
			// Append the command, ensuring it's on a new line and executable
			// Add a newline if the file doesn't end with one
			contentToAppend := fmt.Sprintf("%s commitlint\n", runCmdPrefix)
			if len(existingContent) > 0 && existingContent[len(existingContent)-1] != '\n' {
				contentToAppend = "\n" + contentToAppend
			}

			newContent := string(existingContent) + contentToAppend

			// Write back the updated content
			// Use 0755 permissions to ensure the hook is executable
			writeErr := os.WriteFile(huskyCommitMsgPath, []byte(newContent), 0755)
			if writeErr != nil {
				fmt.Printf("Error writing to %s: %v\n", huskyCommitMsgPath, writeErr)
			} else {
				fmt.Printf("Command '%s commitlint' appended to %s.\n", runCmdPrefix, huskyCommitMsgPath)
			}
		}
	} else if !os.IsNotExist(err) {
		// Some other error occurred checking the file
		fmt.Printf("Error checking for %s: %v\n", huskyCommitMsgPath, err)
	} else {
		// File does not exist, Husky is likely not set up via this tool or manually
		fmt.Println("Husky commit-msg hook not found. Skipping integration.")
	}

	fmt.Println("commitlint setup complete.")
}

func init() {
	rootCmd.AddCommand(commitlintCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// commitlintCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// commitlintCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
