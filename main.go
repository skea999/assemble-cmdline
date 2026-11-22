// assemble-cmdline assembles the kernel command line from fragment files
// in a directory (/etc/kernel/cmdline.d by default):
//
//   - files are read in lexicographic order
//   - lines whose first non-blank character is '#' are dropped
//   - remaining lines are joined with single spaces
//
// The result goes to stdout. Writing /etc/kernel/cmdline and rebuilding the
// UKI is the caller's job (installer, kernel hook, apk trigger) — this tool
// is pure: it never touches the filesystem.
//
// FeralOS packages ship fragments (e.g. hardening-config ships
// 50-hardening.conf); upgrading such a package changes the fragment, and
// the next kernel hook / trigger run picks it up.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const defaultDir = "/etc/kernel/cmdline.d"

// usage: assemble-cmdline [dir]
func main() {
	dir := defaultDir
	if len(os.Args) > 2 {
		fmt.Fprintln(os.Stderr, "usage: assemble-cmdline [dir]")
		os.Exit(2)
	}
	if len(os.Args) == 2 {
		dir = os.Args[1]
	}
	out, err := Assemble(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "assemble-cmdline: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(out)
}

// Assemble joins the non-comment lines of the sorted *.conf files in dir
// with single spaces. Internal spacing of a line is preserved (e.g.
// "kernel.printk=3 3 3 3").
func Assemble(dir string) (string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", fmt.Errorf("read %s: %w", dir, err)
	}
	// ReadDir sorts by name; keep only *.conf fragments.
	var files []string
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".conf" {
			continue
		}
		files = append(files, filepath.Join(dir, e.Name()))
	}
	sort.Strings(files)

	var lines []string
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			return "", fmt.Errorf("read %s: %w", f, err)
		}
		for _, line := range strings.Split(string(data), "\n") {
			if strings.TrimSpace(strings.TrimLeft(line, " \t")) == "" {
				continue
			}
			if strings.HasPrefix(strings.TrimLeft(line, " \t"), "#") {
				continue
			}
			lines = append(lines, strings.TrimSpace(line))
		}
	}
	return strings.Join(lines, " "), nil
}
