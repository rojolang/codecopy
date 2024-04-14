package ccopy

import (
	"fmt"
	"github.com/fatih/color"
	"os"
	"path/filepath"
	"strings"

	"codecopy/constants"
	"codecopy/helpers"
	"codecopy/ui"
)

// Run is the main entry point for the codecopy command.
// Run is the main entry point for the codecopy command.
func Run(args []string) error {
	rootDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %v", err)
	}

	projectType, err := helpers.DetectProjectType(rootDir)
	if err != nil {
		return fmt.Errorf("failed to detect project type: %v", err)
	}

	manualMode := helpers.ContainsFlag(args, "-m")
	languageFlags := []string{"-py", "-rs", "-go", "-js", "-php", "-java", "-rb", "-cs"}
	selectedLanguage := helpers.GetSelectedLanguage(args, languageFlags)

	var selectedFiles []string
	if manualMode {
		selectedFiles, err = helpers.SelectFiles(rootDir)
		if err != nil {
			return fmt.Errorf("failed to perform manual file selection: %v", err)
		}
	} else {
		selectedFiles, err = helpers.GetRelevantFiles(rootDir, projectType, selectedLanguage)
		if err != nil {
			return fmt.Errorf("failed to get relevant files: %v", err)
		}
	}

	if len(selectedFiles) == 0 {
		treeOutput, err := helpers.GenerateTree(rootDir)
		if err != nil {
			return fmt.Errorf("failed to generate tree: %v", err)
		}
		selectedFiles = helpers.ExtractFilesFromTree(treeOutput)
	}

	codeContext, totalTokens, fileTokenCounts, err := generateCodeContext(rootDir, selectedFiles)
	if err != nil {
		return fmt.Errorf("failed to generate code context: %v", err)
	}

	if totalTokens > constants.TokenLimit {
		ui.DisplayTokenWarning(totalTokens)
		selectedFiles, err = helpers.SelectFilesToRemove(selectedFiles)
		if err != nil {
			return fmt.Errorf("failed to select files to remove: %v", err)
		}
		codeContext, totalTokens, fileTokenCounts, err = generateCodeContext(rootDir, selectedFiles)
		if err != nil {
			return fmt.Errorf("failed to generate code context: %v", err)
		}
	}

	ui.DisplayProjectInfo(projectType, selectedFiles, fileTokenCounts)

	excludedFiles := getExcludedFiles(rootDir)
	if len(excludedFiles) > 0 {
		color.New(color.FgYellow).Printf("🚫 Excluded files: %s\n\n", strings.Join(excludedFiles, ", "))
	}

	treeWithTokenCounts := helpers.BuildTreeWithTokenCounts(rootDir, selectedFiles, fileTokenCounts)
	ui.DisplayTreeWithTokenCounts(treeWithTokenCounts)
	ui.DisplayTotalTokens(totalTokens)

	if err := helpers.CopyToClipboard(codeContext); err != nil {
		ui.DisplayError(fmt.Errorf("failed to copy code context to clipboard: %v", err))
		if err := helpers.WriteToFile(codeContext, "code_context.txt"); err != nil {
			return fmt.Errorf("failed to write code context to file: %v", err)
		}
		ui.DisplaySuccess("Code context generated and written to code_context.txt")
		return nil
	}

	ui.DisplayCopySuccess()
	ui.DisplaySuccess("Code context generated and copied successfully!")

	if helpers.ContainsFlag(args, "--help") {
		ui.DisplayHelp()
	}

	return nil
}

// generateCodeContext generates the code context and calculates token counts for the selected files.
func generateCodeContext(rootDir string, selectedFiles []string) (string, int, map[string]int, error) {
	var codeContext strings.Builder
	totalTokens := 0
	fileTokenCounts := make(map[string]int)

	// Calculate token counts for selected files
	for _, file := range selectedFiles {
		content, err := helpers.ReadFileContent(file)
		if err != nil {
			fmt.Printf("Warning: failed to read file %s: %v\n", file, err)
			continue
		}

		tokenCount, err := helpers.CountTokens(content)
		if err != nil {
			fmt.Printf("Warning: failed to count tokens for file %s: %v\n", file, err)
			continue
		}

		totalTokens += tokenCount
		fileTokenCounts[file] = tokenCount
	}

	// Build code context
	for _, file := range selectedFiles {
		content, err := helpers.ReadFileContent(file)
		if err != nil {
			fmt.Printf("Warning: failed to read file %s: %v\n", file, err)
			continue
		}

		relPath := strings.TrimPrefix(file, rootDir+"/")
		codeContext.WriteString(fmt.Sprintf("\n%s\n\n", relPath))
		codeContext.WriteString(content)
		codeContext.WriteString("\n")
	}

	var output strings.Builder
	output.WriteString(fmt.Sprintf("Root Directory: %s\n\n", rootDir))
	output.WriteString(fmt.Sprintf("Total Tokens: %d\n\n", totalTokens))
	output.WriteString("Code Context:\n")
	output.WriteString(codeContext.String())

	return output.String(), totalTokens, fileTokenCounts, nil
}

// getExcludedFiles retrieves the excluded files based on the ignored directories.
func getExcludedFiles(rootDir string) []string {
	var excludedFiles []string

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			for _, dir := range constants.IgnoredDirs {
				if strings.Contains(path, dir) {
					excludedFiles = append(excludedFiles, path)
					break
				}
			}
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Warning: failed to get excluded files: %v\n", err)
	}

	return excludedFiles
}
