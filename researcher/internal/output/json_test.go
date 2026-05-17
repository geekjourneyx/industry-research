package output

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestWriteJSON(t *testing.T) {
	var buf bytes.Buffer
	err := WriteJSON(&buf, map[string]string{"status": "ok"}, false)
	if err != nil {
		t.Fatalf("WriteJSON returned error: %v", err)
	}
	var decoded map[string]string
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if decoded["status"] != "ok" {
		t.Fatalf("expected status ok, got %q", decoded["status"])
	}
}
