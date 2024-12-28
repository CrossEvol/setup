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

const AUTO_TYPE = "AutoType"

// pyCmd represents the py command
var pyCmd = &cobra.Command{
	Use:   "py",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		var tools []string

		fmt.Println("Python project initializing....")
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewMultiSelect[string]().Title("Setup Scripts").
					Options(
						huh.NewOption(AUTO_TYPE, AUTO_TYPE),
					).
					Description("Choose your script").
					Value(&tools),
			),
		)

		err := form.Run()

		if err != nil {
			fmt.Println("Uh oh:", err)
			os.Exit(1)
		}

		fmt.Printf("Choose scripts: %s\n", tools)

		if slices.Contains(tools, AUTO_TYPE) {
			generateAutoTypeScript()
		}
	},
}

func generateAutoTypeScript() {
	const pythonFilename = `auto_type.py`
	const pythonFileContent = `
import argparse
import subprocess

def run_autotyping(filename):
  """
  Runs the autotyping command with the given filename.

  Args:
    filename: The path to the Python file to analyze.

  Returns:
    The output of the autotyping command.
  """

  command = [
      "python", "-m", "autotyping", 
      filename, 
      "--none-return", 
      "--scalar-return", 
      "--bool-param", 
      "--int-param", 
      "--float-param", 
      "--str-param", 
      "--bytes-param", 
      "--annotate-optional", 
      "foo:bar.Baz", 
      "--annotate-named-param", 
      "foo:bar.Baz", 
      "--guess-common-names", 
      "--annotate-magics", 
      "--annotate-imprecise-magics"
  ]

  try:
    result = subprocess.run(command, capture_output=True, text=True, check=True)
    return result.stdout
  except subprocess.CalledProcessError as e:
    print(f"Error running autotyping command: {e}")
    return e.stderr

if __name__ == "__main__":
  # Parse arguments
  parser = argparse.ArgumentParser(description="Run autotyping on a Python file")
  parser.add_argument("filename", help="Path to the Python file to analyze")
  args = parser.parse_args()

  # Run the autotyping command
  output = run_autotyping(args.filename)

  # Print the output
  print(output)
`

	err := os.WriteFile(pythonFilename, []byte(pythonFileContent), 0644)
	if err != nil {
		fmt.Printf("Error creating %s: %v\n", pythonFilename, err)
	} else {
		fmt.Printf("%s created successfully\n", pythonFilename)
		fmt.Println("You can read the document on https://github.com/JelleZijlstra/autotyping")
	}
}

func init() {
	rootCmd.AddCommand(pyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
