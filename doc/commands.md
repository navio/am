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
