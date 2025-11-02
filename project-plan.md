Readme Linter: Project Todo List

This document breaks down the development of the readme-lint tool into concrete, actionable steps.

Phase 1: Project Setup & CLI Structure

The first phase is setting up your Go module and building the "skeleton" of the command-line application. We will use spf13/cobra, a powerful library for building modern Go CLIs.

[ ] Initialize the Go module: go mod init github.com/your-username/readme-lint

[ ] Add cobra as a dependency: go get -u github.com/spf13/cobra@latest

[ ] Create the cobra-cli recommended structure:

main.go (the entry point)

cmd/root.go (defines the main readme-lint command)

[ ] In cmd/root.go, add logic to accept one optional argument: the path to the README file.

[ ] If no argument is provided, set a default value of ./README.md.

[ ] Implement the Run function in cmd/root.go to print the file path it's supposed to lint (e.g., Linting file: ./README.md).

Phase 2: Core Linter Engine

This phase involves creating the business logic for reading the file and running checks.

[ ] Create a new package: pkg/linter.

[ ] In the linter package, create a Linter struct.

[S] In pkg/linter/linter.go, create a LintResult struct to hold errors (e.g., type LintResult struct { Line int; Message string }).

[ ] Create a NewLinter() function that initializes the linter.

[ ] Create a public method: func (l *Linter) Run(filePath string) ([]LintResult, error).

[ ] Inside Run, use os.ReadFile(filePath) to read the file's content.

[ ] Handle errors gracefully (e.g., return a "File not found" LintResult).

[ ] Store the file content as an array of strings (one for each line) to make line-based checks easier.

[ ] In cmd/root.go, call linter.Run() and print any results.

Phase 3: Implementing Linting Rules

This is the core of the project. We'll add private methods to the Linter struct for each rule.

[ ] Rule 1: Title Check

Create func (l *Linter) checkTitle() LintResult.

Check if the first non-empty line of the file starts with # .

If not, return a LintResult{Line: 1, Message: "No H1 title found. (e.g., # Project Name)"}.

[ ] Rule 2: Section Check

Create func (l *Linter) checkSections() []LintResult.

Define a slice of required sections (e.g., []string{"## Usage", "## Installation", "## License"}).

Loop through the file content and check if each of these sections exists (case-insensitive string matching is fine).

For each missing section, add a LintResult (e.g., Message: "Missing required section: ## Usage").

[ ] Rule 3: Placeholder Check

Create func (l *Linter) checkPlaceholders() []LintResult.

Define a slice of "bad" words (e.g., "TODO", "coming soon").

Loop through every line and check for these words (case-insensitive).

For each one found, return a LintResult{Line: X, Message: "Found placeholder text: 'TODO'"}.

[ ] Rule 4: License File Check

Create func (l *Linter) checkLicenseFile() LintResult.

Use os.Stat("LICENSE") to check if a LICENSE file exists in the current directory.

If it doesn't (and os.IsNotExist(err) is true), return a LintResult{Message: "No LICENSE file found in repository root."}.

[ ] Integrate Rules: Update the Linter.Run() method to call all the new check...() methods and aggregate their results into a single []LintResult slice.

Phase 4: Reporting & Exit Codes

This makes the tool useful for humans and for CI.

[ ] In main.go / cmd/root.go, after getting the results from the linter, loop through them and print them to the console in a clean format (e.g., [FAIL] README.md:1: No H1 title found.).

[ ] After printing, check if the slice of results has a length greater than 0.

[ ] If len(results) > 0, exit with a non-zero status code: os.Exit(1). This is critical for failing the GitHub Action.

[ ] If len(results) == 0, print a "Success!" message and os.Exit(0).

Phase 5: "V2" Advanced Features (Future Ideas)

[ ] Broken Link Checker:

Use regex to find all markdown links: \[[^\]]+\]\((http[^\)]+)\).

For each URL found, make an HTTP HEAD request using the net/http package.

To make it fast, run these checks concurrently using goroutines and wait for them all to finish using a sync.WaitGroup.

Report any links that return a 404 or other error status code.

[ ] Configuration File:

Add support for a .readme-lint.yml file.

Use a YAML parsing library (like gopkg.in/yaml.v3) to read it.

Allow users to customize the list of required sections.
