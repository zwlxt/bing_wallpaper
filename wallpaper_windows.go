package main

/*
#include <windows.h>
void change_wallpaper(char path[255])
{
    SystemParametersInfo(0x0014, 0, path, 1);
}
*/
import "C"
import (
	"fmt"
	"strings"
	"unsafe"
)

func setWindowsWallPaper(path string) {
	path = strings.Replace(path, "/", "\\", -1)
	fmt.Println(path)

	cs := C.CString(path)
	defer C.free(unsafe.Pointer(cs))
	C.change_wallpaper(cs)
}
