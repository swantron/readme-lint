package linter

import (
	"os"
	"strings"
)

type LintResult struct {
	Line    int
	Message string
}

type Linter struct {
	lines []string
}

func NewLinter() *Linter {
	return &Linter{}
}

func (l *Linter) Run(filePath string) ([]LintResult, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []LintResult{{Line: 0, Message: "File not found"}}, nil
		}
		return nil, err
	}

	l.lines = strings.Split(string(content), "\n")

	var results []LintResult

	if result := l.checkTitle(); result.Message != "" {
		results = append(results, result)
	}

	results = append(results, l.checkSections()...)
	results = append(results, l.checkPlaceholders()...)

	if result := l.checkLicenseFile(); result.Message != "" {
		results = append(results, result)
	}

	return results, nil
}

func (l *Linter) checkTitle() LintResult {
	for i, line := range l.lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if strings.HasPrefix(trimmed, "# ") {
			return LintResult{}
		}
		return LintResult{Line: i + 1, Message: "No H1 title found (e.g., # Project Name)"}
	}
	return LintResult{Line: 1, Message: "No H1 title found (e.g., # Project Name)"}
}

func (l *Linter) checkSections() []LintResult {
	requiredSections := []string{"## Usage", "## Installation", "## License"}
	var results []LintResult

	for _, section := range requiredSections {
		found := false
		for _, line := range l.lines {
			if strings.Contains(strings.ToLower(line), strings.ToLower(section)) {
				found = true
				break
			}
		}
		if !found {
			results = append(results, LintResult{
				Line:    0,
				Message: "Missing required section: " + section,
			})
		}
	}

	return results
}

func (l *Linter) checkPlaceholders() []LintResult {
	placeholders := []string{"TODO", "coming soon"}
	var results []LintResult

	for i, line := range l.lines {
		lineLower := strings.ToLower(line)
		for _, placeholder := range placeholders {
			if strings.Contains(lineLower, strings.ToLower(placeholder)) {
				results = append(results, LintResult{
					Line:    i + 1,
					Message: "Found placeholder text: '" + placeholder + "'",
				})
				break
			}
		}
	}

	return results
}

func (l *Linter) checkLicenseFile() LintResult {
	hasLicenseSection := false
	for _, line := range l.lines {
		if strings.Contains(strings.ToLower(line), "## license") {
			hasLicenseSection = true
			break
		}
	}

	if !hasLicenseSection {
		return LintResult{}
	}

	if _, err := os.Stat("LICENSE"); os.IsNotExist(err) {
		return LintResult{
			Line:    0,
			Message: "No LICENSE file found in repository root",
		}
	}

	return LintResult{}
}
