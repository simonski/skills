package project_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/simonski/skills/internal/project"
)

func TestInstallAndGet(t *testing.T) {
	root := t.TempDir()
	content := "---\nid: test-skill\nversion: 1.2.3\ndescription: A test skill\n---\n\nContent here.\n"

	if err := project.Install(root, "test-skill", content); err != nil {
		t.Fatalf("Install() error: %v", err)
	}

	is, err := project.Get(root, "test-skill")
	if err != nil {
		t.Fatalf("Get() error: %v", err)
	}
	if is == nil {
		t.Fatal("Get() returned nil, expected installed skill")
	}
	if is.ID != "test-skill" {
		t.Errorf("expected ID %q, got %q", "test-skill", is.ID)
	}
	if is.Version != "1.2.3" {
		t.Errorf("expected version %q, got %q", "1.2.3", is.Version)
	}
}

func TestGet_NotInstalled(t *testing.T) {
	root := t.TempDir()
	is, err := project.Get(root, "missing-skill")
	if err != nil {
		t.Fatalf("Get() unexpected error: %v", err)
	}
	if is != nil {
		t.Errorf("expected nil for uninstalled skill, got %+v", is)
	}
}

func TestList_Empty(t *testing.T) {
	root := t.TempDir()
	skills, err := project.List(root)
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}
	if len(skills) != 0 {
		t.Errorf("expected empty list, got %v", skills)
	}
}

func TestList_WithSkills(t *testing.T) {
	root := t.TempDir()
	for _, id := range []string{"skill-a", "skill-b"} {
		content := "---\nid: " + id + "\nversion: 1.0.0\ndescription: Test\n---\n\nContent.\n"
		if err := project.Install(root, id, content); err != nil {
			t.Fatalf("Install(%q) error: %v", id, err)
		}
	}

	skills, err := project.List(root)
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}
	if len(skills) != 2 {
		t.Errorf("expected 2 skills, got %d", len(skills))
	}
}

func TestRemove(t *testing.T) {
	root := t.TempDir()
	content := "---\nid: rm-skill\nversion: 1.0.0\ndescription: Test\n---\n\nContent.\n"
	if err := project.Install(root, "rm-skill", content); err != nil {
		t.Fatalf("Install() error: %v", err)
	}

	if err := project.Remove(root, "rm-skill"); err != nil {
		t.Fatalf("Remove() error: %v", err)
	}

	// File should be gone.
	path := filepath.Join(root, ".skills", "rm-skill.md")
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Errorf("expected file to be deleted, but it still exists at %s", path)
	}
}

func TestRemove_NotInstalled(t *testing.T) {
	root := t.TempDir()
	err := project.Remove(root, "ghost-skill")
	if err == nil {
		t.Fatal("expected error when removing non-existent skill")
	}
}
