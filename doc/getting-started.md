---
title: Getting Started
description: Install Alias Manager and run the first commands.
---

# Getting Started

Alias Manager (`am`) is a lightweight CLI for managing shell aliases directly in your dotfiles. Install it, initialize once, and start updating aliases safely.

## Install

### Homebrew (recommended)

```bash
brew install navio/tap/am
```

### Build from source

```bash
git clone https://github.com/navio/am.git
cd am
go build -o am
sudo mv am /usr/local/bin/   # or rely on GOPATH/bin via go install
```

## One-time init

Run `am init` once to add a tiny wrapper that automatically reloads your shell after any alias change:

```bash
am init
```

Then re-source your shell (pick one):

```bash
source ~/.zshrc
source ~/.bashrc
source ~/.config/fish/config.fish
```

From now on, alias changes take effect immediately.

## Quick smoke test

```bash
am add ll 'ls -la'
am list
```

You should see `ll → ls -la` listed. If not, confirm your shell profile path matches your shell and that `am init` ran successfully.

## Optional: sync across machines

If you want your aliases mirrored to other machines, install the [GitHub CLI](https://cli.github.com) (e.g. `brew install gh`), authenticate it once, and then run:

```bash
gh auth login
am sync init --gist
am sync push
```

See the [Sync section in Command Reference](/commands#sync-across-machines) for the full workflow.
