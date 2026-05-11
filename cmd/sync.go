package cmd

import (
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/navio/am/internal/alias"
	"github.com/navio/am/internal/config"
	syncpkg "github.com/navio/am/internal/sync"
	"github.com/spf13/cobra"
)

var (
	syncInitGist        bool
	syncInitGistID      string
	syncInitDescription string
	syncPullOverwrite   bool
	syncPullDryRun      bool
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync aliases across machines via a private GitHub Gist",
	Long: `Sync your aliases between machines using a private GitHub Gist.

REQUIREMENTS:
  The GitHub CLI (gh) MUST be installed and authenticated for 'am sync' to
  work. Install it from https://cli.github.com and run 'gh auth login' once.
  Every subcommand below (init / push / pull / status*) shells out to 'gh'.

  *'am sync status' works without gh - it only reads local config.

The sync feature reuses the JSON format produced by 'am export' and stores
its configuration in ~/.config/am/config.json (override with $AM_CONFIG_DIR
or $XDG_CONFIG_HOME).

Examples:
  am sync init --gist                  # create a new private gist
  am sync init --gist --id <gistID>    # attach to an existing gist
  am sync push                         # upload local aliases
  am sync pull                         # download (skips conflicts)
  am sync pull --overwrite             # download and replace conflicts
  am sync status                       # show remote config & last sync`,
}

var syncInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Configure a remote sync destination",
	RunE: func(cmd *cobra.Command, args []string) error {
		if !syncInitGist {
			return fmt.Errorf("only --gist is supported at this time")
		}

		cfg, err := config.Load()
		if err != nil {
			return err
		}
		if cfg.Sync != nil && cfg.Sync.GistID != "" && syncInitGistID == "" {
			return fmt.Errorf("sync is already configured for gist %s. Re-run with --id to switch, or edit %s", cfg.Sync.GistID, configPathOrEmpty())
		}

		client, err := syncpkg.NewGistClient()
		if err != nil {
			return err
		}

		filename := syncpkg.DefaultGistFilename
		var gistID string

		if syncInitGistID != "" {
			if _, err := client.View(syncInitGistID, filename); err != nil {
				return fmt.Errorf("could not access gist %s with file %s: %w", syncInitGistID, filename, err)
			}
			gistID = syncInitGistID
		} else {
			manager := alias.NewManager()
			aliases, err := manager.List("")
			if err != nil {
				return err
			}
			payload, err := syncpkg.Marshal(syncpkg.AliasesToMap(aliases))
			if err != nil {
				return err
			}
			desc := syncInitDescription
			if desc == "" {
				desc = "am - shell aliases sync"
			}
			id, err := client.Create(filename, payload, desc)
			if err != nil {
				return err
			}
			gistID = id
		}

		cfg.Sync = &config.SyncConfig{
			Provider:     "gist",
			GistID:       gistID,
			GistFilename: filename,
		}
		if err := config.Save(cfg); err != nil {
			return err
		}

		green := color.New(color.FgGreen, color.Bold)
		green.Print("✓ ")
		fmt.Printf("Sync configured for gist %s\n", gistID)
		fmt.Println("  Next: 'am sync push' to upload, or 'am sync pull' to download.")
		return nil
	},
}

var syncPushCmd = &cobra.Command{
	Use:   "push",
	Short: "Upload local aliases to the remote",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, sync, err := requireSyncConfig()
		if err != nil {
			return err
		}

		client, err := syncpkg.NewGistClient()
		if err != nil {
			return err
		}

		manager := alias.NewManager()
		aliases, err := manager.List("")
		if err != nil {
			return err
		}
		payload, err := syncpkg.Marshal(syncpkg.AliasesToMap(aliases))
		if err != nil {
			return err
		}

		if err := client.Edit(sync.GistID, sync.GistFilename, payload); err != nil {
			return err
		}

		sync.LastSynced = time.Now().UTC().Format(time.RFC3339)
		sync.LastHash = syncpkg.Hash(payload)
		cfg.Sync = sync
		if err := config.Save(cfg); err != nil {
			return err
		}

		green := color.New(color.FgGreen, color.Bold)
		green.Print("✓ ")
		fmt.Printf("Pushed %d aliases to gist %s\n", len(aliases), sync.GistID)
		return nil
	},
}

var syncPullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Download aliases from the remote",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, sync, err := requireSyncConfig()
		if err != nil {
			return err
		}

		client, err := syncpkg.NewGistClient()
		if err != nil {
			return err
		}

		data, err := client.View(sync.GistID, sync.GistFilename)
		if err != nil {
			return err
		}
		entries, err := syncpkg.Unmarshal(data)
		if err != nil {
			return err
		}
		if len(entries) == 0 {
			yellow := color.New(color.FgYellow)
			yellow.Println("Remote gist is empty - nothing to pull.")
			return nil
		}

		manager := alias.NewManager()
		result, err := manager.ImportMany(entries, syncPullOverwrite, syncPullDryRun)
		if err != nil {
			return err
		}

		green := color.New(color.FgGreen, color.Bold)
		yellow := color.New(color.FgYellow)
		cyan := color.New(color.FgCyan)

		if syncPullDryRun {
			cyan.Println("Dry run - no changes were written.")
			fmt.Println()
		} else {
			green.Print("✓ ")
			fmt.Print("Pull complete!\n\n")
		}

		fmt.Printf("  Added:     %d\n", len(result.Added))
		fmt.Printf("  Updated:   %d\n", len(result.Updated))
		fmt.Printf("  Skipped:   %d\n", len(result.Skipped))
		fmt.Printf("  Conflicts: %d\n", len(result.Conflicts))

		if len(result.Conflicts) > 0 {
			fmt.Println()
			yellow.Println("Conflicts (pass --overwrite to replace):")
			for _, name := range result.Conflicts {
				fmt.Printf("  - %s\n", name)
			}
		}

		if !syncPullDryRun {
			sync.LastSynced = time.Now().UTC().Format(time.RFC3339)
			sync.LastHash = syncpkg.Hash(data)
			cfg.Sync = sync
			if err := config.Save(cfg); err != nil {
				return err
			}
		}
		return nil
	},
}

var syncStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show sync configuration and last sync info",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		path, _ := config.Path()
		cyan := color.New(color.FgCyan, color.Bold)

		if cfg.Sync == nil || cfg.Sync.GistID == "" {
			cyan.Println("Sync is not configured.")
			fmt.Println("Run 'am sync init --gist' to get started.")
			fmt.Printf("Config file: %s\n", path)
			return nil
		}

		cyan.Println("Sync status:")
		fmt.Printf("  Provider:    %s\n", cfg.Sync.Provider)
		fmt.Printf("  Gist ID:     %s\n", cfg.Sync.GistID)
		fmt.Printf("  Filename:    %s\n", cfg.Sync.GistFilename)
		if cfg.Sync.LastSynced != "" {
			fmt.Printf("  Last synced: %s\n", cfg.Sync.LastSynced)
		} else {
			fmt.Println("  Last synced: never")
		}
		if cfg.Sync.LastHash != "" {
			fmt.Printf("  Last hash:   %s\n", cfg.Sync.LastHash)
		}
		fmt.Printf("  Config file: %s\n", path)
		return nil
	},
}

func requireSyncConfig() (*config.Config, *config.SyncConfig, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, nil, err
	}
	if cfg.Sync == nil || cfg.Sync.GistID == "" {
		return nil, nil, fmt.Errorf("sync is not configured. Run 'am sync init --gist' first")
	}
	sync := *cfg.Sync
	if sync.GistFilename == "" {
		sync.GistFilename = syncpkg.DefaultGistFilename
	}
	return cfg, &sync, nil
}

func configPathOrEmpty() string {
	p, err := config.Path()
	if err != nil {
		return ""
	}
	return p
}

func init() {
	syncInitCmd.Flags().BoolVar(&syncInitGist, "gist", false, "Use a GitHub Gist as the sync target")
	syncInitCmd.Flags().StringVar(&syncInitGistID, "id", "", "Attach to an existing gist by ID instead of creating one")
	syncInitCmd.Flags().StringVar(&syncInitDescription, "desc", "", "Description for the newly created gist")

	syncPullCmd.Flags().BoolVar(&syncPullOverwrite, "overwrite", false, "Replace local aliases on conflict")
	syncPullCmd.Flags().BoolVar(&syncPullDryRun, "dry-run", false, "Preview changes without writing to disk")

	syncCmd.AddCommand(syncInitCmd)
	syncCmd.AddCommand(syncPushCmd)
	syncCmd.AddCommand(syncPullCmd)
	syncCmd.AddCommand(syncStatusCmd)

	rootCmd.AddCommand(syncCmd)
}
