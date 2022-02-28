package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
)

var pkg = regexp.MustCompile("(.*) (.*)\n")

func main() {
  // TODO:  don't do this...
  // just do `asdf plugin list` and then `asdf current ${pkg}`
	home := os.Getenv("HOME")
	toolVersions := path.Join(home, ".tool-versions")

	contents, err := os.ReadFile(toolVersions)
	if err != nil {
		panic(fmt.Errorf("Failed to read tool versions:  %w", err))
	}

	var pkgErrors []string
	matches := pkg.FindAllStringSubmatch(string(contents), -1)
	for _, match := range matches {
		plugin := string(match[1])
		currentVersion := strings.TrimSpace(match[2])

    if currentVersion == "system" {
      fmt.Printf("Skipping system plugin %s\n", plugin)
      continue
    }

		latestVersion, err := exec.Command("asdf", "latest", plugin).Output()
		if err != nil {
			fmt.Println(fmt.Errorf("Failed to get latest for %s:  %w", plugin, err))
			pkgErrors = append(pkgErrors, plugin)
			continue
		}
		latest := strings.TrimSpace(string(latestVersion))
    if latest == currentVersion {
      fmt.Printf("Not updating %s (%s)\n", plugin, currentVersion)
      continue
    }

		var (
			install   = exec.Command("asdf", "install", plugin, latest)
			uninstall = exec.Command("asdf", "uninstall", plugin, string(currentVersion))
			global    = exec.Command("asdf", "global", plugin, latest)
		)

		fmt.Printf("Updating %s to %s\n", plugin, latest)

    install.Stdout = os.Stdout
    install.Stderr = os.Stderr
		err = install.Run()
		if err != nil {
			fmt.Println(fmt.Errorf("Failed to install:  %w", err))
			continue
		}
		err = uninstall.Run()
		if err != nil {
			fmt.Println(fmt.Errorf("Failed to uninstall:  %w", err))
			continue
		}
		err = global.Run()
		if err != nil {
			fmt.Println(fmt.Errorf("Failed to set global:  %w", err))
			continue
		}
	}

	if len(pkgErrors) != 0 {
    os.Exit(1)
	}
}
