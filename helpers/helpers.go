package helpers

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"codecopy/constants"
	"github.com/manifoldco/promptui"
	"github.com/tiktoken-go/tokenizer"
)

// ReadFileContent reads the content of a file.
func ReadFileContent(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %v", filePath, err)
	}
	return string(content), nil
}

// CountTokens counts the number of tokens in the given content.
func CountTokens(content string) (int, error) {
	enc, err := tokenizer.Get(tokenizer.Cl100kBase)
	if err != nil {
		return 0, fmt.Errorf("failed to get encoding: %v", err)
	}

	ids, _, err := enc.Encode(content)
	if err != nil {
		return 0, fmt.Errorf("failed to encode content: %v", err)
	}

	return len(ids), nil
}

// WriteToFile writes the given content to a file.
func WriteToFile(content, filename string) error {
	err := os.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write to file %s: %v", filename, err)
	}
	return nil
}

// GenerateTree generates a tree representation of the directory structure.
func GenerateTree(rootDir string) (string, error) {
	cmd := exec.Command("tree", "-F", "-I", strings.Join(constants.IgnoredDirs, "|"))
	cmd.Dir = rootDir
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to generate tree: %v", err)
	}
	return string(output), nil
}

// CopyToClipboard copies the given content to the system clipboard.
func CopyToClipboard(content string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("pbcopy")
	case "linux":
		cmd = exec.Command("xclip", "-selection", "clipboard")
	case "windows":
		cmd = exec.Command("clip")
	default:
		return fmt.Errorf("unsupported platform")
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	defer stdin.Close()

	if err := cmd.Start(); err != nil {
		return err
	}

	if _, err := io.WriteString(stdin, content); err != nil {
		return err
	}

	if err := stdin.Close(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}

// Contains checks if a string is present in a slice of strings.
func Contains(slice []string, item string) bool {
	for _, val := range slice {
		if val == item {
			return true
		}
	}
	return false
}

// ContainsFlag checks if a flag is present in the command-line arguments.
func ContainsFlag(args []string, flag string) bool {
	for _, arg := range args {
		if arg == flag {
			return true
		}
	}
	return false
}

// GetSelectedLanguage retrieves the selected language flag from the command-line arguments.
func GetSelectedLanguage(args []string, languageFlags []string) string {
	for _, arg := range args {
		if Contains(languageFlags, arg) {
			return arg
		}
	}
	return ""
}

// DisplayHelp displays the help message for the codecopy command.

// ExtractFilesFromTree extracts the file paths from the generated tree output.
func ExtractFilesFromTree(treeOutput string) []string {
	var files []string

	lines := strings.Split(treeOutput, "\n")
	for _, line := range lines {
		if strings.Contains(line, "├──") || strings.Contains(line, "└──") {
			file := strings.TrimSpace(line[4:])
			if !strings.HasSuffix(file, "/") {
				files = append(files, file)
			}
		}
	}

	return files
}

// SelectFiles prompts the user to select files or directories to include.
func SelectFiles(rootDir string) ([]string, error) {
	var files []string

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && !Contains(constants.IgnoredDirs, filepath.Base(path)) {
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

// BuildTreeWithTokenCounts constructs the project directory tree with token counts for each file.
func BuildTreeWithTokenCounts(rootDir string, selectedFiles []string, fileTokenCounts map[string]int) []string {
	var treeWithTokenCounts []string

	// Create a map to store directory paths and their corresponding files
	dirMap := make(map[string][]string)
	for _, file := range selectedFiles {
		dirPath, fileName := filepath.Split(file)
		dirMap[dirPath] = append(dirMap[dirPath], fileName)
	}

	// Sort the directory paths to ensure consistent ordering
	dirPaths := make([]string, 0, len(dirMap))
	for dirPath := range dirMap {
		dirPaths = append(dirPaths, dirPath)
	}
	sort.Strings(dirPaths)

	// Build the tree with token counts
	totalTokens := 0
	for _, dirPath := range dirPaths {
		files := dirMap[dirPath]
		sort.Strings(files)

		// Add directory path to the tree
		depth := strings.Count(dirPath, "/") - strings.Count(rootDir, "/")
		indent := strings.Repeat("  ", depth)
		if depth == 0 {
			treeWithTokenCounts = append(treeWithTokenCounts, fmt.Sprintf("%s./", indent))
		} else {
			treeWithTokenCounts = append(treeWithTokenCounts, fmt.Sprintf("%s%s/", indent, filepath.Base(dirPath)))
		}

		// Add files and their token counts to the tree
		for _, fileName := range files {
			filePath := filepath.Join(dirPath, fileName)
			tokenCount := fileTokenCounts[filePath]
			totalTokens += tokenCount
			if fileName == filepath.Base(os.Args[0]) {
				continue // Skip the executable file
			}
			line := fmt.Sprintf("%s  └── %s", indent, fileName)
			treeWithTokenCounts = append(treeWithTokenCounts, fmt.Sprintf("%-25s | %d", line, tokenCount))
		}
	}

	// Add the total token count and directory/file count to the tree
	dirCount := len(dirPaths)
	fileCount := len(selectedFiles)
	if _, err := os.Stat(filepath.Join(rootDir, os.Args[0])); err == nil {
		fileCount-- // Exclude the executable file from the count
	}
	treeWithTokenCounts = append(treeWithTokenCounts, fmt.Sprintf("%d directories, %d files   | %d", dirCount, fileCount, totalTokens))

	return treeWithTokenCounts
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

// DetectProjectType detects the project type based on the file extensions in the root directory.
func DetectProjectType(rootDir string) (string, error) {
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

// GetRelevantFiles retrieves the relevant files based on the selected language and project type.
func GetRelevantFiles(rootDir, projectType, selectedLanguage string) ([]string, error) {
	var relevantFiles []string

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			ext := filepath.Ext(path)
			if IsRelevantFile(ext, selectedLanguage, projectType) {
				relevantFiles = append(relevantFiles, path)
			}
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk the directory: %v", err)
	}

	return relevantFiles, nil
}

// IsRelevantFile checks if a file is relevant based on the selected language and project type.
func IsRelevantFile(ext, selectedLanguage, projectType string) bool {
	switch selectedLanguage {
	case "-py":
		return Contains(constants.PythonFiles, ext)
	case "-rs":
		return Contains(constants.RustFiles, ext)
	case "-go":
		return Contains(constants.GoFiles, ext)
	case "-js":
		return Contains(constants.JavaScriptFiles, ext)
	case "-php":
		return Contains(constants.PHPFiles, ext)
	case "-java":
		return Contains(constants.JavaFiles, ext)
	case "-rb":
		return Contains(constants.RubyFiles, ext)
	case "-cs":
		return Contains(constants.CSharpFiles, ext)
	default:
		switch projectType {
		case "Go":
			return Contains(constants.GoFiles, ext)
		case "Python":
			return Contains(constants.PythonFiles, ext)
		case "JavaScript/TypeScript":
			return Contains(constants.JavaScriptFiles, ext)
		case "Rust":
			return Contains(constants.RustFiles, ext)
		case "PHP":
			return Contains(constants.PHPFiles, ext)
		case "Java":
			return Contains(constants.JavaFiles, ext)
		case "Ruby":
			return Contains(constants.RubyFiles, ext)
		case "C#":
			return Contains(constants.CSharpFiles, ext)
		default:
			return false
		}
	}
}
