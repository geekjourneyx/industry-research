package project

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestCreateWorkspaceWritesQuestion(t *testing.T) {
	root := t.TempDir()
	ws, err := CreateWorkspace(root, "瑞幸咖啡 2026 年门店数目标是否可信？", "chain-brand", "standard")
	if err != nil {
		t.Fatalf("CreateWorkspace error: %v", err)
	}
	if ws.Dir == "" {
		t.Fatalf("workspace dir empty")
	}

	path := filepath.Join(ws.Dir, "question.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("question.json missing: %v", err)
	}

	var got Question
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("question.json invalid: %v", err)
	}
	if got.UserInput != "瑞幸咖啡 2026 年门店数目标是否可信？" {
		t.Fatalf("UserInput = %q", got.UserInput)
	}
	if got.Domain != "chain-brand" || got.Depth != "standard" {
		t.Fatalf("unexpected domain/depth: %#v", got)
	}
	if got.CreatedAt.IsZero() {
		t.Fatalf("CreatedAt is zero")
	}
}

func TestWriteJSONCreatesPrettyJSON(t *testing.T) {
	path := filepath.Join(t.TempDir(), "value.json")
	if err := WriteJSON(path, map[string]string{"status": "ok"}); err != nil {
		t.Fatalf("WriteJSON error: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read JSON: %v", err)
	}
	if string(data) != "{\n  \"status\": \"ok\"\n}\n" {
		t.Fatalf("unexpected JSON: %q", string(data))
	}
}

func TestSlugifyASCIIQuestion(t *testing.T) {
	got := slugify(" Luckin 2026 Store Count?! ")
	if got != "luckin-2026-store-count" {
		t.Fatalf("slugify = %q", got)
	}
}

func TestCreateWorkspaceDoesNotOverwriteExistingWorkspace(t *testing.T) {
	root := t.TempDir()
	first, err := CreateWorkspace(root, "Luckin 2026 Store Count", "chain-brand", "standard")
	if err != nil {
		t.Fatalf("first CreateWorkspace error: %v", err)
	}
	second, err := CreateWorkspace(root, "Luckin 2026 Store Count", "chain-brand", "standard")
	if err != nil {
		t.Fatalf("second CreateWorkspace error: %v", err)
	}
	if first.Dir == second.Dir {
		t.Fatalf("expected distinct workspace dirs, got %q", first.Dir)
	}
}
