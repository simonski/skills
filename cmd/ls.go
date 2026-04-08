package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/simonski/skills/internal/catalog"
	"github.com/simonski/skills/internal/project"
)

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all skills in the catalog",
	Long: `List all skills in the catalog and show their installation status
in the current project.

Status indicators:
  ` + color.GreenString("INSTALLED") + `        — the latest version is installed
  ` + color.YellowString("UPDATE AVAILABLE") + ` — an older version is installed
  ` + color.RedString("NOT INSTALLED") + `    — skill is not installed in this project`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runLS()
	},
}

func runLS() error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting current directory: %w", err)
	}

	catalogSkills, err := catalog.All()
	if err != nil {
		return err
	}

	green := color.New(color.FgGreen, color.Bold).SprintFunc()
	yellow := color.New(color.FgYellow, color.Bold).SprintFunc()
	red := color.New(color.FgRed, color.Bold).SprintFunc()

	fmt.Printf("%-20s %-12s %-20s %s\n", "ID", "VERSION", "STATUS", "DESCRIPTION")
	fmt.Printf("%-20s %-12s %-20s %s\n", "--------------------", "------------", "--------------------", "-----------")

	for _, s := range catalogSkills {
		installed, err := project.Get(cwd, s.ID)
		if err != nil {
			return err
		}

		var statusLabel string
		switch {
		case installed == nil:
			statusLabel = red("NOT INSTALLED")
		case installed.Version == s.Version:
			statusLabel = green("INSTALLED")
		default:
			statusLabel = yellow("UPDATE AVAILABLE")
		}

		fmt.Printf("%-20s %-12s %-30s %s\n", s.ID, s.Version, statusLabel, s.Description)
	}
	return nil
}
