package ui

import (
	"fmt"
	"strings"

	"codecopy/constants"
	"github.com/fatih/color"
)

// DisplayProjectInfo displays the project information, including the detected project type,
// selected files, and token counts for each file.
func DisplayProjectInfo(projectType string, selectedFiles []string, fileTokenCounts map[string]int) {
	color.New(color.FgGreen, color.Bold).Printf("🚀 Detected project type: %s\n", projectType)

	color.New(color.FgBlue).Println("📂 Selected files:")
	for _, file := range selectedFiles {
		tokenCount := fileTokenCounts[file]
		color.New(color.FgGreen).Printf("📊 File: %s | Token Count: %d\n", file, tokenCount)
	}
	fmt.Println()
}

// DisplayTokenWarning prints a warning message when the token count exceeds the limit.
func DisplayTokenWarning(totalTokens int) {
	color.New(color.FgYellow).Printf("⚠️ Warning: The total token count (%d) exceeds the limit of %d tokens.\n", totalTokens, constants.TokenLimit)
	color.New(color.FgYellow).Println("Consider reducing the number of files or their contents.")
}

// DisplayProjectType prints the detected project type with color and formatting.

// DisplayTreeWithTokenCounts displays the project directory tree with token counts for each file.
func DisplayTreeWithTokenCounts(treeWithTokenCounts []string) {
	color.New(color.FgCyan).Println("🌳 Project Directory Tree | Token Count")
	color.New(color.FgCyan).Println(strings.Repeat("-", 24) + " | " + strings.Repeat("-", 12))

	for _, line := range treeWithTokenCounts {
		color.New(color.FgGreen).Println(line)
	}
}

// DisplayTotalTokens displays the total token count.
func DisplayTotalTokens(totalTokens int) {
	color.New(color.FgCyan).Printf("\n📊 Total Tokens: %d\n", totalTokens)
}

// DisplayCopySuccess prints a success message when the code context is copied to the clipboard.
func DisplayCopySuccess() {
	color.New(color.FgGreen).Println("✅ Code context copied to clipboard!")
}

// DisplaySuccess displays a success message.
func DisplaySuccess(message string) {
	color.New(color.FgGreen, color.Bold).Println(message)
}

// DisplayError displays an error message with additional information.
func DisplayError(err error) {
	color.New(color.FgRed, color.Bold).Printf("❌ Error: %v\n", err)
	color.New(color.FgYellow).Println("If the issue persists, please file an issue at https://github.com/yourusername/codecopy/issues")
}

// DisplayHelpInfo displays information about additional options and functionality.
func DisplayHelpInfo() {
	color.New(color.FgYellow).Println("\n💡 For more options and functionality, run:")
	color.New(color.FgGreen, color.Bold).Println("codecopy --help")
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
