// ccopy/codecopy.go

package ccopy

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"codecopy/constants"
	"codecopy/helpers"
	"codecopy/ui"
)

func Run(args []string) error {
	rootDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %v", err)
	}

	projectType, err := detectProjectType(rootDir)
	if err != nil {
		return fmt.Errorf("failed to detect project type: %v", err)
	}
	ui.DisplayProjectType(projectType)

	manualMode := helpers.ContainsFlag(args, "-m")
	languageFlags := []string{"-py", "-rs", "-go", "-js", "-php", "-java", "-rb", "-cs"}
	selectedLanguage := helpers.GetSelectedLanguage(args, languageFlags)

	var selectedFiles []string
	if manualMode {
		selectedFiles, err = ui.SelectFiles(rootDir)
		if err != nil {
			return fmt.Errorf("failed to perform manual file selection: %v", err)
		}
	} else {
		selectedFiles, err = getRelevantFiles(rootDir, projectType, selectedLanguage)
		if err != nil {
			return fmt.Errorf("failed to get relevant files: %v", err)
		}
	}

	if len(selectedFiles) == 0 {
		treeOutput, err := helpers.GenerateTree(rootDir)
		if err != nil {
			return fmt.Errorf("failed to generate tree: %v", err)
		}

		treeFiles := helpers.ExtractFilesFromTree(treeOutput)
		selectedFiles = treeFiles
	}

	codeContext, totalTokens, err := generateCodeContext(rootDir, selectedFiles)
	if err != nil {
		return fmt.Errorf("failed to generate code context: %v", err)
	}

	treeOutput, err := helpers.GenerateTree(rootDir)
	if err != nil {
		return fmt.Errorf("failed to generate tree: %v", err)
	}

	if totalTokens > constants.TokenLimit {
		ui.DisplayTokenWarning(totalTokens)

		selectedFiles, err = ui.SelectFilesToRemove(selectedFiles)
		if err != nil {
			return fmt.Errorf("failed to select files to remove: %v", err)
		}

		codeContext, totalTokens, err = generateCodeContext(rootDir, selectedFiles)
		if err != nil {
			return fmt.Errorf("failed to generate code context: %v", err)
		}
	}

	ui.DisplaySelectedFiles(selectedFiles)
	ui.DisplayTreeAndTokens(rootDir, treeOutput, selectedFiles, totalTokens)

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

	return nil
}

func detectProjectType(rootDir string) (string, error) {
	fileTypes := make(map[string]int)

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			ext := filepath.Ext(path)
			fileTypes[ext]++
		}

		return nil
	})

	if err != nil {
		return "", fmt.Errorf("failed to walk the directory: %v", err)
	}

	var projectType string
	var maxCount int

	for ext, count := range fileTypes {
		switch ext {
		case ".go":
			if count > maxCount {
				maxCount = count
				projectType = "Go"
			}
		case ".py":
			if count > maxCount {
				maxCount = count
				projectType = "Python"
			}
		case ".js", ".ts":
			if count > maxCount {
				maxCount = count
				projectType = "JavaScript/TypeScript"
			}
		case ".rs":
			if count > maxCount {
				maxCount = count
				projectType = "Rust"
			}
		case ".php":
			if count > maxCount {
				maxCount = count
				projectType = "PHP"
			}
		}
	}

	if projectType == "" {
		projectType = "Unknown"
	}

	return projectType, nil
}

func getRelevantFiles(rootDir, projectType, selectedLanguage string) ([]string, error) {
	var relevantFiles []string

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			ext := filepath.Ext(path)
			switch selectedLanguage {
			case "-py":
				if helpers.Contains(constants.PythonFiles, ext) {
					relevantFiles = append(relevantFiles, path)
				}
			case "-rs":
				if helpers.Contains(constants.RustFiles, ext) {
					relevantFiles = append(relevantFiles, path)
				}
			case "-go":
				if helpers.Contains(constants.GoFiles, ext) {
					relevantFiles = append(relevantFiles, path)
				}
			case "-js":
				if helpers.Contains(constants.JavaScriptFiles, ext) {
					relevantFiles = append(relevantFiles, path)
				}
			case "-php":
				if helpers.Contains(constants.PHPFiles, ext) {
					relevantFiles = append(relevantFiles, path)
				}
			case "-java":
				if helpers.Contains(constants.JavaFiles, ext) {
					relevantFiles = append(relevantFiles, path)
				}
			case "-rb":
				if helpers.Contains(constants.RubyFiles, ext) {
					relevantFiles = append(relevantFiles, path)
				}
			case "-cs":
				if helpers.Contains(constants.CSharpFiles, ext) {
					relevantFiles = append(relevantFiles, path)
				}
			default:
				switch projectType {
				case "Go":
					if helpers.Contains(constants.GoFiles, ext) {
						relevantFiles = append(relevantFiles, path)
					}
				case "Python":
					if helpers.Contains(constants.PythonFiles, ext) {
						relevantFiles = append(relevantFiles, path)
					}
				case "JavaScript/TypeScript":
					if helpers.Contains(constants.JavaScriptFiles, ext) {
						relevantFiles = append(relevantFiles, path)
					}
				case "Rust":
					if helpers.Contains(constants.RustFiles, ext) {
						relevantFiles = append(relevantFiles, path)
					}
				case "PHP":
					if helpers.Contains(constants.PHPFiles, ext) {
						relevantFiles = append(relevantFiles, path)
					}
				case "Java":
					if helpers.Contains(constants.JavaFiles, ext) {
						relevantFiles = append(relevantFiles, path)
					}
				case "Ruby":
					if helpers.Contains(constants.RubyFiles, ext) {
						relevantFiles = append(relevantFiles, path)
					}
				case "C#":
					if helpers.Contains(constants.CSharpFiles, ext) {
						relevantFiles = append(relevantFiles, path)
					}
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk the directory: %v", err)
	}

	return relevantFiles, nil
}

func generateCodeContext(rootDir string, selectedFiles []string) (string, int, error) {
	var codeContext strings.Builder
	totalTokens := 0

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

		relPath := strings.TrimPrefix(file, rootDir+"/")
		codeContext.WriteString(fmt.Sprintf("\n%s\n\n", relPath))
		codeContext.WriteString(content)
		codeContext.WriteString("\n")
	}

	tree, err := helpers.GenerateTree(rootDir)
	if err != nil {
		return "", 0, fmt.Errorf("failed to generate directory tree: %v", err)
	}

	var output strings.Builder
	output.WriteString(fmt.Sprintf("Root Directory: %s\n\n", rootDir))
	output.WriteString(fmt.Sprintf("Total Tokens: %d\n\n", totalTokens))
	output.WriteString("Tree:\n")
	output.WriteString(tree)
	output.WriteString("\n")
	output.WriteString(codeContext.String())

	return output.String(), totalTokens, nil
}
