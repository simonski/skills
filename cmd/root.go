package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	ver "github.com/simonski/skills/internal/version"
)

// Version is set at build time via -ldflags.
var Version = "dev"

var rootCmd = &cobra.Command{
	Use:   "skills",
	Short: "Manage agentic skills (SKILL.md) in your project",
	Long: `skills — agentic skills manager

Manage a catalog of AI-agent skill definitions (SKILL.md files) in your project.
Skills can be listed, searched, added and removed.

Run 'skills <command> --help' for more information about a specific command.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Skip the version check for the version command itself.
		if cmd.Use == "version" {
			return
		}
		checkForUpdates()
	},
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.AddCommand(lsCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(rmCmd)
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(versionsCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(getCmd)
}

// checkForUpdates prints a notice when a newer release is available.
func checkForUpdates() {
	outdated, latest, err := ver.IsOutdated(Version)
	if err != nil {
		// Silently ignore network errors during update check.
		return
	}
	if outdated {
		yellow := color.New(color.FgYellow).SprintfFunc()
		fmt.Fprintf(os.Stderr, "%s\n\n",
			yellow("A newer version of skills is available: v%s (current: v%s)\nRun: brew upgrade simonski/tap/skills", latest, Version))
	}
}
