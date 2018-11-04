package main

import (
	"fmt"
	"strings"
	"syscall"
	"unsafe"
)

func setWindowsWallPaper(path string) {
	path = strings.Replace(path, "/", "\\", -1)
	fmt.Println(path)

	mod := syscall.NewLazyDLL("user32.dll")
	proc := mod.NewProc("SystemParametersInfoW")
	proc.Call(
		0x0014,
		0,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(path))),
		1,
	)
}
