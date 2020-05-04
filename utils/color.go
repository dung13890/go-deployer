package utils

import "fmt"

const (
	ColorBlack   = "\033[30m"
	ColorRed     = "\033[31m"
	ColorGreen   = "\033[32m"
	ColorYellow  = "\033[33m"
	ColorBlue    = "\033[34m"
	ColorMagenta = "\033[35m"
	ColorCyan    = "\033[36m"
	ColorReset   = "\033[0m"
)

var ClientColors = []string{
	"\033[36m", // cyan
	"\033[35m", // magenta
	"\033[34m", // blue
	"\033[33m", // yellow
	"\033[32m", // green
}

func FillColor(text string, color string) string {
	return fmt.Sprintf("%s%s%s", color, text, ColorReset)
}
