---
title: Internals & Safety
description: How Alias Manager edits your dotfiles while keeping them safe.
---

# Internals & Safety

Alias Manager is intentionally conservative when touching your shell configuration.

## What it does

1. Detects your shell via `$SHELL`.
2. Locates the matching config file (`~/.zshrc`, `~/.bashrc`, or `~/.config/fish/config.fish`).
3. Creates a `.bak` backup before any write.
4. Parses alias lines, leaving comments and spacing untouched.
5. Writes changes atomically to avoid partial updates.

## Reliability features

- **Backups on every write**: restore with `cp ~/.zshrc.bak ~/.zshrc` (or bash/fish equivalents).
- **Duplicate detection**: prevents adding an alias that already exists.
- **Validation**: alias names must be alphanumeric with `_` or `-`.
- **Structure preservation**: only alias lines move; comments and blank lines stay.

## Recovery

If something looks off, copy the backup back into place and re-source your shell:

```bash
cp ~/.zshrc.bak ~/.zshrc
source ~/.zshrc
```

The same approach works for bash and fish.
