---
title: Usage
description: Core Alias Manager workflow and everyday commands.
---

# Usage

Alias Manager edits your existing shell config, keeping comments and spacing intact. The typical flow is add → verify → adjust.

## Add aliases

```bash
am add ll 'ls -la'
am add gs 'git status'
am add serve 'python -m http.server 8000'
```

Name rules: alphanumeric with `_` or `-`. Duplicate names are rejected to prevent accidental overwrites.

## List and search

```bash
am list            # all aliases
am list git        # filter by substring
```

Example output:

```
Aliases matching 'git' (1 found):

  gs → git status
```

## Update

```bash
am update ll 'ls -lah'
```

Updates happen in-place; structure and comments around the alias stay untouched.

## Delete

```bash
am delete serve
```

## Help and version

```bash
am --help
am add --help
am --version
```

## Backup, restore, and sync

```bash
am export -o aliases.json       # back up
am import aliases.json          # restore (conflicts reported)
am sync init --gist             # one-time setup; requires the gh CLI
am sync push                    # upload local aliases
am sync pull                    # bring another machine in sync
```

The [Sync section](/commands#sync-across-machines) covers the full workflow, including how `am sync` reuses import semantics so pulls never clobber unsynced local changes.

## Tips

- Run `am init` once so shells auto-reload after changes.
- Keep related aliases grouped with comments in your dotfile; `am` will preserve them.
- Backups live next to your dotfile with a `.bak` suffix for quick recovery.
- `am sync` needs the [`gh`](https://cli.github.com) CLI — install and `gh auth login` before running `am sync init --gist`.
