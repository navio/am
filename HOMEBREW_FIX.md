# Fix for Homebrew Xcode 26.0 Installation Error

## Problem

Users installing `am` via Homebrew get this error:

```
Error: Your Command Line Tools are too outdated.
Update them from Software Update in System Settings.
...
You should download the Command Line Tools for Xcode 26.0.
```

## Root Cause

Homebrew's automated bottle (pre-compiled binary) building system created bottles with incorrect SDK version requirements that reference a non-existent Xcode 26.0. When users try to install, Homebrew attempts to use these bottles and fails with the SDK version error.

## Solution

Disable bottles and force building from source by adding `bottle :unneeded` to the formula. This is appropriate for `am` because:
- It's a lightweight Go application that compiles quickly (< 10 seconds)
- Go produces static binaries that work across macOS versions
- Building from source avoids bottle SDK compatibility issues

## Steps to Fix

### 1. Update the formula in your homebrew-tap repository

Clone or navigate to your tap repository:
```bash
cd ~/path/to/homebrew-tap
```

Edit `Formula/am.rb` and add the `bottle :unneeded` line after the license declaration:

```ruby
class Am < Formula
  desc "Lightweight CLI tool for managing shell aliases"
  homepage "https://github.com/navio/am"
  url "https://github.com/navio/am/archive/refs/tags/v1.1.0.tar.gz"
  sha256 "4ab5ad31b6e885dae8ea2af1e97e7f225fbb5a029f6875937d7437cd369a7781"
  license "MIT"

  bottle :unneeded  # ← Add this line

  depends_on "go" => :build

  # ... rest of the formula
end
```

### 2. Commit and push the fix

```bash
git add Formula/am.rb
git commit -m "Fix Xcode 26.0 error by disabling bottles"
git push origin main
```

### 3. Users need to update and reinstall

After pushing the fix, users experiencing the issue should run:

```bash
# Update Homebrew to get the latest formula
brew update

# Uninstall the problematic version
brew uninstall am

# Clean up any cached bottles
brew cleanup am

# Reinstall (will now build from source)
brew install navio/tap/am
```

## Prevention for Future Releases

Always include `bottle :unneeded` in the formula for future version updates to prevent this issue from recurring. The updated formula template should include this line by default.

## Technical Details

- **Bottle**: Pre-compiled binary packages that Homebrew creates for faster installation
- **bottle :unneeded**: Directive that tells Homebrew to always build from source instead of using bottles
- **Why it happened**: Homebrew's bottle building infrastructure may have used an incorrect or future SDK version (26.0) that doesn't exist
- **Why disable bottles**: Go applications build quickly and produce static binaries, making source builds practical

## Testing

After applying the fix, test that installation works:

```bash
# In a clean environment
brew tap navio/tap
brew install am --verbose  # Verbose shows build from source

# Verify it works
am --version
am --help
```

You should see compilation output instead of bottle download messages.
