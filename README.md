# Lorren

![Go Version](https://img.shields.io/badge/go-1.27%2B-00ADD8?logo=go&logoColor=white)
![License](https://img.shields.io/github/license/idoceb00/lorren)
![Status](https://img.shields.io/badge/status-in%20development-yellow)

A CLI wizard that interviews you about your daily habits and training sessions, then writes the results as structured markdown files into your Obsidian vault.

## Features

- **`lorren day`** — interactive wizard for daily habits, meals, and a short day evaluation. Writes a markdown file per day, with Dataview-ready frontmatter.
- **First-run setup** — asks for your Obsidian vault path once and remembers it. Validates the folder exists, never creates it.

## Installation

Requires Go 1.27+.

```bash
git clone https://github.com/idoceb00/lorren.git
cd lorren
go build -o lorren ./cmd/lorren
```

## Usage

```bash
./lorren day
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

## Status

`lorren day` works end to end. `lorren train` (training session logging) is planned, not started yet.

## Built with

- [Go](https://go.dev/) 1.27+
- [Cobra](https://github.com/spf13/cobra) — CLI command routing
- [Viper](https://github.com/spf13/viper) — configuration
- [huh](https://github.com/charmbracelet/huh) — interactive terminal forms

## License

MIT