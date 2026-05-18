package project

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type Workspace struct {
	Dir string
}

type Question struct {
	UserInput string    `json:"user_input"`
	Domain    string    `json:"domain"`
	Depth     string    `json:"depth"`
	CreatedAt time.Time `json:"created_at"`
}

func CreateWorkspace(root string, question string, domain string, depth string) (Workspace, error) {
	if root == "" {
		root = "researcher-workspace"
	}
	slug := slugify(question)
	if slug == "" {
		slug = "research-" + time.Now().Format("20060102-150405")
	}
	dir := uniqueDir(root, slug)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return Workspace{}, err
	}
	q := Question{UserInput: question, Domain: domain, Depth: depth, CreatedAt: time.Now()}
	if err := WriteJSON(filepath.Join(dir, "question.json"), q); err != nil {
		return Workspace{}, err
	}
	return Workspace{Dir: dir}, nil
}

func WriteJSON(path string, v any) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

func uniqueDir(root string, slug string) string {
	dir := filepath.Join(root, slug)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return dir
	}
	suffix := time.Now().Format("20060102-150405-000000000")
	return filepath.Join(root, slug+"-"+suffix)
}

func slugify(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	re := regexp.MustCompile(`[^a-z0-9]+`)
	value = re.ReplaceAllString(value, "-")
	value = strings.Trim(value, "-")
	if len(value) > 64 {
		value = strings.Trim(value[:64], "-")
	}
	return value
}
