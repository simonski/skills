package catalog

import (
	"bufio"
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

//go:embed skills/*.md
var skillFiles embed.FS

// Skill represents a skill in the catalog.
type Skill struct {
	ID          string
	Version     string
	Description string
	Content     string
}

// All returns all skills in the catalog.
func All() ([]*Skill, error) {
	var skills []*Skill

	entries, err := fs.ReadDir(skillFiles, "skills")
	if err != nil {
		return nil, fmt.Errorf("reading catalog: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}
		data, err := skillFiles.ReadFile(filepath.Join("skills", entry.Name()))
		if err != nil {
			return nil, fmt.Errorf("reading skill %s: %w", entry.Name(), err)
		}
		skill, err := parse(data)
		if err != nil {
			return nil, fmt.Errorf("parsing skill %s: %w", entry.Name(), err)
		}
		skills = append(skills, skill)
	}

	return skills, nil
}

// Get returns a single skill by its ID, or an error if not found.
func Get(id string) (*Skill, error) {
	skills, err := All()
	if err != nil {
		return nil, err
	}
	for _, s := range skills {
		if s.ID == id {
			return s, nil
		}
	}
	return nil, fmt.Errorf("skill %q not found in catalog", id)
}

// parse reads a skill file with YAML-like front matter (--- key: value ---).
func parse(data []byte) (*Skill, error) {
	s := &Skill{}

	scanner := bufio.NewScanner(bytes.NewReader(data))
	inFrontMatter := false
	bodyLines := []string{}
	frontMatterDone := false
	lineNum := 0

	for scanner.Scan() {
		line := scanner.Text()
		lineNum++

		if lineNum == 1 && line == "---" {
			inFrontMatter = true
			continue
		}

		if inFrontMatter {
			if line == "---" {
				inFrontMatter = false
				frontMatterDone = true
				continue
			}
			// parse key: value
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				switch key {
				case "id":
					s.ID = value
				case "version":
					s.Version = value
				case "description":
					s.Description = value
				}
			}
			continue
		}

		if frontMatterDone {
			bodyLines = append(bodyLines, line)
		}
	}

	s.Content = strings.TrimSpace(strings.Join(bodyLines, "\n"))
	return s, scanner.Err()
}
