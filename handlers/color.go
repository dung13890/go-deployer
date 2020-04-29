package handlers

import "fmt"

var (
	Colors = []string{
		"\033[32m", // green
		"\033[33m", // yellow
		"\033[36m", // cyan
		"\033[35m", // magenta
		"\033[31m", // red
		"\033[34m", // blue
	}
	ResetColor = "\033[0m"
)

const (
	ColorBlack  = "\u001b[30m"
	ColorRed    = "\u001b[31m"
	ColorGreen  = "\u001b[32m"
	ColorYellow = "\u001b[33m"
	ColorBlue   = "\u001b[34m"
	ColorReset  = "\u001b[0m"
)

func FillColor(text string, color string) string {
	return fmt.Sprintf("%s%s%s", color, text, ColorReset)
}
