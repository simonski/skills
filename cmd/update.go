package cmd

import (
"fmt"
"os"

"github.com/fatih/color"
"github.com/spf13/cobra"

"github.com/simonski/skills/internal/catalog"
"github.com/simonski/skills/internal/project"
)

var updateConfirm bool

var updateCmd = &cobra.Command{
Use:   "update [<skill-id>]",
Short: "Update installed skills to their latest catalog versions",
Long: `Update installed skills to their latest catalog versions.

Without -y, shows what would be updated (dry run).
With -y, applies the updates.

With no arguments, all installed skills are checked.
With a skill ID, only that skill is checked.

Examples:

  skills update          # show what would be updated
  skills update -y       # update all installed skills
  skills update go       # show whether 'go' would be updated
  skills update go -y    # update only the 'go' skill`,
Args: cobra.MaximumNArgs(1),
RunE: func(cmd *cobra.Command, args []string) error {
cwd, err := os.Getwd()
if err != nil {
return fmt.Errorf("getting current directory: %w", err)
}
if len(args) == 1 {
return runUpdateOne(cwd, args[0], updateConfirm)
}
return runUpdateAll(cwd, updateConfirm)
},
}

func init() {
updateCmd.Flags().BoolVarP(&updateConfirm, "yes", "y", false, "Apply updates (without this flag, only shows what would change)")
}

func runUpdateAll(root string, apply bool) error {
installed, err := project.List(root)
if err != nil {
return err
}

if len(installed) == 0 {
fmt.Println("No skills are installed in this project.")
return nil
}

yellow := color.New(color.FgYellow, color.Bold).SprintFunc()
green := color.New(color.FgGreen, color.Bold).SprintFunc()
cyan := color.New(color.FgCyan, color.Bold).SprintFunc()

type pending struct {
ins    *project.InstalledSkill
latest *catalog.Skill
}

var toUpdate []pending
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
fmt.Printf("  %s %s  v%s → v%s\n", cyan("↑"), ins.ID, ins.Version, latest.Version)
toUpdate = append(toUpdate, pending{ins, latest})
}

fmt.Println()

if len(toUpdate) == 0 {
fmt.Printf("All %d skill(s) are up to date.\n", upToDate)
return nil
}

if !apply {
fmt.Printf("%d skill(s) would be updated. Run with -y to apply.\n", len(toUpdate))
return nil
}

updated := 0
for _, p := range toUpdate {
content := buildSkillFile(p.latest)
if err := project.Install(root, p.ins.ID, content); err != nil {
return err
}
fmt.Printf("  %s Updated %q v%s → v%s\n", green("✓"), p.ins.ID, p.ins.Version, p.latest.Version)
updated++
}

fmt.Println()
fmt.Printf("Done. %d skill(s) updated, %d already up to date.\n", updated, upToDate)
return nil
}

func runUpdateOne(root, id string, apply bool) error {
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

yellow := color.New(color.FgYellow, color.Bold).SprintFunc()
green := color.New(color.FgGreen, color.Bold).SprintFunc()
cyan := color.New(color.FgCyan, color.Bold).SprintFunc()

if ins.Version == latest.Version {
fmt.Printf("  %s %s is already up to date (v%s)\n", yellow("–"), id, ins.Version)
return nil
}

fmt.Printf("  %s %s  v%s → v%s\n", cyan("↑"), id, ins.Version, latest.Version)

if !apply {
fmt.Println()
fmt.Println("Run with -y to apply.")
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
