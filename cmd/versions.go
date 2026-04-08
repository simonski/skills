package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/simonski/skills/internal/catalog"
)

var versionsCmd = &cobra.Command{
	Use:   "versions <skill-id>",
	Short: "List all available versions of a skill",
	Long: `List all versions of a skill available in the catalog, sorted oldest to newest.

Example:

  skills versions go`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runVersions(args[0])
	},
}

func runVersions(id string) error {
	versions, err := catalog.Versions(id)
	if err != nil {
		return err
	}
	if len(versions) == 0 {
		fmt.Printf("No versions found for skill %q\n", id)
		return nil
	}
	for _, v := range versions {
		fmt.Println(v)
	}
	return nil
}
