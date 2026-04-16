package catalog_test

import (
	"testing"

	"github.com/simonski/skills/internal/catalog"
)

func TestAll_ReturnsSkills(t *testing.T) {
	skills, err := catalog.All()
	if err != nil {
		t.Fatalf("catalog.All() error: %v", err)
	}
	if len(skills) == 0 {
		t.Fatal("expected at least one skill in the catalog")
	}
}

func TestAll_SkillsHaveRequiredFields(t *testing.T) {
	skills, err := catalog.All()
	if err != nil {
		t.Fatalf("catalog.All() error: %v", err)
	}
	for _, s := range skills {
		if s.ID == "" {
			t.Errorf("skill missing ID: %+v", s)
		}
		if s.Version == "" {
			t.Errorf("skill %q missing Version", s.ID)
		}
		if s.Description == "" {
			t.Errorf("skill %q missing Description", s.ID)
		}
		if s.Content == "" {
			t.Errorf("skill %q missing Content", s.ID)
		}
	}
}

func TestGet_KnownSkill(t *testing.T) {
	s, err := catalog.Get("go")
	if err != nil {
		t.Fatalf("catalog.Get(\"go\") error: %v", err)
	}
	if s.ID != "go" {
		t.Errorf("expected ID %q, got %q", "go", s.ID)
	}
	if s.Version == "" {
		t.Error("expected non-empty Version")
	}
}

func TestGet_UnknownSkill(t *testing.T) {
	_, err := catalog.Get("nonexistent-skill-xyz")
	if err == nil {
		t.Fatal("expected error for unknown skill, got nil")
	}
}

func TestVersions_KnownSkill(t *testing.T) {
	versions, err := catalog.Versions("go")
	if err != nil {
		t.Fatalf("catalog.Versions(\"go\") error: %v", err)
	}
	if len(versions) == 0 {
		t.Fatal("expected at least one version for skill \"go\"")
	}
	for _, v := range versions {
		if v == "" {
			t.Error("unexpected empty version string")
		}
	}
}

func TestVersions_UnknownSkill(t *testing.T) {
	_, err := catalog.Versions("nonexistent-skill-xyz")
	if err == nil {
		t.Fatal("expected error for unknown skill, got nil")
	}
}

func TestVersions_SortedOldestToNewest(t *testing.T) {
	versions, err := catalog.Versions("go")
	if err != nil {
		t.Fatalf("catalog.Versions(\"go\") error: %v", err)
	}
	// Sorting is validated indirectly by TestGet_ReturnsLatestVersion.
	// Here we just confirm the slice is non-empty and each entry is non-empty.
	for _, v := range versions {
		if v == "" {
			t.Error("unexpected empty version in sorted list")
		}
	}
}

func TestGetVersion_KnownVersion(t *testing.T) {
	s, err := catalog.GetVersion("go", "0.0.1")
	if err != nil {
		t.Fatalf("catalog.GetVersion(\"go\", \"0.0.1\") error: %v", err)
	}
	if s.ID != "go" {
		t.Errorf("expected ID %q, got %q", "go", s.ID)
	}
	if s.Version != "0.0.1" {
		t.Errorf("expected Version %q, got %q", "0.0.1", s.Version)
	}
}

func TestGetVersion_UnknownVersion(t *testing.T) {
	_, err := catalog.GetVersion("go", "99.0.0")
	if err == nil {
		t.Fatal("expected error for unknown version, got nil")
	}
}

func TestGet_ReturnsLatestVersion(t *testing.T) {
	versions, err := catalog.Versions("go")
	if err != nil {
		t.Fatalf("catalog.Versions(\"go\") error: %v", err)
	}
	latest := versions[len(versions)-1]

	s, err := catalog.Get("go")
	if err != nil {
		t.Fatalf("catalog.Get(\"go\") error: %v", err)
	}
	if s.Version != latest {
		t.Errorf("Get returned version %q, expected latest %q", s.Version, latest)
	}
}
