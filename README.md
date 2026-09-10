# 📝 Lorren

![Go Version](https://img.shields.io/badge/go-1.27%2B-00ADD8?logo=go&logoColor=white)
![CI](https://github.com/idoceb00/lorren/actions/workflows/ci.yml/badge.svg)
![License](https://img.shields.io/github/license/idoceb00/lorren)
![Status](https://img.shields.io/badge/status-in%20development-yellow)

A CLI wizard that interviews you about your daily habits and training sessions, then writes the results as structured markdown files with YAML frontmatter, compatible with Obsidian and its Dataview plugin.

## Features

- **`lorren day`** — interactive wizard for daily habits, meals, and a short day evaluation. Writes one markdown file per day, with Dataview-ready frontmatter.
- **`lorren train`** — interactive wizard for a training session, driven by your own plan template. Asks only what changes: the weight and reps you actually did.
- **`lorren plan`** — prints your active plan as Lorren understands it, so you can check a template after editing it by hand.
- **Plans you write yourself** — a plan is a markdown file in your vault listing the sessions available to log. No weekly schedule, no prescribed frequency: what you actually trained is whatever the logs end up saying.
- **First-run setup** — asks for your vault path once and remembers it. Validates the folder exists, never creates it.

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

## Configuration

First run asks for your vault path and saves it to `~/.config/lorren/config.yaml`, along with the folders it writes to:

```yaml
vault_path: /path/to/your/vault
plans_dir: plans
trainings_dir: trainings
daily_notes_dir: diary
active_plan: hybrid
```

All directories are relative to the vault. Lorren creates the folders it writes into, but never the vault itself and never the plans folder — that one holds files you write.

## Usage

### Daily habits

```bash
lorren day
```

Writes `YYYY-MM-DD.md` in your daily notes folder. Running it again the same day preloads the wizard with what you already logged, so you only change what's new:

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

### Training

A plan is a markdown file in your plans folder. The frontmatter is what Lorren reads; the body is yours, for whatever you want to check mid-workout:

```markdown
---
plan: Strength and boxing
sessions:
  - name: Gym A
    kind: strength
    modality: strength
    exercises:
      - { name: Front squat, block: strength, sets: 3, reps_min: 4, reps_max: 6 }
      - { name: Vertical jump, block: power, sets: 3, reps_min: 3, body_weight: true }
  - name: Boxing
    kind: sport
    modality: boxing
---
```

Sessions come in three kinds, and the kind decides what the wizard asks: `strength` walks your prescribed exercises asking for weight and reps, `cardio` asks what you did, and `sport` just records duration and notes. Descriptive labels like *zone 2* or *submaximal strength* go in `modality`, which is free text and needs no code.

```bash
lorren train
```

Picks a session from your active plan, asks what you did, and writes `YYYY-MM-DD <session>.md` in your trainings folder. Prescribed sets and reps are copied into the log, so a note still makes sense years later even if the plan has changed since.

```bash
lorren plan
```

Prints the active plan. Useful right after editing a template by hand: if a session is invalid, this tells you which one and why.

## How it works

Lorren follows a hexagonal (ports & adapters) architecture: a plain Go domain at the center that knows nothing about Cobra, huh, or the filesystem — just interfaces describing what it needs. Everything else plugs into those interfaces, so the storage backend or the wizard library can be swapped without touching business logic.

Two ideas shape the training side. Plans are data, not code: adding a session, a sport, or a whole training block means editing a markdown file, never touching Go. And what a session *structurally is* (`kind`) is kept apart from what you *call* it (`modality`), so a new training style costs nothing.

## Development

```bash
make test   # run tests
make lint   # run golangci-lint
make fmt    # list files with formatting issues
```

## Status

## Status

## Status

| Command | State | Tests |
|---|---|---|
| `lorren day` | ✅ Working | ✅ domain, repository |
| `lorren train` | ✅ Working | ✅ domain, repository |
| `lorren plan` | ✅ Working | ✅ repository |
| `lorren read` | 💭 Idea | — |
| `lorren cook` | 💭 Idea | — |
| `lorren bars` | 💭 Idea | — |
| Progress dashboard | 💭 Idea | — |

`internal/interviewer` is deliberately left out of the test suite, as it depends on interactive terminal I/O.

**Known gaps:** `lorren train` overwrites when the same session is logged twice in one day, and does not yet seed the wizard with the weights from your last session.

## Built with

- [Go](https://go.dev/) 1.27+
- [Cobra](https://github.com/spf13/cobra) — CLI command routing
- [Viper](https://github.com/spf13/viper) — configuration
- [huh](https://github.com/charmbracelet/huh) — interactive terminal forms
- [GitHub Actions](https://github.com/features/actions) — CI (build, test, lint on every push to main)

## License

MIT