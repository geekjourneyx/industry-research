package validate

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateWorkspaceRequiresFiles(t *testing.T) {
	dir := t.TempDir()
	err := Workspace(dir)
	if err == nil {
		t.Fatalf("expected missing files error")
	}
	for _, name := range RequiredFiles {
		content := []byte(`{}`)
		if filepath.Ext(name) == ".md" {
			content = []byte("# Report\n")
		}
		if err := os.WriteFile(filepath.Join(dir, name), content, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if err := Workspace(dir); err != nil {
		t.Fatalf("expected valid workspace, got %v", err)
	}
}
