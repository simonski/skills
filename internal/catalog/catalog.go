package catalog

import (
	"bufio"
	"bytes"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

//go:embed skills/*/*.md
var skillFiles embed.FS

// Skill represents a skill in the catalog.
type Skill struct {
	ID          string
	Version     string
	Description string
	Content     string
}

// All returns the latest version of each skill in the catalog.
func All() ([]*Skill, error) {
	ids, err := skillIDs()
	if err != nil {
		return nil, err
	}

	var skills []*Skill
	for _, id := range ids {
		s, err := Get(id)
		if err != nil {
			return nil, err
		}
		skills = append(skills, s)
	}
	return skills, nil
}

// Get returns the latest version of a skill by its ID.
func Get(id string) (*Skill, error) {
	versions, err := Versions(id)
	if err != nil {
		return nil, err
	}
	if len(versions) == 0 {
		return nil, fmt.Errorf("skill %q has no versions in catalog", id)
	}
	return GetVersion(id, versions[len(versions)-1])
}

// GetVersion returns a specific version of a skill by ID and version string.
func GetVersion(id, version string) (*Skill, error) {
	path := filepath.Join("skills", id, version+".md")
	data, err := skillFiles.ReadFile(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, fmt.Errorf("skill %q version %q not found in catalog", id, version)
		}
		return nil, fmt.Errorf("reading skill %s@%s: %w", id, version, err)
	}

	skill, err := parse(data)
	if err != nil {
		return nil, fmt.Errorf("parsing skill %s@%s: %w", id, version, err)
	}
	if skill.ID == "" {
		skill.ID = id
	}
	if skill.Version == "" {
		skill.Version = version
	}
	return skill, nil
}

// Versions returns all available versions of a skill, sorted oldest to newest.
func Versions(id string) ([]string, error) {
	dir := filepath.Join("skills", id)
	entries, err := fs.ReadDir(skillFiles, dir)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, fmt.Errorf("skill %q not found in catalog", id)
		}
		return nil, fmt.Errorf("reading skill %s: %w", id, err)
	}

	var versions []string
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}
		if !isVersionFile(entry.Name()) {
			continue
		}
		versions = append(versions, strings.TrimSuffix(entry.Name(), ".md"))
	}

	sort.Slice(versions, func(i, j int) bool {
		return compareVersions(versions[i], versions[j]) < 0
	})
	return versions, nil
}

// skillIDs returns the IDs of all skills in the catalog (subdirectory names).
func skillIDs() ([]string, error) {
	entries, err := fs.ReadDir(skillFiles, "skills")
	if err != nil {
		return nil, fmt.Errorf("reading catalog: %w", err)
	}
	var ids []string
	for _, entry := range entries {
		if entry.IsDir() {
			ids = append(ids, entry.Name())
		}
	}
	return ids, nil
}

// isVersionFile returns true if the filename looks like a version file (e.g., "1.0.0.md").
func isVersionFile(name string) bool {
	return len(name) > 0 && name[0] >= '0' && name[0] <= '9'
}

// compareVersions compares two semver strings (e.g., "1.0.0" vs "1.1.0").
// Returns negative if a < b, zero if equal, positive if a > b.
func compareVersions(a, b string) int {
	ap := parseVersion(a)
	bp := parseVersion(b)
	for i := range ap {
		if ap[i] != bp[i] {
			return ap[i] - bp[i]
		}
	}
	return 0
}

// parseVersion splits a semver string into its numeric components.
func parseVersion(v string) [3]int {
	parts := strings.SplitN(v, ".", 3)
	var result [3]int
	for i, p := range parts {
		n, _ := strconv.Atoi(p)
		result[i] = n
	}
	return result
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
