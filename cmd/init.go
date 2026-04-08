package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/simonski/skills/internal/catalog"
	"github.com/simonski/skills/internal/project"
)

// knownAgents is a list of coding agents with their detection heuristics.
var knownAgents = []agentDef{
	{
		Name: "GitHub Copilot",
		Paths: []string{
			".github/copilot-instructions.md",
			".github/copilot",
		},
	},
	{
		Name: "Cursor",
		Paths: []string{
			".cursor",
			".cursorignore",
			".cursorrules",
		},
	},
	{
		Name: "Claude (Anthropic)",
		Paths: []string{
			"CLAUDE.md",
			".claude",
		},
	},
	{
		Name: "Windsurf",
		Paths: []string{
			".windsurf",
			".windsurfrules",
		},
	},
	{
		Name: "Codeium",
		Paths: []string{
			".codeium",
		},
	},
	{
		Name: "Continue",
		Paths: []string{
			".continue",
		},
	},
	{
		Name: "Aider",
		Paths: []string{
			".aider.conf.yml",
			".aiderignore",
		},
	},
}

type agentDef struct {
	Name  string
	Paths []string
}

// skillEntry pairs a catalog skill with its installed state.
type skillEntry struct {
	skill     *catalog.Skill
	installed *project.InstalledSkill
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Interactive wizard to set up skills for this project",
	Long: `Inspect the current folder for coding agents and guide you through
selecting which skills to install or upgrade.

Run this command at any time to review and update your installed skills.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runInit()
	},
}

func runInit() error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting current directory: %w", err)
	}

	bold := color.New(color.Bold).SprintFunc()
	green := color.New(color.FgGreen, color.Bold).SprintFunc()
	yellow := color.New(color.FgYellow, color.Bold).SprintFunc()
	cyan := color.New(color.FgCyan, color.Bold).SprintFunc()

	fmt.Println()
	fmt.Printf("%s\n", bold("skills init — project setup wizard"))
	fmt.Printf("Working directory: %s\n", cwd)
	fmt.Println()

	// 1. Detect coding agents.
	detected := detectAgents(cwd)
	if len(detected) > 0 {
		fmt.Printf("%s\n", bold("Detected coding agents:"))
		for _, name := range detected {
			fmt.Printf("  %s %s\n", green("✓"), name)
		}
	} else {
		fmt.Printf("%s\n", yellow("No coding agents detected in this directory."))
		fmt.Println("  (Skills can still be installed for any agent you plan to use.)")
	}
	fmt.Println()

	// 2. Load catalog and current installation state.
	catalogSkills, err := catalog.All()
	if err != nil {
		return err
	}

	entries := make([]skillEntry, 0, len(catalogSkills))
	for _, s := range catalogSkills {
		ins, err := project.Get(cwd, s.ID)
		if err != nil {
			return err
		}
		entries = append(entries, skillEntry{skill: s, installed: ins})
	}

	// 3. Display current status table.
	fmt.Printf("%s\n", bold("Available skills:"))
	fmt.Printf("  %-4s %-20s %-12s %-32s %s\n", "#", "ID", "VERSION", "STATUS", "DESCRIPTION")
	fmt.Printf("  %-4s %-20s %-12s %-32s %s\n",
		"----", "--------------------", "------------", "--------------------------------", "-----------")

	for i, e := range entries {
		statusLabel := skillStatusStr(e.skill, e.installed)
		fmt.Printf("  %-4s %-20s %-12s %-32s %s\n",
			cyan(fmt.Sprintf("[%d]", i+1)),
			e.skill.ID,
			e.skill.Version,
			statusLabel,
			e.skill.Description,
		)
	}
	fmt.Println()

	// 4. Build default selection: already-installed skills are pre-selected.
	selected := make([]bool, len(entries))
	for i, e := range entries {
		if e.installed != nil {
			selected[i] = true
		}
	}

	// 5. Interactive selection loop.
	reader := bufio.NewReader(os.Stdin)
	for {
		printInitSelection(entries, selected)
		fmt.Println()
		fmt.Println("Enter a skill number to toggle, 'a' to select all, 'n' to select none,")
		fmt.Println("'i' to install/update selected, or 'q' to quit without changes.")
		fmt.Print("> ")

		line, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("reading input: %w", err)
		}
		input := strings.TrimSpace(line)

		switch strings.ToLower(input) {
		case "q", "quit", "exit":
			fmt.Println("Exiting without changes.")
			return nil
		case "a", "all":
			for i := range selected {
				selected[i] = true
			}
		case "n", "none":
			for i := range selected {
				selected[i] = false
			}
		case "i", "install":
			return applyInitSelection(cwd, entries, selected)
		default:
			// Parse as a number or comma-separated list of numbers.
			toggled := false
			for _, part := range strings.Split(input, ",") {
				part = strings.TrimSpace(part)
				num, parseErr := strconv.Atoi(part)
				if parseErr != nil || num < 1 || num > len(entries) {
					if part != "" {
						fmt.Printf("  Unknown input %q — enter a number between 1 and %d.\n", part, len(entries))
					}
					continue
				}
				selected[num-1] = !selected[num-1]
				toggled = true
			}
			if toggled {
				fmt.Println()
			}
		}
	}
}

// detectAgents returns the names of coding agents detected in root.
func detectAgents(root string) []string {
	var found []string
	for _, agent := range knownAgents {
		for _, rel := range agent.Paths {
			if _, err := os.Stat(filepath.Join(root, rel)); err == nil {
				found = append(found, agent.Name)
				break
			}
		}
	}
	return found
}

// skillStatusStr returns a coloured status label for a skill.
func skillStatusStr(s *catalog.Skill, installed *project.InstalledSkill) string {
	green := color.New(color.FgGreen, color.Bold).SprintFunc()
	yellow := color.New(color.FgYellow, color.Bold).SprintFunc()
	red := color.New(color.FgRed, color.Bold).SprintFunc()

	switch {
	case installed == nil:
		return red("NOT INSTALLED")
	case installed.Version == s.Version:
		return green("INSTALLED")
	default:
		return yellow(fmt.Sprintf("UPDATE AVAILABLE (v%s→v%s)", installed.Version, s.Version))
	}
}

// printInitSelection prints the current skill selection state.
func printInitSelection(entries []skillEntry, selected []bool) {
	bold := color.New(color.Bold).SprintFunc()
	green := color.New(color.FgGreen, color.Bold).SprintFunc()
	cyan := color.New(color.FgCyan, color.Bold).SprintFunc()

	fmt.Printf("%s\n", bold("Current selection:"))
	for i, e := range entries {
		mark := "[ ]"
		if selected[i] {
			mark = green("[✓]")
		}
		statusLabel := skillStatusStr(e.skill, e.installed)
		fmt.Printf("  %s %s %-20s %-12s %-32s %s\n",
			mark,
			cyan(fmt.Sprintf("[%d]", i+1)),
			e.skill.ID,
			e.skill.Version,
			statusLabel,
			e.skill.Description,
		)
	}
}

// applyInitSelection installs all selected skills, skipping those already up to date.
func applyInitSelection(cwd string, entries []skillEntry, selected []bool) error {
	green := color.New(color.FgGreen, color.Bold).SprintFunc()
	yellow := color.New(color.FgYellow, color.Bold).SprintFunc()

	installed := 0
	skipped := 0
	fmt.Println()

	for i, e := range entries {
		if !selected[i] {
			continue
		}
		// Skip if already at the correct version.
		if e.installed != nil && e.installed.Version == e.skill.Version {
			fmt.Printf("  %s %s is already up to date (v%s)\n", yellow("–"), e.skill.ID, e.skill.Version)
			skipped++
			continue
		}

		content := buildSkillFile(e.skill)
		if err := project.Install(cwd, e.skill.ID, content); err != nil {
			return err
		}

		if e.installed != nil {
			fmt.Printf("  %s Updated %q to v%s\n", green("✓"), e.skill.ID, e.skill.Version)
		} else {
			fmt.Printf("  %s Added %q (v%s)\n", green("✓"), e.skill.ID, e.skill.Version)
		}
		installed++
	}

	fmt.Println()
	if installed == 0 && skipped == 0 {
		fmt.Println("No skills were selected.")
	} else {
		fmt.Printf("Done. %d skill(s) installed/updated, %d already up to date.\n", installed, skipped)
	}
	return nil
}
