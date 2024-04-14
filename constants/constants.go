package constants

const (
	TokenLimit = 10000
)

var (
	GoFiles = []string{
		".go", ".mod", ".sum", ".toml", ".yaml", ".yml", ".json", ".md", ".txt",
	}

	PythonFiles = []string{
		".py", ".pyc", ".pyd", ".pyo", ".pyw", ".pyz", ".pyi", ".ini", ".toml",
		".yaml", ".yml", ".json", ".md", ".txt",
	}

	JavaScriptFiles = []string{
		".js", ".mjs", ".cjs", ".ts", ".tsx", ".jsx", ".es6", ".es", ".json",
		".jsonc", ".json5", ".css", ".scss", ".sass", ".less", ".styl", ".html",
		".htm", ".xhtml", ".vue", ".svelte", ".angular", ".yaml", ".yml", ".toml",
		".ini", ".md", ".txt",
	}

	RustFiles = []string{
		".rs", ".toml", ".lock", ".yaml", ".yml", ".json", ".md", ".txt",
	}

	PHPFiles = []string{
		".php", ".phtml", ".php3", ".php4", ".php5", ".php7", ".phps", ".ini",
		".json", ".xml", ".yaml", ".yml", ".toml", ".md", ".txt",
	}

	JavaFiles = []string{
		".java", ".class", ".jar", ".xml", ".json", ".yaml", ".yml", ".toml",
		".md", ".txt",
	}

	RubyFiles = []string{
		".rb", ".rbw", ".rake", ".gemspec", ".ru", ".erb", ".yml", ".yaml",
		".json", ".toml", ".md", ".txt",
	}

	CSharpFiles = []string{
		".cs", ".csx", ".sln", ".csproj", ".vbproj", ".xml", ".json", ".yaml",
		".yml", ".toml", ".md", ".txt",
	}

	IgnoredDirs = []string{
		"node_modules", ".git", ".vscode", ".idea", "__pycache__", "venv", "vendor",
		"build", "dist", "bin", "obj", "target", "debug", "release", "tmp", "temp",
		"cache", "logs", "log",
	}
)
