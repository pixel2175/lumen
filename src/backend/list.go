package backend

import (
	"os"
	"path/filepath"
	"strings"
)

type Monitor struct {
	Kind string
	Name string
}

func List() []Monitor {
	var list []Monitor

	dirs, _ := filepath.Glob("/sys/class/backlight/*")
	for _, d := range dirs {
		name := filepath.Base(d)
		for _, c := range connectors() {
			if isInternal(c.name) {
				name = c.name
				break
			}
		}
		list = append(list, Monitor{Kind: "sysfs", Name: name})
	}

	for _, c := range connectors() {
		if !isInternal(c.name) {
			list = append(list, Monitor{Kind: "ddc", Name: c.name})
		}
	}
	return list
}

func connectors() []struct{ dir, name string } {
	var out []struct{ dir, name string }
	dirs, _ := filepath.Glob("/sys/class/drm/card*-*")
	for _, dir := range dirs {
		s, err := os.ReadFile(dir + "/status")
		if err != nil || strings.TrimSpace(string(s)) != "connected" {
			continue
		}
		n := filepath.Base(dir)
		out = append(out, struct{ dir, name string }{dir, n[strings.Index(n, "-")+1:]})
	}
	return out
}
