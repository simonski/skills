package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/simonski/skills/internal/catalog"
	"github.com/simonski/skills/internal/project"
)

var addCmd = &cobra.Command{
	Use:   "add <skill-id>[@<version>]",
	Short: "Add a skill to the current project",
	Long: `Add a skill from the catalog to the current project.

The skill file is written to .skills/<skill-id>/SKILL.md.
If the skill is already installed, it is updated to the specified version.

To install the latest version:

  skills add go

To pin to a specific version:

  skills add go@0.0.1`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAdd(args[0])
	},
}

func runAdd(idArg string) error {
	id, version, err := parseSkillArg(idArg)
	if err != nil {
		return err
	}

	var skill *catalog.Skill
	if version != "" {
		skill, err = catalog.GetVersion(id, version)
	} else {
		skill, err = catalog.Get(id)
	}
	if err != nil {
		return err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting current directory: %w", err)
	}

	existing, err := project.Get(cwd, id)
	if err != nil {
		return err
	}

	content := buildSkillFile(skill)
	if err := project.Install(cwd, id, content); err != nil {
		return err
	}

	green := color.New(color.FgGreen, color.Bold).SprintFunc()
	if existing != nil {
		fmt.Printf("%s Updated skill %q to v%s\n", green("✓"), id, skill.Version)
	} else {
		fmt.Printf("%s Added skill %q (v%s)\n", green("✓"), id, skill.Version)
	}
	fmt.Printf("  Written to: %s\n", project.SkillPath(cwd, id))
	return nil
}

// buildSkillFile renders the full SKILL.md content including front matter.
func buildSkillFile(s *catalog.Skill) string {
	return fmt.Sprintf("---\nid: %s\nversion: %s\ndescription: %s\n---\n\n%s\n", s.ID, s.Version, s.Description, s.Content)
}

// parseSkillArg splits "id@version" into (id, version).
// version is empty when no @ is present (meaning: use latest).
func parseSkillArg(arg string) (id, version string, err error) {
	parts := strings.SplitN(arg, "@", 2)
	id = parts[0]
	if id == "" {
		return "", "", fmt.Errorf("skill id must not be empty")
	}
	if len(parts) == 2 {
		version = parts[1]
		if version == "" {
			return "", "", fmt.Errorf("version must not be empty after @")
		}
	}
	return id, version, nil
}
