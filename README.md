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

## Installation

### Option 1: Using `go install` (Easiest)

If you have Go installed, you can install directly:

```bash
go install github.com/swantron/readme-lint@latest
```

This will install the binary to `$GOPATH/bin` (or `$HOME/go/bin` by default). Make sure that directory is in your `PATH`.

### Option 2: Download Pre-built Binary

Download the latest release binary for your platform from the [Releases](https://github.com/swantron/readme-lint/releases) page:

**Linux:**
```bash
curl -L -o readme-lint https://github.com/swantron/readme-lint/releases/latest/download/readme-lint-linux-amd64
chmod +x readme-lint
# Optional: move to PATH for global access
sudo mv readme-lint /usr/local/bin/
```

**macOS (Intel):**
```bash
curl -L -o readme-lint https://github.com/swantron/readme-lint/releases/latest/download/readme-lint-darwin-amd64
chmod +x readme-lint
# Remove macOS quarantine attribute
xattr -c readme-lint
# Optional: move to PATH for global access
sudo mv readme-lint /usr/local/bin/
```

**macOS (Apple Silicon):**
```bash
curl -L -o readme-lint https://github.com/swantron/readme-lint/releases/latest/download/readme-lint-darwin-arm64
chmod +x readme-lint
# Remove macOS quarantine attribute
xattr -c readme-lint
# Optional: move to PATH for global access
sudo mv readme-lint /usr/local/bin/
```

**Windows:**
```powershell
curl -L -o readme-lint.exe https://github.com/swantron/readme-lint/releases/latest/download/readme-lint-windows-amd64.exe
```

### Option 3: Build from Source

If you want to build from source:

```bash
# Clone the repository
git clone https://github.com/swantron/readme-lint.git
cd readme-lint

# Build the binary
go build -o readme-lint .

# The binary is now ready to use
./readme-lint --help
```

## Usage

### Running Locally

```bash
# Lint the default README.md in the current directory
./readme-lint

# Specify a path to a different file
./readme-lint ./docs/README.md

# Get help
./readme-lint --help
```

## Using in CI/CD

### Option 1: Download Pre-built Binary (Recommended)

This is the fastest and most efficient approach for CI:

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

      - name: Download readme-lint
        run: |
          curl -L -o readme-lint https://github.com/swantron/readme-lint/releases/latest/download/readme-lint-linux-amd64
          chmod +x readme-lint

      - name: Run README Linter
        run: ./readme-lint ./README.md
```

### Option 2: Install via `go install` in CI

This is simpler than building from source and doesn't require downloading binaries:

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

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.21'

      - name: Install readme-lint
        run: go install github.com/swantron/readme-lint@latest

      - name: Run README Linter
        run: ~/go/bin/readme-lint ./README.md
```

### Option 3: Build from Source in CI

If you need a specific commit or branch:

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

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.21'

      - name: Build readme-lint
        run: |
          git clone https://github.com/swantron/readme-lint.git /tmp/readme-lint
          cd /tmp/readme-lint
          go build -o readme-lint .

      - name: Run README Linter
        run: /tmp/readme-lint/readme-lint ./README.md
```

## Buildkite Pipeline

This repository includes a Buildkite pipeline (`.buildkite/pipeline.yml`) that mirrors the GitHub Actions workflows while showcasing Buildkite-specific capabilities.

### Architecture

```
Push / PR
  └── Buildkite
        └── Self-hosted agent (GCP e2-micro, us-west1)
              ├── fmt check
              ├── vet
              ├── build + self-lint  (depends on fmt + vet)
              └── annotate           (surfaces result in Buildkite UI)

Tag push (v*)
  └── Buildkite
        ├── linux/amd64   ┐
        ├── linux/arm64   │ parallel — all four dispatch simultaneously
        ├── darwin/amd64  │ to the agent pool
        ├── darwin/arm64  ┘
        └── windows/amd64
              └── publish (downloads artifacts, creates GitHub release)
```

### Why Buildkite instead of GitHub Actions for this

**Parallel steps as first-class citizens.** The cross-platform release builds run as five independent steps dispatched in parallel. In GitHub Actions, these are sequential `run` statements in a single job (or require a separate matrix job with YAML overhead). Buildkite's step model maps more naturally to distributing work across an agent pool.

**Annotations.** The CI stage writes a structured annotation directly into the Buildkite build UI — a pass/fail summary visible without opening individual job logs. At scale (many pipelines, many engineers), annotations are the difference between a usable build dashboard and one that requires drilling into every log.

**Agent targeting.** All steps specify `queue: gcp`, routing jobs to the self-hosted agent on GCP. In a real platform engineering context this pattern extends naturally: GPU workloads route to `queue: gpu`, macOS builds to `queue: macos`, etc. — all managed via the same pipeline syntax.

### Pipeline setup

The pipeline connects to this repo via Buildkite → New Pipeline → point at `github.com/swantron/readme-lint`. Buildkite reads `.buildkite/pipeline.yml` automatically.

The agent is provisioned separately via [buildkite-gcp-agent](https://github.com/swantron/buildkite-gcp-agent) — a Terraform config that provisions a free-tier GCP instance and registers it with Buildkite. Infrastructure changes are applied automatically via GitHub Actions on merge to main.

### Using in Buildkite pipelines

```yaml
# .buildkite/pipeline.yml
steps:
  - label: "Lint README"
    command: |
      curl -fsSL -o readme-lint https://github.com/swantron/readme-lint/releases/latest/download/readme-lint-linux-amd64
      chmod +x readme-lint
      ./readme-lint ./README.md
    agents:
      queue: gcp
```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

