# Go README Linter (readme-lint)

readme-lint is a fast, standalone command-line tool written in Go, designed to be run in CI/CD pipelines (like GitHub Actions) to enforce quality and completeness standards for your README.md files.

It is built as a single, static binary with no external dependencies (like the gh CLI or Node.js), making it lightweight and extremely fast to execute.

## Purpose

The main goal of this tool is to programmatically check for common "hygiene" issues in a repository's README. It ensures that every project provides a basic, consistent level of documentation for new visitors and contributors.

When run, it will scan a README.md file, report any errors, and exit with a non-zero status code if any checks fail. This non-zero exit code is what allows it to fail a GitHub Actions step.

## Features (v1.0)

- **File Existence**: Checks that a README.md file (or specified file) exists.
- **Title Check**: Verifies that the file starts with a main title (e.g., `# Project Name`).
- **Section Checks**: Ensures that key sections are present, such as:
  - `## Usage`
  - `## Installation`
  - `## License`
- **Placeholder Check**: Scans for incomplete placeholder text that should be finished before publication.
- **License File Check**: If a `## License` section is found, it will also check that a LICENSE file exists in the repository's root.

## Command-Line Usage

```bash
# Lint the default README.md in the current directory
./readme-lint

# Specify a path to a different file
./readme-lint ./docs/README.md

# Get help
./readme-lint --help
```

## Example GitHub Actions Workflow

This is how you would use the compiled readme-lint binary in a real workflow.

```yaml
# .github/workflows/lint.yml
name: Lint Documentation

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  readme-lint:
    name: Check README
    runs-on: ubuntu-latest
    steps:
      - name: Check out code
        uses: actions/checkout@v4

      # Set up Go to build the tool
      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.2x' # Or your project's version

      # Build the binary
      - name: Build linter
        run: go build -o readme-lint .

      # Run the linter
      # This step will fail if the linter finds any issues
      - name: Run README Linter
        run: ./readme-lint ./README.md
```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
