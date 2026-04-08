package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/simonski/skills/internal/project"
)

var rmCmd = &cobra.Command{
	Use:   "rm <skill-id>",
	Short: "Remove a skill from the current project",
	Long: `Remove an installed skill from the current project.

Deletes the .skills/<skill-id>.md file.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runRM(args[0])
	},
}

func runRM(id string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting current directory: %w", err)
	}

	if err := project.Remove(cwd, id); err != nil {
		return err
	}

	red := color.New(color.FgRed, color.Bold).SprintFunc()
	fmt.Printf("%s Removed skill %q\n", red("✗"), id)
	return nil
}
