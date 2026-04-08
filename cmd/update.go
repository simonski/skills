package cmd

import (
"fmt"
"os"

"github.com/fatih/color"
"github.com/spf13/cobra"

"github.com/simonski/skills/internal/catalog"
"github.com/simonski/skills/internal/project"
)

var updateCmd = &cobra.Command{
Use:   "update [<skill-id>]",
Short: "Update installed skills to their latest catalog versions",
Long: `Update installed skills to their latest catalog versions.

With no arguments, all installed skills that have a newer catalog version are updated.
With a skill ID, only that skill is updated.

Examples:

  skills update          # update all installed skills
  skills update go       # update only the 'go' skill`,
Args: cobra.MaximumNArgs(1),
RunE: func(cmd *cobra.Command, args []string) error {
cwd, err := os.Getwd()
if err != nil {
return fmt.Errorf("getting current directory: %w", err)
}
if len(args) == 1 {
return runUpdateOne(cwd, args[0])
}
return runUpdateAll(cwd)
},
}

func runUpdateAll(root string) error {
installed, err := project.List(root)
if err != nil {
return err
}

if len(installed) == 0 {
fmt.Println("No skills are installed in this project.")
return nil
}

green := color.New(color.FgGreen, color.Bold).SprintFunc()
yellow := color.New(color.FgYellow, color.Bold).SprintFunc()

updated := 0
upToDate := 0

for _, ins := range installed {
latest, err := catalog.Get(ins.ID)
if err != nil {
fmt.Fprintf(os.Stderr, "  warning: skill %q not found in catalog, skipping\n", ins.ID)
continue
}

if ins.Version == latest.Version {
fmt.Printf("  %s %s is already up to date (v%s)\n", yellow("–"), ins.ID, ins.Version)
upToDate++
continue
}

content := buildSkillFile(latest)
if err := project.Install(root, ins.ID, content); err != nil {
return err
}
fmt.Printf("  %s Updated %q v%s → v%s\n", green("✓"), ins.ID, ins.Version, latest.Version)
updated++
}

fmt.Println()
if updated == 0 && upToDate == 0 {
fmt.Println("Nothing to update.")
} else {
fmt.Printf("Done. %d skill(s) updated, %d already up to date.\n", updated, upToDate)
}
return nil
}

func runUpdateOne(root, id string) error {
ins, err := project.Get(root, id)
if err != nil {
return err
}
if ins == nil {
return fmt.Errorf("skill %q is not installed", id)
}

latest, err := catalog.Get(id)
if err != nil {
return err
}

green := color.New(color.FgGreen, color.Bold).SprintFunc()
yellow := color.New(color.FgYellow, color.Bold).SprintFunc()

if ins.Version == latest.Version {
fmt.Printf("  %s %s is already up to date (v%s)\n", yellow("–"), id, ins.Version)
return nil
}

content := buildSkillFile(latest)
if err := project.Install(root, id, content); err != nil {
return err
}
fmt.Printf("  %s Updated %q v%s → v%s\n", green("✓"), id, ins.Version, latest.Version)
fmt.Printf("  Written to: %s\n", project.SkillPath(root, id))
return nil
}
