package linter

import (
	"os"
	"testing"
)

func TestCheckTitle(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		hasError bool
	}{
		{
			name:     "Valid title",
			content:  "# My Project\n\nDescription",
			hasError: false,
		},
		{
			name:     "Missing title",
			content:  "## Subtitle\n\nDescription",
			hasError: true,
		},
		{
			name:     "Empty file",
			content:  "",
			hasError: true,
		},
		{
			name:     "Title after blank lines",
			content:  "\n\n# My Project\n\nDescription",
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Linter{
				lines: splitLines(tt.content),
			}
			result := l.checkTitle()
			if tt.hasError && result.Message == "" {
				t.Errorf("Expected error but got none")
			}
			if !tt.hasError && result.Message != "" {
				t.Errorf("Expected no error but got: %s", result.Message)
			}
		})
	}
}

func TestCheckSections(t *testing.T) {
	tests := []struct {
		name          string
		content       string
		missingCount  int
	}{
		{
			name: "All sections present",
			content: `# Project
## Installation
## Usage
## License`,
			missingCount: 0,
		},
		{
			name: "Missing one section",
			content: `# Project
## Installation
## Usage`,
			missingCount: 1,
		},
		{
			name:         "Missing all sections",
			content:      "# Project\n\nDescription",
			missingCount: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Linter{
				lines: splitLines(tt.content),
			}
			results := l.checkSections()
			if len(results) != tt.missingCount {
				t.Errorf("Expected %d missing sections, got %d", tt.missingCount, len(results))
			}
		})
	}
}

func TestCheckPlaceholders(t *testing.T) {
	tests := []struct {
		name          string
		content       string
		expectedCount int
	}{
		{
			name:          "No placeholders",
			content:       "# Project\n\nThis is complete.",
			expectedCount: 0,
		},
		{
			name:          "Has TODO",
			content:       "# Project\n\nTODO: Add docs",
			expectedCount: 1,
		},
		{
			name:          "Has coming soon",
			content:       "# Project\n\nComing soon",
			expectedCount: 1,
		},
		{
			name:          "Multiple placeholders",
			content:       "# Project\n\nTODO: Fix\n\nComing soon",
			expectedCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Linter{
				lines: splitLines(tt.content),
			}
			results := l.checkPlaceholders()
			if len(results) != tt.expectedCount {
				t.Errorf("Expected %d placeholders, got %d", tt.expectedCount, len(results))
			}
		})
	}
}

func TestCheckLicenseFile(t *testing.T) {
	tmpDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	
	os.Chdir(tmpDir)

	tests := []struct {
		name         string
		content      string
		createLicense bool
		hasError     bool
	}{
		{
			name:         "No license section",
			content:      "# Project",
			createLicense: false,
			hasError:     false,
		},
		{
			name:         "Has license section and file",
			content:      "# Project\n\n## License\n\nMIT",
			createLicense: true,
			hasError:     false,
		},
		{
			name:         "Has license section but no file",
			content:      "# Project\n\n## License\n\nMIT",
			createLicense: false,
			hasError:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Remove("LICENSE")
			
			if tt.createLicense {
				os.WriteFile("LICENSE", []byte("MIT License"), 0644)
			}

			l := &Linter{
				lines: splitLines(tt.content),
			}
			result := l.checkLicenseFile()
			
			if tt.hasError && result.Message == "" {
				t.Errorf("Expected error but got none")
			}
			if !tt.hasError && result.Message != "" {
				t.Errorf("Expected no error but got: %s", result.Message)
			}
		})
	}
}

func splitLines(content string) []string {
	if content == "" {
		return []string{}
	}
	lines := []string{}
	current := ""
	for _, ch := range content {
		if ch == '\n' {
			lines = append(lines, current)
			current = ""
		} else {
			current += string(ch)
		}
	}
	if current != "" || len(content) > 0 {
		lines = append(lines, current)
	}
	return lines
}
