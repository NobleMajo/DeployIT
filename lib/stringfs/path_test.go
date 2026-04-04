package stringfs

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetRelativePaths(t *testing.T) {
	t.Parallel()
	home, wd, err := GetRelativePaths()
	if err != nil {
		t.Fatal(err)
	}
	wantHome, e1 := os.UserHomeDir()
	wantWd, e2 := os.Getwd()
	if e1 != nil || e2 != nil {
		t.Fatalf("os.UserHomeDir/Getwd: %v %v", e1, e2)
	}
	if home != wantHome {
		t.Fatalf("homeDir got %q, want %q", home, wantHome)
	}
	if wd != wantWd {
		t.Fatalf("workDir got %q, want %q", wd, wantWd)
	}
}

func TestAbsolutePathRefFrom_table(t *testing.T) {
	workDir := filepath.Clean("/test/cwd")
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("UserHomeDir: %v", err)
	}

	tests := []struct {
		input    string
		expected string
	}{
		{"", workDir},
		{"~", homeDir},
		{"~/test", filepath.Join(homeDir, "test")},
		{"/absolute/path", "/absolute/path"},
		{"/absolute/path/../test/test/../../path", filepath.Clean("/absolute/path")},
		{"relative/path", filepath.Join(workDir, "relative", "path")},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			t.Parallel()
			s := test.input
			err := AbsolutePathRefFrom(homeDir, workDir, &s)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if s != test.expected {
				t.Fatalf("expected %q, got %q", test.expected, s)
			}
		})
	}
}

func TestAbsolutePathRef_table(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd: %v", err)
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("UserHomeDir: %v", err)
	}

	tests := []struct {
		input    string
		expected string
	}{
		{"", cwd},
		{"~", homeDir},
		{"~/test", filepath.Join(homeDir, "test")},
		{"/absolute/path", "/absolute/path"},
		{"/absolute/path/../test/test/../../path", filepath.Clean("/absolute/path")},
		{"relative/path", filepath.Join(cwd, "relative", "path")},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			t.Parallel()
			s := test.input
			err := AbsolutePathRef(&s)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if s != test.expected {
				t.Fatalf("expected %q, got %q", test.expected, s)
			}
		})
	}
}

func TestAbsolutePathFrom_relativeUnderWorkDir(t *testing.T) {
	t.Parallel()
	const homeDir = "/fake/home"
	const workDir = "/project/root"
	got, err := AbsolutePathFrom(homeDir, workDir, "src/pkg")
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(workDir, "src", "pkg")
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestAbsolutePathFrom_absoluteIgnoresWorkDir(t *testing.T) {
	t.Parallel()
	got, err := AbsolutePathFrom("/any/home", "/other/cwd", "/var/log/app")
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Clean("/var/log/app")
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestAbsolutePathFrom_tildeUsesParamHomeDir(t *testing.T) {
	t.Parallel()
	const fakeHome = "/fake/home"
	got, err := AbsolutePathFrom(fakeHome, "/w", "~/sub")
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(fakeHome, "sub")
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestAbsolutePath_absolute(t *testing.T) {
	t.Parallel()
	got, err := AbsolutePath("/var/log/myapp")
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Clean("/var/log/myapp")
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestAbsolutePath_relativeUsesChdir(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	got, err := AbsolutePath("nested/sub")
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(dir, "nested", "sub")
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestAbsolutePathRefFrom_nilPath(t *testing.T) {
	t.Parallel()
	err := AbsolutePathRefFrom("/h", "/w", nil)
	if err == nil || !strings.Contains(err.Error(), "path is nil") {
		t.Fatalf("got %v", err)
	}
}

func TestAbsolutePathRef_nilPath(t *testing.T) {
	t.Parallel()
	err := AbsolutePathRef(nil)
	if err == nil || !strings.Contains(err.Error(), "path is nil") {
		t.Fatalf("got %v", err)
	}
}
