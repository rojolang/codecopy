package helpers

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"codecopy/constants"
	"github.com/tiktoken-go/tokenizer"
)

func ReadFileContent(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %v", filePath, err)
	}
	return string(content), nil
}

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

func WriteToFile(content, filename string) error {
	err := os.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write to file %s: %v", filename, err)
	}
	return nil
}

func GenerateTree(rootDir string) (string, error) {
	cmd := exec.Command("tree", "-F", "-I", strings.Join(constants.IgnoredDirs, "|"))
	cmd.Dir = rootDir
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to generate tree: %v", err)
	}
	return string(output), nil
}

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

func Contains(slice []string, item string) bool {
	for _, val := range slice {
		if val == item {
			return true
		}
	}
	return false
}

func ContainsFlag(args []string, flag string) bool {
	for _, arg := range args {
		if arg == flag {
			return true
		}
	}
	return false
}

func GetSelectedLanguage(args []string, languageFlags []string) string {
	for _, arg := range args {
		if Contains(languageFlags, arg) {
			return arg
		}
	}
	return ""
}

func RemoveFromSlice(slice []string, items ...string) []string {
	var result []string
	for _, item := range slice {
		if !Contains(items, item) {
			result = append(result, item)
		}
	}
	return result
}

func DisplayHelp() {
	fmt.Println("\nUsage:")
	fmt.Println("  codecopy [options]")
	fmt.Println("\nOptions:")
	fmt.Println("  -m    Enable manual file selection mode")
	fmt.Println("  -py   Generate code context for Python projects")
	fmt.Println("  -rs   Generate code context for Rust projects")
	fmt.Println("  -go   Generate code context for Go projects")
	fmt.Println("  -js   Generate code context for JavaScript projects")
	fmt.Println("  -php  Generate code context for PHP projects")
	fmt.Println("  -java Generate code context for Java projects")
	fmt.Println("  -rb   Generate code context for Ruby projects")
	fmt.Println("  -cs   Generate code context for C# projects")
	fmt.Println("  --help Display this help message")
}

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
