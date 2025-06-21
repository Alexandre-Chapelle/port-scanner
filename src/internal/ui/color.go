package ui

import "fmt"

const (
	colorRed   = "\033[31m"
	colorGreen = "\033[32m"
	colorBlue  = "\033[34m"
	colorReset = "\033[0m"
)

func PrintfErr(format string, a ...any) {
	fmt.Printf("%s"+format+"%s\n", append([]any{colorRed}, append(a, colorReset)...)...)
}

func PrintfSuc(format string, a ...any) {
	fmt.Printf("%s"+format+"%s\n", append([]any{colorGreen}, append(a, colorReset)...)...)
}

func PrintfInfo(format string, a ...any) {
	fmt.Printf("%s"+format+"%s\n", append([]any{colorBlue}, append(a, colorReset)...)...)
}
