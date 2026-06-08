package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCanReadCertAndKey(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(certPath, keyPath string)
		wantOK    bool
		wantError bool
	}{
		{
			name:   "neither file exists",
			setup:  func(_, _ string) {},
			wantOK: false,
		},
		{
			name: "both files exist",
			setup: func(certPath, keyPath string) {
				writeFile(t, certPath, "cert")
				writeFile(t, keyPath, "key")
			},
			wantOK: true,
		},
		{
			name: "only cert exists",
			setup: func(certPath, _ string) {
				writeFile(t, certPath, "cert")
			},
			wantError: true,
		},
		{
			name: "only key exists",
			setup: func(_, keyPath string) {
				writeFile(t, keyPath, "key")
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()
			certPath := filepath.Join(dir, "cert.pem")
			keyPath := filepath.Join(dir, "key.pem")
			tt.setup(certPath, keyPath)

			ok, err := CanReadCertAndKey(certPath, keyPath)
			if tt.wantError {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if ok != tt.wantOK {
				t.Fatalf("CanReadCertAndKey() = %v, want %v", ok, tt.wantOK)
			}
		})
	}
}

func TestCanReadFile(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	existing := filepath.Join(dir, "exists.txt")
	writeFile(t, existing, "data")

	if !canReadFile(existing) {
		t.Fatal("canReadFile() = false, want true for existing file")
	}
	if canReadFile(filepath.Join(dir, "missing.txt")) {
		t.Fatal("canReadFile() = true, want false for missing file")
	}
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write file %s: %v", path, err)
	}
}
