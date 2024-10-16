package common

import (
	"fmt"
	"os"
	"strings"
)

func UpdatePackageJSON(newEntries string) {
	// Read package.json
	data, err := os.ReadFile("package.json")
	if err != nil {
		fmt.Println("Error reading package.json:", err)
		return
	}

	// Convert data to string
	content := string(data)

	// Find the "scripts" section
	scriptsStart := strings.Index(content, `"scripts"`)
	if scriptsStart == -1 {
		fmt.Println("Error: 'scripts' section not found")
		return
	}

	// Find the opening brace of the scripts object
	braceStart := strings.Index(content[scriptsStart:], "{")
	if braceStart == -1 {
		fmt.Println("Error: Opening brace for 'scripts' not found")
		return
	}
	braceStart += scriptsStart

	//// Insert new entries
	//newEntries := `
	//"a": "1",
	//"b": 2,`
	updatedContent := content[:braceStart+1] + newEntries + content[braceStart+1:]

	// Write the updated content to package2.json
	err = os.WriteFile("package.json", []byte(updatedContent), 0644)
	if err != nil {
		fmt.Println("Error writing package.json:", err)
		return
	}

	fmt.Println("Successfully updated and wrote package2.json")
}
