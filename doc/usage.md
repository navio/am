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

## Tips

- Run `am init` once so shells auto-reload after changes.
- Keep related aliases grouped with comments in your dotfile; `am` will preserve them.
- Backups live next to your dotfile with a `.bak` suffix for quick recovery.
