package cli

import (
	"io"
	"os"

	"golang.org/x/term"
)

type colors bool

func colorSupport(w io.Writer) (colors, func()) {
	file, ok := w.(*os.File)
	if !ok || os.Getenv("NO_COLOR") != "" || os.Getenv("TERM") == "dumb" || !term.IsTerminal(int(file.Fd())) {
		return false, func() {}
	}
	return terminalColor(file)
}

func (c colors) label(code, text string) string {
	if c {
		return "\x1b[" + code + "m" + text + "\x1b[0m"
	}
	return text
}
