package version

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	githubReleaseAPI = "https://api.github.com/repos/simonski/skills/releases/latest"
)

// LatestRelease queries the GitHub API for the latest release tag.
func LatestRelease() (string, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(githubReleaseAPI)
	if err != nil {
		return "", fmt.Errorf("checking for updates: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("checking for updates: HTTP %d", resp.StatusCode)
	}

	var release struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", fmt.Errorf("parsing release response: %w", err)
	}
	return strings.TrimPrefix(release.TagName, "v"), nil
}

// IsNewerThan returns true if version a is newer than b.
func IsNewerThan(a, b string) bool {
	return semverGT(a, b)
}

// IsOutdated returns true when a newer release is available.
// It returns (false, "", nil) when the current version is "dev" (development builds).
func IsOutdated(current string) (bool, string, error) {
	if current == "dev" || current == "" {
		return false, "", nil
	}
	latest, err := LatestRelease()
	if err != nil {
		return false, "", err
	}
	if latest == "" {
		return false, "", nil
	}
	if latest != current && semverGT(latest, current) {
		return true, latest, nil
	}
	return false, latest, nil
}

// semverGT returns true if a > b using simple semver comparison.
func semverGT(a, b string) bool {
	aParts := splitSemver(a)
	bParts := splitSemver(b)

	for i := 0; i < 3; i++ {
		if aParts[i] > bParts[i] {
			return true
		}
		if aParts[i] < bParts[i] {
			return false
		}
	}
	return false
}

func splitSemver(v string) [3]int {
	v = strings.TrimPrefix(v, "v")
	parts := strings.SplitN(v, ".", 3)
	var result [3]int
	for i, p := range parts {
		if i >= 3 {
			break
		}
		n := 0
		for _, c := range p {
			if c < '0' || c > '9' {
				break
			}
			n = n*10 + int(c-'0')
		}
		result[i] = n
	}
	return result
}
