//go:build linux || darwin

package cli

import "os"

func terminalColor(_ *os.File) (colors, func()) { return true, func() {} }
