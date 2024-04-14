package ui

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"codecopy/constants"
	"codecopy/helpers"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
)

// DisplayProjectType prints the detected project type with color and formatting.
func DisplayProjectType(projectType string) {
	color.New(color.FgGreen, color.Bold).Printf("🚀 Detected project type: %s\n", projectType)
}

// DisplayTokenWarning prints a warning message when the token count exceeds the limit.
func DisplayTokenWarning(totalTokens int) {
	color.New(color.FgYellow).Printf("⚠️ Warning: The total token count (%d) exceeds the limit of %d tokens.\n", totalTokens, constants.TokenLimit)
	color.New(color.FgYellow).Println("Consider reducing the number of files or their contents.")
}

// DisplayCopySuccess prints a success message when the code context is copied to the clipboard.
func DisplayCopySuccess() {
	color.New(color.FgGreen).Println("✅ Code context copied to clipboard!")
}

// DisplaySelectedFiles prints the list of selected files.
func DisplaySelectedFiles(selectedFiles []string) {
	color.New(color.FgBlue).Println("📂 Selected files:")
	for _, file := range selectedFiles {
		color.New(color.FgGreen).Printf("  - %s\n", file)
	}
}

// SelectFiles prompts the user to select files or directories to include.
func SelectFiles(rootDir string) ([]string, error) {
	var files []string

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && !helpers.Contains(constants.IgnoredDirs, filepath.Base(path)) {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk the directory: %v", err)
	}

	selectedFiles, err := multiSelectPrompt("Select files/directories to include", files)
	if err != nil {
		return nil, fmt.Errorf("failed to perform file selection: %v", err)
	}

	return selectedFiles, nil
}

// SelectFilesToRemove prompts the user to select files or directories to remove.
func SelectFilesToRemove(selectedFiles []string) ([]string, error) {
	removedFiles, err := multiSelectPrompt("Select files/directories to remove", selectedFiles)
	if err != nil {
		return nil, fmt.Errorf("failed to perform file selection: %v", err)
	}

	return removedFiles, nil
}

// DisplayTreeAndTokens displays the project directory tree with token counts for each file.
// CalculateFileTokenCounts calculates the token counts for each file and returns a map of file paths to token counts.
func CalculateFileTokenCounts(selectedFiles []string) (map[string]int, error) {
	fileTokenCounts := make(map[string]int)

	for _, file := range selectedFiles {
		content, err := helpers.ReadFileContent(file)
		if err != nil {
			color.New(color.FgYellow).Printf("⚠️ Warning: failed to read file %s: %v\n", file, err)
			continue
		}

		tokenCount, err := helpers.CountTokens(content)
		if err != nil {
			color.New(color.FgYellow).Printf("⚠️ Warning: failed to count tokens for file %s: %v\n", file, err)
			continue
		}

		fileTokenCounts[file] = tokenCount
	}

	return fileTokenCounts, nil
}

// DisplayTreeWithTokenCounts displays the project directory tree with token counts for each file.
func DisplayTreeWithTokenCounts(rootDir, treeOutput string, fileTokenCounts map[string]int) {
	color.New(color.FgCyan).Println("🌳 Project Directory Tree | Token Count")
	color.New(color.FgCyan).Println(strings.Repeat("-", 24) + " | " + strings.Repeat("-", 12))

	treeLines := strings.Split(treeOutput, "\n")
	for _, line := range treeLines {
		if line == "" {
			continue
		}

		if strings.HasSuffix(line, "/") {
			fmt.Printf("%-23s |\n", line)
		} else {
			// Extract the file path from the tree output
			filePath := strings.TrimSpace(strings.TrimSuffix(line, "*"))
			filePath = filepath.Join(rootDir, filePath)

			// Get the token count for the file path
			tokenCount := fileTokenCounts[filePath]
			color.New(color.FgGreen).Printf("%-23s | %d\n", line, tokenCount)
		}
	}
}

// DisplayTreeAndTokens displays the project directory tree with token counts for each file.
func DisplayTreeAndTokens(rootDir, treeOutput string, selectedFiles []string, totalTokens int) error {
	fileTokenCounts, err := CalculateFileTokenCounts(selectedFiles)
	if err != nil {
		return fmt.Errorf("failed to calculate file token counts: %v", err)
	}

	DisplayTreeWithTokenCounts(rootDir, treeOutput, fileTokenCounts)
	color.New(color.FgCyan).Printf("\n📊 Total Tokens: %d\n", totalTokens)

	return nil
}

// DisplayHelpInfo displays information about additional options and functionality.
func DisplayHelpInfo() {
	color.New(color.FgYellow).Println("\n💡 For more options and functionality, run:")
	color.New(color.FgGreen, color.Bold).Println("codecopy --help")
}

// ConfirmAction prompts the user for a yes/no confirmation.
func ConfirmAction(message string) (bool, error) {
	prompt := promptui.Prompt{
		Label:     message,
		IsConfirm: true,
	}

	result, err := prompt.Run()
	if err != nil {
		return false, fmt.Errorf("failed to get user confirmation: %v", err)
	}

	return strings.ToLower(result) == "y", nil
}

// DisplayError displays an error message with additional information.
func DisplayError(err error) {
	color.New(color.FgRed, color.Bold).Printf("❌ Error: %v\n", err)
	color.New(color.FgYellow).Println("If the issue persists, please file an issue at https://github.com/yourusername/codecopy/issues")
}

// DisplaySuccess displays a success message and additional information if the --help flag is provided.
func DisplaySuccess(message string) {
	color.New(color.FgGreen, color.Bold).Println(message)
	if helpers.ContainsFlag(os.Args[1:], "--help") {
		helpers.DisplayHelp()
	}
}

func multiSelectPrompt(label string, items []string) ([]string, error) {
	searcher := func(input string, index int) bool {
		return strings.Contains(items[index], input)
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "\U0001F449 {{ . | cyan }}",
		Inactive: "  {{ . | cyan }}",
		Selected: "\U00002705 {{ . | green }}",
	}

	prompt := promptui.Select{
		Label:     label,
		Items:     items,
		Templates: templates,
		Size:      10,
		HideHelp:  true,
		IsVimMode: true,
		Searcher:  searcher,
	}

	var selectedItems []string

	for {
		index, _, err := prompt.Run()
		if err != nil {
			if errors.Is(err, promptui.ErrInterrupt) {
				break
			}
			return nil, err
		}

		selectedItem := items[index]
		selectedItems = append(selectedItems, selectedItem)

		items = append(items[:index], items[index+1:]...)

		if len(items) == 0 {
			break
		}

		prompt.Items = items
		prompt.Label = "Select another file/directory or press Enter to continue"
	}

	return selectedItems, nil
}

// DisplayHelp displays the help information for the codecopy command.
func DisplayHelp() {
	color.New(color.FgGreen, color.Bold).Println("codecopy - Copy code context to clipboard")
	color.New(color.FgYellow).Println("\nUsage:")
	color.New(color.FgCyan).Println("  codecopy [options]")
	color.New(color.FgYellow).Println("\nOptions:")
	color.New(color.FgCyan).Println("  -m    Enable manual file selection mode")
	color.New(color.FgCyan).Println("  -py   Generate code context for Python projects")
	color.New(color.FgCyan).Println("  -rs   Generate code context for Rust projects")
	color.New(color.FgCyan).Println("  -go   Generate code context for Go projects")
	color.New(color.FgCyan).Println("  -js   Generate code context for JavaScript projects")
	color.New(color.FgCyan).Println("  -php  Generate code context for PHP projects")
	color.New(color.FgCyan).Println("  -java Generate code context for Java projects")
	color.New(color.FgCyan).Println("  -rb   Generate code context for Ruby projects")
	color.New(color.FgCyan).Println("  -cs   Generate code context for C# projects")
	color.New(color.FgCyan).Println("  --help Display this help message")
}
