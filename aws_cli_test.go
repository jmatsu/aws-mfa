package main

import (
	"testing"
)

func Test_existAwsCli(t *testing.T) {
	withTempDir(t, func(dir string) {
		dir1 := makeDir(t, dir, "foo")

		t.Setenv("PATH", dir1)

		t.Run("executable aws is found", func(t *testing.T) {
			makeFile(t, dir1, "aws", 0755)

			if !existAwsCli() {
				t.Errorf("existAwsCli() = false, want true")
			}
		})

		t.Run("non-executable aws is found", func(t *testing.T) {
			chmod(t, dir1, "aws", 0644)

			if existAwsCli() {
				t.Errorf("existAwsCli() = true, want false")
			}
		})
	})
}
