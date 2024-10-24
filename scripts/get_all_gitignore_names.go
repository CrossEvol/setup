package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func (r *Result) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Result struct {
	SHA       string `json:"sha"`
	URL       string `json:"url"`
	Tree      []Tree `json:"tree"`
	Truncated bool   `json:"truncated"`
}

type Tree struct {
	Path string `json:"path"`
	Mode string `json:"mode"`
	Type string `json:"type"`
	SHA  string `json:"sha"`
	URL  string `json:"url"`
}

func main() {
	// Fetch data from the API
	resp, err := http.Get("https://api.github.com/repos/github/gitignore/git/trees/main?recursive=1")
	if err != nil {
		fmt.Println("Error fetching data:", err)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// Unmarshal the JSON data
	var result Result
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

	// Collect paths with .gitignore suffix
	gitignoreMap := make(map[string]string)
	for _, tree := range result.Tree {
		if strings.HasSuffix(tree.Path, ".gitignore") {
			key := strings.TrimSuffix(tree.Path, ".gitignore")
			gitignoreMap[key] = tree.Path
		}
	}

	// Marshal the map to JSON
	jsonData, err := json.MarshalIndent(gitignoreMap, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling map to JSON:", err)
		return
	}

	// Write JSON data to a file
	err = os.WriteFile("assets/gitignore_pairs.json", jsonData, 0644)
	if err != nil {
		fmt.Println("Error writing JSON to file:", err)
		return
	}

	fmt.Println("gitignore paths have been written to gitignore_paths.json")
}
