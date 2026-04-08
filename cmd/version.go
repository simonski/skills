package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	ver "github.com/simonski/skills/internal/version"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the skills version",
	Long:  `Print the current version of the skills binary.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runVersion()
	},
}

func runVersion() error {
	fmt.Printf("skills version %s\n", Version)

	latest, err := ver.LatestRelease()
	if err != nil {
		// Non-fatal: silently skip update check if offline.
		return nil
	}
	if latest == "" {
		return nil
	}

	if Version != "dev" && ver.IsNewerThan(latest, Version) {
		yellow := color.New(color.FgYellow).SprintFunc()
		fmt.Printf("%s\n", yellow(fmt.Sprintf("A newer version is available: v%s\nRun: brew upgrade simonski/tap/skills", latest)))
	} else if Version != "dev" {
		green := color.New(color.FgGreen).SprintFunc()
		fmt.Printf("%s\n", green("You are running the latest version."))
	}
	return nil
}
