class Am < Formula
  desc "Lightweight CLI tool for managing shell aliases"
  homepage "https://github.com/navio/am"
  url "https://github.com/navio/am/archive/refs/tags/v1.1.0.tar.gz"
  sha256 "4ab5ad31b6e885dae8ea2af1e97e7f225fbb5a029f6875937d7437cd369a7781"
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
    assert_match "1.1.0", output

    # Test that help works
    help_output = shell_output("#{bin}/am --help")
    assert_match "Alias Manager", help_output
  end
end
