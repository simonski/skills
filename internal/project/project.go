package project

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const skillsDir = ".skills"

// InstalledSkill represents a skill installed in the current project.
type InstalledSkill struct {
	ID      string
	Version string
}

// Dir returns the path to the skills directory for the given project root.
func Dir(root string) string {
	return filepath.Join(root, skillsDir)
}

// SkillPath returns the path to the installed skill file.
func SkillPath(root, id string) string {
	return filepath.Join(Dir(root), id+".md")
}

// List returns all skills installed in the project at root.
func List(root string) ([]*InstalledSkill, error) {
	dir := Dir(root)
	entries, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("reading skills directory: %w", err)
	}

	var skills []*InstalledSkill
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}
		path := filepath.Join(dir, entry.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("reading %s: %w", path, err)
		}
		is := parseInstalled(data)
		if is.ID == "" {
			// Fall back to filename-derived ID
			is.ID = strings.TrimSuffix(entry.Name(), ".md")
		}
		skills = append(skills, is)
	}
	return skills, nil
}

// Get returns a single installed skill by ID, or nil if not installed.
func Get(root, id string) (*InstalledSkill, error) {
	path := SkillPath(root, id)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", path, err)
	}
	is := parseInstalled(data)
	if is.ID == "" {
		is.ID = id
	}
	return is, nil
}

// Install writes a skill file into the project.
func Install(root, id, content string) error {
	dir := Dir(root)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("creating skills directory: %w", err)
	}
	path := SkillPath(root, id)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return fmt.Errorf("writing skill %s: %w", id, err)
	}
	return nil
}

// Remove deletes a skill file from the project.
func Remove(root, id string) error {
	path := SkillPath(root, id)
	if err := os.Remove(path); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("skill %q is not installed", id)
		}
		return fmt.Errorf("removing skill %s: %w", id, err)
	}
	return nil
}

// parseInstalled extracts id and version from the front matter of a skill file.
func parseInstalled(data []byte) *InstalledSkill {
	is := &InstalledSkill{}
	scanner := bufio.NewScanner(bytes.NewReader(data))
	inFrontMatter := false
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
				break
			}
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				switch key {
				case "id":
					is.ID = value
				case "version":
					is.Version = value
				}
			}
		}
	}
	return is
}
