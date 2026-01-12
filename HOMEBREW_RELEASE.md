# Homebrew Release Guide

## Step-by-Step Guide to Release `am` on Homebrew

### Step 1: Create a GitHub Release

1. **Go to your repository's releases page:**
   https://github.com/navio/am/releases

2. **Click "Create a new release"**

3. **Fill in the release details:**
   - **Tag:** `v1.0.0` (already exists, select it)
   - **Title:** `v1.0.0 - Initial Release`
   - **Description:** Use this template:

```markdown
# Alias Manager (am) v1.0.0

A lightweight CLI tool for managing shell aliases with direct dotfile manipulation.

## Features
- ✅ Add, list, update, and delete aliases
- ✅ Auto-detect shell (bash/zsh)
- ✅ Automatic backups before modifications
- ✅ Auto-source wrapper (`am init`) for immediate alias activation
- ✅ Search aliases by name
- ✅ Preserve comments and file structure

## Installation

### Via Homebrew
```bash
brew tap navio/tap
brew install am
```

### From Source
```bash
go install github.com/navio/am@latest
```

## Quick Start

1. Initialize auto-sourcing:
   ```bash
   am init
   source ~/.zshrc  # or ~/.bashrc
   ```

2. Start using:
   ```bash
   am add ll 'ls -la'
   am list
   ```

## What's New
Initial release with core functionality for managing shell aliases.

## Test Coverage
- Shell Detection: 95.7%
- Alias Operations: 88.2%
- File Handling: 61.6%
```

4. **Publish the release**

### Step 2: Get the SHA256 Hash

Run this command to get the SHA256 of your release tarball:

```bash
curl -L https://github.com/navio/am/archive/refs/tags/v1.0.0.tar.gz | shasum -a 256
```

Copy the resulting hash (it will look like: `a1b2c3d4e5f6...`)

### Step 3: Create a Homebrew Tap Repository

1. **Create a new GitHub repository:**
   - Repository name: `homebrew-tap`
   - Full name will be: `navio/homebrew-tap`
   - Make it public
   - Don't add README, .gitignore, or license (we'll add them manually)

2. **Clone it locally:**
```bash
cd ~/Development
git clone https://github.com/navio/homebrew-tap.git
cd homebrew-tap
```

3. **Create the Formula directory:**
```bash
mkdir Formula
```

### Step 4: Create the Formula File

Create the file `Formula/am.rb` with this content (replace the SHA256):

```ruby
class Am < Formula
  desc "Lightweight CLI tool for managing shell aliases"
  homepage "https://github.com/navio/am"
  url "https://github.com/navio/am/archive/refs/tags/v1.0.0.tar.gz"
  sha256 "PASTE_YOUR_SHA256_HERE"  # ← Replace this with the actual SHA256
  license "MIT"

  bottle :unneeded

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w"), "-o", bin/"am"
  end

  def caveats
    <<~EOS
      To enable auto-sourcing of aliases (recommended), run:
        am init

      Then restart your shell or run:
        source ~/.zshrc    # for zsh
        source ~/.bashrc   # for bash

      From that point forward, all alias changes will take effect immediately!
    EOS
  end

  test do
    # Test that the binary runs and shows version
    output = shell_output("#{bin}/am --version")
    assert_match "1.0.0", output

    # Test that help works
    help_output = shell_output("#{bin}/am --help")
    assert_match "Alias Manager", help_output
  end
end
```

### Step 5: Test the Formula Locally

Before publishing, test that it works:

```bash
# Test installation
brew install --build-from-source ./Formula/am.rb

# Test the installed binary
am --version
am --help

# Uninstall after testing
brew uninstall am
```

### Step 6: Publish Your Tap

```bash
# Add and commit the formula
git add Formula/am.rb
git commit -m "Add am formula v1.0.0"
git push origin main
```

### Step 7: Add a README to Your Tap

Create `README.md` in your tap repository:

```markdown
# Homebrew Tap for navio

## Installation

```bash
brew tap navio/tap
brew install am
```

## Formulas

### am - Alias Manager

A lightweight CLI tool for managing shell aliases.

[GitHub Repository](https://github.com/navio/am)
```

Then commit and push:

```bash
git add README.md
git commit -m "Add README"
git push origin main
```

### Step 8: Test Installation from Your Tap

Now test that users can install it:

```bash
# Untap if you already tapped it
brew untap navio/tap 2>/dev/null || true

# Tap your repository
brew tap navio/tap

# Install
brew install am

# Verify it works
am --version
am --help
```

### Step 9: Update Your Main Repository README

The Homebrew installation section in your main `am` repository README should already be correct:

```markdown
### Via Homebrew

```bash
brew tap navio/tap
brew install am
```
```

## For Future Releases

When you release v1.1.0, v2.0.0, etc.:

1. **Create a new git tag and release:**
   ```bash
   git tag v1.1.0
   git push origin v1.1.0
   ```

2. **Get the new SHA256:**
   ```bash
   curl -L https://github.com/navio/am/archive/refs/tags/v1.1.0.tar.gz | shasum -a 256
   ```

3. **Update the formula in homebrew-tap:**
   ```bash
   cd ~/Development/homebrew-tap
   # Edit Formula/am.rb
   # - Update the version in the URL: v1.0.0 → v1.1.0
   # - Update the sha256
   git add Formula/am.rb
   git commit -m "Update am to v1.1.0"
   git push origin main
   ```

4. **Users can then upgrade:**
   ```bash
   brew update
   brew upgrade am
   ```

## Troubleshooting

### Formula Audit

Before publishing, you can audit your formula:

```bash
brew audit --strict Formula/am.rb
```

### Testing in a Clean Environment

Test installation in a clean environment:

```bash
# Create a temporary directory
mkdir -p /tmp/brew-test
cd /tmp/brew-test

# Install from your tap
brew tap navio/tap
brew install am

# Test it
am --version
```

## Quick Command Reference

```bash
# Get SHA256 for release
curl -L https://github.com/navio/am/archive/refs/tags/v1.0.0.tar.gz | shasum -a 256

# Test formula locally
brew install --build-from-source ./Formula/am.rb

# Audit formula
brew audit --strict Formula/am.rb

# Uninstall
brew uninstall am

# Untap
brew untap navio/tap
```

## Resources

- [Homebrew Formula Cookbook](https://docs.brew.sh/Formula-Cookbook)
- [Homebrew Acceptable Formulae](https://docs.brew.sh/Acceptable-Formulae)
- [How to Create Homebrew Tap](https://docs.brew.sh/How-to-Create-and-Maintain-a-Tap)
