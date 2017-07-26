package main

import (
	"os/exec"
)

// for gnome3 and unity desktop
func setWindowsWallPaper(fileName string) {
	cmd := exec.Command("gsettings",
		"set", "org.gnome.desktop.background", "picture-uri", "file://"+fileName)
	cmd.Run()
}
