package cli

import (
	"os"

	"golang.org/x/sys/windows"
)

func terminalColor(f *os.File) (colors, func()) {
	handle := windows.Handle(f.Fd())
	var mode uint32
	if windows.GetConsoleMode(handle, &mode) != nil {
		return false, func() {}
	}
	next := mode | windows.ENABLE_PROCESSED_OUTPUT | windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING
	if next == mode {
		return true, func() {}
	}
	if windows.SetConsoleMode(handle, next) != nil {
		return false, func() {}
	}
	return true, func() { _ = windows.SetConsoleMode(handle, mode) }
}
