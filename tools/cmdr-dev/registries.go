package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var screenRegisterRowPattern = regexp.MustCompile("^\\| \\x60([^\\x60]+)\\x60 \\|")

func loadActiveScreenIDs(root, specRel string) (map[string]struct{}, error) {
	path := filepath.Join(root, filepath.FromSlash(specRel), "00-governance", "registers", "screen-register.md")
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open screen register: %w", err)
	}
	defer f.Close()

	active := map[string]struct{}{}
	inActiveSection := false
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "## ") {
			switch {
			case strings.HasPrefix(trimmed, "## Écrans actifs"):
				inActiveSection = true
				continue
			case inActiveSection:
				inActiveSection = false
			}
		}
		if !inActiveSection {
			continue
		}
		match := screenRegisterRowPattern.FindStringSubmatch(trimmed)
		if len(match) == 2 {
			active[match[1]] = struct{}{}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read screen register: %w", err)
	}
	if len(active) == 0 {
		return nil, fmt.Errorf("screen register yielded no active screen ids")
	}
	return active, nil
}
