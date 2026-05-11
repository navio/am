---
title: Command Reference
description: All CLI commands for Alias Manager with examples.
---

# Command Reference

| Command | Description | Example |
| --- | --- | --- |
| `am init` | Add auto-source wrapper (run once) | `am init` |
| `am add <name> <command>` | Add a new alias | `am add ll 'ls -la'` |
| `am list [search]` | List all aliases or filter by substring | `am list git` |
| `am update <name> <command>` | Update an existing alias | `am update gs 'git status -sb'` |
| `am delete <name>` | Remove an alias | `am delete serve` |
| `am export [--format json\|shell] [-o FILE]` | Export aliases (JSON to stdout by default) | `am export -o aliases.json` |
| `am import <file> [--overwrite] [--dry-run]` | Import aliases from a JSON file or stdin | `am import aliases.json` |
| `am sync init --gist [--id <id>]` | Configure a private Gist as sync target<br/>**requires [`gh`](https://cli.github.com)** | `am sync init --gist` |
| `am sync push` | Upload local aliases to the configured Gist | `am sync push` |
| `am sync pull [--overwrite] [--dry-run]` | Download aliases (reuses import semantics) | `am sync pull --dry-run` |
| `am sync status` | Show remote config and last-sync info | `am sync status` |
| `am --version` | Show version | `am --version` |
| `am --help` | Help for any command | `am add --help` |

## Examples

### Development aliases

```bash
am add gc "git commit -m"
am add gp "git push"
am add gl "git log --oneline --graph"
```

### Navigation shortcuts

```bash
am add .. 'cd ..'
am add ... 'cd ../..'
am add home 'cd ~'
```

### Utility helpers

```bash
am add serve 'python -m http.server 8000'
am add ports "lsof -i -P | grep LISTEN"
```

## Sync across machines

`am sync` syncs your aliases between machines through a private GitHub Gist. It reuses the JSON format produced by `am export`, and the import semantics from `am import`, so pulls never silently clobber local changes.

::: warning Requires the GitHub CLI
**`am sync` requires the [`gh`](https://cli.github.com) CLI to be installed and authenticated.** Run `gh auth login` once before using any sync subcommand (except `status`, which only reads local config).
:::

### Typical flow

```bash
# Machine A — set up sync (creates a new private gist)
am sync init --gist
am sync push

# Machine B — attach to the same gist and pull
am sync init --gist --id <gistID>
am sync pull
```

### Subcommands

- **`am sync init --gist`** — creates a new private gist via `gh gist create`. Pass `--id <gistID>` to attach to an existing one (handy for additional machines). Optionally set `--desc "..."` to label the gist.
- **`am sync push`** — uploads the current aliases as canonical, sorted JSON to the configured gist.
- **`am sync pull`** — downloads the remote payload and merges via the same logic as `am import`. By default conflicts are reported, not applied; use `--overwrite` to replace local aliases on conflict, or `--dry-run` to preview.
- **`am sync status`** — prints provider, gist ID, last sync timestamp, content hash, and config-file path. Safe to run without `gh`.

### Configuration

Sync state lives in `~/.config/am/config.json` (respects `$XDG_CONFIG_HOME`; `AM_CONFIG_DIR` overrides everything, mainly used in tests). Only the gist ID, filename, and last-sync metadata are stored — no tokens.
