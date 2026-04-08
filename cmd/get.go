package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/simonski/skills/internal/catalog"
)

var getCmd = &cobra.Command{
	Use:   "get <skill-id>",
	Short: "Print a skill from the catalog to STDOUT",
	Long: `Print the full content of a skill from the catalog to STDOUT.

The output includes the YAML front matter and the skill body, suitable
for piping into other tools or inspecting a skill before adding it.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runGet(args[0])
	},
}

func runGet(id string) error {
	skill, err := catalog.Get(id)
	if err != nil {
		return err
	}
	fmt.Print(buildSkillFile(skill))
	return nil
}
