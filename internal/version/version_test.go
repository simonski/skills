package version_test

import (
	"testing"

	"github.com/simonski/skills/internal/version"
)

func TestIsNewerThan(t *testing.T) {
	tests := []struct {
		a, b string
		want bool
	}{
		{"1.1.0", "1.0.0", true},
		{"2.0.0", "1.9.9", true},
		{"1.0.1", "1.0.0", true},
		{"1.0.0", "1.0.0", false},
		{"0.9.0", "1.0.0", false},
		{"1.0.0", "1.0.1", false},
	}

	for _, tc := range tests {
		got := version.IsNewerThan(tc.a, tc.b)
		if got != tc.want {
			t.Errorf("IsNewerThan(%q, %q) = %v, want %v", tc.a, tc.b, got, tc.want)
		}
	}
}

func TestIsOutdated_DevVersion(t *testing.T) {
	outdated, latest, err := version.IsOutdated("dev")
	if err != nil {
		t.Fatalf("IsOutdated(\"dev\") unexpected error: %v", err)
	}
	if outdated {
		t.Error("dev builds should never be considered outdated")
	}
	if latest != "" {
		t.Errorf("expected empty latest for dev, got %q", latest)
	}
}
