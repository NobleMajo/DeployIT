package bun

import (
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

//go:generate go run ./interface

//go:embed embed/interface*
var embeddedBunFS embed.FS

const MinEmbeddedBunBytes = 1024

func EmbeddedBunBytes() ([]byte, error) {
	entries, err := fs.ReadDir(embeddedBunFS, "embed")
	if err != nil {
		return nil, fmt.Errorf("read embedded Bun CLI directory: %w", err)
	}

	var name string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		n := e.Name()
		if !strings.HasPrefix(n, "interface") {
			continue
		}
		if name != "" {
			return nil, fmt.Errorf("multiple interface embeds under embed/: %q and %q", name, n)
		}
		name = n
	}

	if name == "" {
		return nil, fmt.Errorf("no embedded Bun binary; run go generate ./lib/bun")
	}

	bytes, err := embeddedBunFS.ReadFile(filepath.Join("embed", name))
	if err != nil {
		return nil, fmt.Errorf("read embedded Bun CLI: %w", err)
	}

	if len(bytes) < MinEmbeddedBunBytes {
		return nil, fmt.Errorf(
			"embedded Bun CLI missing or too small (%d bytes); run: go generate ./lib/bun",
			len(bytes),
		)
	}

	return bytes, nil
}
