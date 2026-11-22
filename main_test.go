package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func write(t *testing.T, dir, name, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestAssembleSortedFragments(t *testing.T) {
	dir := t.TempDir()
	write(t, dir, "50-hardening.conf", "sysctl.kernel.dmesg_restrict=1 loglevel=0\n")
	write(t, dir, "00-base.conf", "rootflags=subvol=@\ninit=/usr/libexec/dinit/init\n")
	out, err := Assemble(dir)
	if err != nil {
		t.Fatal(err)
	}
	want := "rootflags=subvol=@ init=/usr/libexec/dinit/init sysctl.kernel.dmesg_restrict=1 loglevel=0"
	if out != want {
		t.Fatalf("got %q, want %q", out, want)
	}
}

func TestAssembleDropsCommentsAndBlanks(t *testing.T) {
	dir := t.TempDir()
	write(t, dir, "a.conf", "# comment\n   # indented comment\n\nrootflags=subvol=@\n")
	write(t, dir, "b.conf", "\n\n")
	out, err := Assemble(dir)
	if err != nil {
		t.Fatal(err)
	}
	if out != "rootflags=subvol=@" {
		t.Fatalf("got %q", out)
	}
}

func TestAssemblePreservesInternalSpacing(t *testing.T) {
	dir := t.TempDir()
	write(t, dir, "94.conf", "kernel.printk=3 3 3 3\n")
	out, err := Assemble(dir)
	if err != nil {
		t.Fatal(err)
	}
	if out != "kernel.printk=3 3 3 3" {
		t.Fatalf("internal spacing lost: %q", out)
	}
}

func TestAssembleNonConfFilesIgnored(t *testing.T) {
	dir := t.TempDir()
	write(t, dir, "README", "rootflags=WRONG\n")
	write(t, dir, "x.conf.bak", "rootflags=WRONG\n")
	write(t, dir, "ok.conf", "quiet=1\n")
	out, err := Assemble(dir)
	if err != nil {
		t.Fatal(err)
	}
	if out != "quiet=1" {
		t.Fatalf("got %q", out)
	}
}

func TestAssembleMissingDir(t *testing.T) {
	if _, err := Assemble(filepath.Join(t.TempDir(), "nope")); err == nil {
		t.Fatal("expected error for missing dir")
	}
}

func TestAssembleEmptyDir(t *testing.T) {
	out, err := Assemble(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(out) != "" {
		t.Fatalf("expected empty, got %q", out)
	}
}
