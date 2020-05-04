package handlers

import "fmt"

const (
	colorBlack   = "\033[30m"
	colorRed     = "\033[31m"
	colorGreen   = "\033[32m"
	colorYellow  = "\033[33m"
	colorBlue    = "\033[34m"
	colorMagenta = "\033[35m"
	colorCyan    = "\033[36m"
	colorReset   = "\033[0m"
)

var clientColors = []string{
	"\033[32m", // green
	"\033[33m", // yellow
	"\033[34m", // blue
	"\033[36m", // cyan
	"\033[35m", // magenta
}

func fillColor(text string, color string) string {
	return fmt.Sprintf("%s%s%s", color, text, colorReset)
}
