class Am < Formula
  desc "Lightweight CLI tool for managing shell aliases"
  homepage "https://github.com/navio/am"
  url "https://github.com/navio/am/archive/refs/tags/v1.2.0.tar.gz"
  sha256 "71f283938375bb36fb5627d5d8dd700e78ba1c50d93db1fa5f0f0d9c8c526805"
  license "MIT"

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
    assert_match "1.2.0", output

    # Test that help works
    help_output = shell_output("#{bin}/am --help")
    assert_match "Alias Manager", help_output
  end
end
