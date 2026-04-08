package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/simonski/skills/internal/catalog"
)

var searchCmd = &cobra.Command{
	Use:   "search <term>",
	Short: "Search the skill catalog",
	Long: `Search the catalog for skills matching the given term.

The search is case-insensitive and matches against skill IDs and descriptions.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runSearch(args[0])
	},
}

func runSearch(term string) error {
	skills, err := catalog.All()
	if err != nil {
		return err
	}

	// Split the term into individual words for heuristic matching.
	words := strings.Fields(strings.ToLower(term))

	var matches []*catalog.Skill
	for _, s := range skills {
		haystack := strings.ToLower(s.ID + " " + s.Description + " " + s.Content)
		for _, word := range words {
			if strings.Contains(haystack, word) {
				matches = append(matches, s)
				break
			}
		}
	}

	if len(matches) == 0 {
		fmt.Printf("No skills found matching %q\n", term)
		return nil
	}

	fmt.Printf("%-20s %-12s %s\n", "ID", "VERSION", "DESCRIPTION")
	fmt.Printf("%-20s %-12s %s\n", "--------------------", "------------", "-----------")
	for _, s := range matches {
		fmt.Printf("%-20s %-12s %s\n", s.ID, s.Version, s.Description)
	}
	return nil
}
