//go:build ignore

// Standalone CI builder for the Batkixni/raglex-plugins community PR repository.
// Usage: go run ./scripts/build_registry.go -verify-download -out registry.json
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

var validID = regexp.MustCompile(`^[a-z0-9][a-z0-9._-]{1,63}$`)

type Entry struct {
	ID               string   `json:"id"`
	Name             string   `json:"name"`
	Version          string   `json:"version"`
	Author           string   `json:"author"`
	Description      string   `json:"description"`
	Icon             string   `json:"icon,omitempty"`
	Accent           string   `json:"accent,omitempty"`
	Category         string   `json:"category,omitempty"`
	Repository       string   `json:"repository,omitempty"`
	DownloadURL      string   `json:"download_url,omitempty"`
	SHA256           string   `json:"sha256,omitempty"`
	MinRaglexVersion string   `json:"min_raglex_version,omitempty"`
	Hooks            []string `json:"hooks,omitempty"`
	Permissions      []string `json:"permissions,omitempty"`
	Official         bool     `json:"official,omitempty"`
	Builtin          bool     `json:"builtin,omitempty"`
}

type Catalog struct {
	SchemaVersion int     `json:"schema_version"`
	UpdatedAt     string  `json:"updated_at"`
	Plugins       []Entry `json:"plugins"`
}

func main() {
	pluginsDir := flag.String("dir", "plugins", "directory containing <id>/plugin.json")
	outPath := flag.String("out", "registry.json", "output registry json path")
	verifyDL := flag.Bool("verify-download", false, "download every zip and verify sha256")
	flag.Parse()

	matches, err := filepath.Glob(filepath.Join(*pluginsDir, "*", "plugin.json"))
	if err != nil {
		fatalf("glob: %v", err)
	}
	seen := map[string]bool{}
	var entries []Entry
	for _, p := range matches {
		raw, err := os.ReadFile(p)
		if err != nil {
			fatalf("read %s: %v", p, err)
		}
		var e Entry
		if err := json.Unmarshal(raw, &e); err != nil {
			fatalf("parse %s: %v", p, err)
		}
		if !validID.MatchString(e.ID) {
			fatalf("%s: invalid id %q", p, e.ID)
		}
		if seen[e.ID] {
			fatalf("%s: duplicate plugin id %q", p, e.ID)
		}
		seen[e.ID] = true
		if strings.TrimSpace(e.Name) == "" || strings.TrimSpace(e.Version) == "" {
			fatalf("%s: name and version are required", p)
		}
		if !e.Builtin {
			if !strings.HasPrefix(e.DownloadURL, "https://") {
				fatalf("%s: download_url must use https://", p)
			}
			if len(strings.TrimSpace(e.SHA256)) != 64 {
				fatalf("%s: sha256 must be 64 hex chars", p)
			}
			if *verifyDL {
				resp, err := http.Get(e.DownloadURL)
				if err != nil || resp.StatusCode != 200 {
					fatalf("%s: download failed (%v)", p, err)
				}
				h := sha256.New()
				_, _ = io.Copy(h, resp.Body)
				_ = resp.Body.Close()
				got := hex.EncodeToString(h.Sum(nil))
				if !strings.EqualFold(got, e.SHA256) {
					fatalf("%s: sha256 mismatch: got %s, want %s", p, got, e.SHA256)
				}
			}
		}
		entries = append(entries, e)
	}

	sort.Slice(entries, func(i, j int) bool { return entries[i].ID < entries[j].ID })
	cat := Catalog{
		SchemaVersion: 1,
		UpdatedAt:     time.Now().UTC().Format(time.RFC3339),
		Plugins:       entries,
	}
	out, _ := json.MarshalIndent(cat, "", "  ")
	if err := os.WriteFile(*outPath, append(out, '\n'), 0o644); err != nil {
		fatalf("write %s: %v", *outPath, err)
	}
	fmt.Printf("Validated %d plugins -> %s\n", len(entries), *outPath)
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "error: "+format+"\n", args...)
	os.Exit(1)
}
