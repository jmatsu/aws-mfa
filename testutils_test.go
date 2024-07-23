package main

import (
	"os"
	"path/filepath"
	"testing"
)

func makeDir(t *testing.T, d ...string) string {
	dir := filepath.Join(d...)

	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}

	return dir
}

func makeFile(t *testing.T, dir, name string, perm os.FileMode) string {
	path := filepath.Join(dir, name)

	f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, perm)

	if err != nil {
		t.Fatal(err)
	}

	if err := f.Close(); err != nil {
		t.Fatal(err)
	}

	return path
}

func chmod(t *testing.T, dir, name string, perm os.FileMode) {
	path := filepath.Join(dir, name)

	err := os.Chmod(path, perm)

	if err != nil {
		t.Fatal(err)
	}
}

func withTempDir(t *testing.T, r func(string)) {
	tmpDir, err := os.MkdirTemp("", "aws-mfa")

	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Error(err)
		}
	}()

	r(tmpDir)
}
