package cmd

import (
"strings"
"testing"

"github.com/simonski/skills/internal/catalog"
"github.com/simonski/skills/internal/project"
)

func latestGoVersion(t *testing.T) string {
t.Helper()
s, err := catalog.Get("go")
if err != nil {
t.Fatalf("catalog.Get(\"go\"): %v", err)
}
return s.Version
}

func TestRunUpdateAll_NoSkillsInstalled(t *testing.T) {
root := t.TempDir()
if err := runUpdateAll(root, false); err != nil {
t.Fatalf("runUpdateAll() error: %v", err)
}
}

func TestRunUpdateAll_AlreadyUpToDate(t *testing.T) {
root := t.TempDir()
latest := latestGoVersion(t)

content := "---\nid: go\nversion: " + latest + "\ndescription: Go best practices\n---\n\nContent.\n"
if err := project.Install(root, "go", content); err != nil {
t.Fatalf("Install() error: %v", err)
}

if err := runUpdateAll(root, true); err != nil {
t.Fatalf("runUpdateAll() error: %v", err)
}

ins, _ := project.Get(root, "go")
if ins.Version != latest {
t.Errorf("expected version %q, got %q", latest, ins.Version)
}
}

func TestRunUpdateAll_DryRun_DoesNotUpdate(t *testing.T) {
root := t.TempDir()

// Install an old version.
content := "---\nid: go\nversion: 0.0.1\ndescription: Go best practices\n---\n\nOld content.\n"
if err := project.Install(root, "go", content); err != nil {
t.Fatalf("Install() error: %v", err)
}

// Dry run — should NOT update.
if err := runUpdateAll(root, false); err != nil {
t.Fatalf("runUpdateAll(dry) error: %v", err)
}

ins, _ := project.Get(root, "go")
if ins.Version != "0.0.1" {
t.Errorf("dry run should not update: expected 0.0.1, got %q", ins.Version)
}
}

func TestRunUpdateAll_Apply_Updates(t *testing.T) {
root := t.TempDir()
latest := latestGoVersion(t)

content := "---\nid: go\nversion: 0.0.1\ndescription: Go best practices\n---\n\nOld content.\n"
if err := project.Install(root, "go", content); err != nil {
t.Fatalf("Install() error: %v", err)
}

if err := runUpdateAll(root, true); err != nil {
t.Fatalf("runUpdateAll(apply) error: %v", err)
}

ins, _ := project.Get(root, "go")
if ins.Version != latest {
t.Errorf("expected updated version %q, got %q", latest, ins.Version)
}
}

func TestRunUpdateOne_NotInstalled(t *testing.T) {
root := t.TempDir()
err := runUpdateOne(root, "go", false)
if err == nil {
t.Fatal("expected error for non-installed skill")
}
if !strings.Contains(err.Error(), "not installed") {
t.Errorf("unexpected error: %v", err)
}
}

func TestRunUpdateOne_AlreadyUpToDate(t *testing.T) {
root := t.TempDir()
latest := latestGoVersion(t)

content := "---\nid: go\nversion: " + latest + "\ndescription: Go best practices\n---\n\nContent.\n"
if err := project.Install(root, "go", content); err != nil {
t.Fatalf("Install() error: %v", err)
}

if err := runUpdateOne(root, "go", true); err != nil {
t.Fatalf("runUpdateOne() error: %v", err)
}

ins, _ := project.Get(root, "go")
if ins.Version != latest {
t.Errorf("expected version %q, got %q", latest, ins.Version)
}
}

func TestRunUpdateOne_DryRun_DoesNotUpdate(t *testing.T) {
root := t.TempDir()

content := "---\nid: go\nversion: 0.0.1\ndescription: Go best practices\n---\n\nOld content.\n"
if err := project.Install(root, "go", content); err != nil {
t.Fatalf("Install() error: %v", err)
}

if err := runUpdateOne(root, "go", false); err != nil {
t.Fatalf("runUpdateOne(dry) error: %v", err)
}

ins, _ := project.Get(root, "go")
if ins.Version != "0.0.1" {
t.Errorf("dry run should not update: expected 0.0.1, got %q", ins.Version)
}
}

func TestRunUpdateOne_Apply_Updates(t *testing.T) {
root := t.TempDir()
latest := latestGoVersion(t)

content := "---\nid: go\nversion: 0.0.1\ndescription: Go best practices\n---\n\nOld content.\n"
if err := project.Install(root, "go", content); err != nil {
t.Fatalf("Install() error: %v", err)
}

if err := runUpdateOne(root, "go", true); err != nil {
t.Fatalf("runUpdateOne(apply) error: %v", err)
}

ins, _ := project.Get(root, "go")
if ins.Version != latest {
t.Errorf("expected updated version %q, got %q", latest, ins.Version)
}
}

func TestRunUpdateOne_UnknownSkillID(t *testing.T) {
root := t.TempDir()
content := "---\nid: nonexistent-skill\nversion: 1.0.0\ndescription: Test\n---\n\nContent.\n"
if err := project.Install(root, "nonexistent-skill", content); err != nil {
t.Fatalf("Install() error: %v", err)
}

err := runUpdateOne(root, "nonexistent-skill", false)
if err == nil {
t.Fatal("expected error for skill not in catalog")
}
}
