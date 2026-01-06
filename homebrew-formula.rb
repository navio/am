# This file is a template for the Homebrew formula
# When you're ready to publish, create a tap repository and use this formula
#
# Usage:
# 1. Create a repository named homebrew-tap (e.g., navio/homebrew-tap)
# 2. Add this formula to the repository as Formula/alias-manager.rb
# 3. Update the URL and SHA256 after creating a release
# 4. Users can then install with:
#    brew tap navio/tap
#    brew install alias-manager

class AliasManager < Formula
  desc "Lightweight CLI tool for managing shell aliases"
  homepage "https://github.com/navio/alias-manager"

  # Update these values when you create a release
  url "https://github.com/navio/alias-manager/archive/refs/tags/v1.0.0.tar.gz"
  sha256 "REPLACE_WITH_ACTUAL_SHA256"
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w -X main.version=#{version}"), "-o", bin/"am"
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
    assert_match version.to_s, output

    # Test that help works
    help_output = shell_output("#{bin}/am --help")
    assert_match "Alias Manager", help_output
  end
end

# Installation Instructions:
# ==========================
#
# 1. Create a GitHub Release:
#    - Tag your code with a version (e.g., v1.0.0)
#    - Create a release on GitHub
#    - GitHub will automatically create a tarball at:
#      https://github.com/navio/alias-manager/archive/refs/tags/v1.0.0.tar.gz
#
# 2. Get the SHA256:
#    curl -L https://github.com/navio/alias-manager/archive/refs/tags/v1.0.0.tar.gz | shasum -a 256
#
# 3. Update the formula:
#    - Replace the URL with your actual release URL
#    - Replace the SHA256 with the value from step 2
#
# 4. Create a tap repository:
#    - Create a new repo: navio/homebrew-tap
#    - Create directory: Formula/
#    - Add this file as: Formula/alias-manager.rb
#
# 5. Test the formula locally:
#    brew install --build-from-source ./homebrew-formula.rb
#
# 6. Users can then install with:
#    brew tap navio/tap
#    brew install alias-manager
#
# 7. To update the formula after a new release:
#    - Update the version number in the URL
#    - Update the SHA256
#    - Commit and push to your tap repository
