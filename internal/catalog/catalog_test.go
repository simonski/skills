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
