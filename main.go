package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var currentPattern = regexp.MustCompile("^.*\\s+([\\d\\w\\.\\-]+).*$")

func main() {
	pluginsOutput, err := exec.Command("asdf", "plugin", "list").Output()
	if err != nil {
		panic(fmt.Errorf("Failed to list `asdf` plugins:  %w", err))
	}

	plugins := strings.TrimSpace(string(pluginsOutput))

	var pkgErrors []string

	for _, plugin := range strings.Split(plugins, "\n") {
		var currentVersion string

		currentOutput, err := exec.Command("asdf", "current", plugin).Output()
		if err != nil {
			fmt.Printf("No current version for %s\n", pluginFmt(plugin))
		} else {
			trimmedOutput := strings.TrimSpace(string(currentOutput))
			matches := currentPattern.FindStringSubmatch(trimmedOutput)
			if len(matches) > 1 {
				currentVersion = matches[1]
			}
		}

		if currentVersion == "system" {
			fmt.Printf("Skipping system plugin %s\n", pluginFmt(plugin))
			continue
		}

		latestVersion, err := exec.Command("asdf", "latest", plugin).Output()
		if err != nil {
			fmt.Println(fmt.Errorf("Failed to get latest for %s:  %w", pluginFmt(plugin), err))
			pkgErrors = append(pkgErrors, plugin)
			continue
		}
		latest := strings.TrimSpace(string(latestVersion))
		if latest == currentVersion {
			fmt.Printf("Not updating %s (%s)\n", pluginFmt(plugin), versionFmt(currentVersion))
			continue
		}

		fmt.Printf("Installing version %s of %s\n", versionFmt(latest), pluginFmt(plugin))

		var (
			install   = exec.Command("asdf", "install", plugin, latest)
			uninstall = exec.Command("asdf", "uninstall", plugin, currentVersion)
			global    = exec.Command("asdf", "global", plugin, latest)
		)

		install.Stdout = os.Stdout
		err = install.Run()
		if err != nil {
			fmt.Println(fmt.Errorf("Failed to install:  %w", err))
			continue
		}
		if currentVersion != "" {
			fmt.Printf("Uninstalling version %s\n", versionFmt(currentVersion))
			err = uninstall.Run()
			if err != nil {
				fmt.Println(fmt.Errorf("Failed to uninstall:  %w", err))
				continue
			}
		}
		fmt.Printf("Setting global version to %s\n", versionFmt(latest))
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

func versionFmt(vsn string) string {
	return fmt.Sprintf("\033[30;32m%s\033[0m", vsn)
}

func pluginFmt(plug string) string {
	return fmt.Sprintf("\033[30;34m%s\033[0m", plug)
}
