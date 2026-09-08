# 📝 Lorren

![Go Version](https://img.shields.io/badge/go-1.27%2B-00ADD8?logo=go&logoColor=white)
![CI](https://github.com/idoceb00/lorren/actions/workflows/ci.yml/badge.svg)
![License](https://img.shields.io/github/license/idoceb00/lorren)
![Status](https://img.shields.io/badge/status-in%20development-yellow)

A CLI wizard that interviews you about your daily habits and training sessions, then writes the results as structured markdown files with YAML frontmatter, compatible with Obsidian and its Dataview plugin.

## Features

- **`lorren day`** — interactive wizard for daily habits, meals, and a short day evaluation. Writes a markdown file per day, with Dataview-ready frontmatter.
- **First-run setup** — asks for your daily notes directory once and remembers it. Validates the folder exists, never creates it.

## Installation

Requires Go 1.27+.

```bash
git clone https://github.com/idoceb00/lorren.git
cd lorren
```

**To try it locally**, build a binary in the project directory:

```bash
make build
./bin/lorren day
```

**To use it day to day**, install it onto your `PATH` instead:

```bash
make install
```

This runs `go install ./cmd/lorren`, which compiles the binary and places it in `$GOBIN` (or `$GOPATH/bin` if `GOBIN` isn't set — typically `~/go/bin`). Make sure that directory is on your `PATH`, then run it from anywhere:

```bash
lorren day
```

### Updating

Pull the latest changes and reinstall — this overwrites the existing binary in place:

```bash
git checkout main && git pull && make install
```

If you're going to contribute, also install the git hooks (requires [lefthook](https://github.com/evilmartians/lefthook)):

```bash
lefthook install
```

This runs `gofmt` and `go vet` before each commit, and `go test ./...` before each push.

## Usage

```bash
lorren day
```

First run asks for your vault path and saves it to `~/.config/lorren/config.yaml`. Every run after that goes straight to the habit wizard, writing (or overwriting) `YYYY-MM-DD.md` in your vault:

```markdown
---
date: 2026-08-29
training: true
reading: true
sleep_hours: 7.50
---

## 🍽️ Meals
**Breakfast:** oatmeal and coffee
...
```

## How it works

Lorren follows a hexagonal (ports & adapters) architecture: a plain Go domain at the center that knows nothing about Cobra, huh, or the filesystem — just interfaces describing what it needs. Everything else plugs into those interfaces, so the storage backend or the wizard library can be swapped without touching business logic.

## Development

```bash
make test   # run tests
make lint   # run golangci-lint
make fmt    # list files with formatting issues
```

## Status

`lorren day` works end to end, with unit tests covering `internal/domain` and `internal/repository`. `lorren train` (training session logging) is planned, not started yet.

## Built with

- [Go](https://go.dev/) 1.27+
- [Cobra](https://github.com/spf13/cobra) — CLI command routing
- [Viper](https://github.com/spf13/viper) — configuration
- [huh](https://github.com/charmbracelet/huh) — interactive terminal forms
- [GitHub Actions](https://github.com/features/actions) — CI (build, test, lint on every push to main)

## License

MIT