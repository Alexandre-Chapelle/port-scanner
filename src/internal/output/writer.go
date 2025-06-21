package output

import (
	"os"

	"github.com/Alexandre-Chapelle/port-scanner/src/internal/ui"
)

func OutputToFile(fileName string, d string, of string) {
	err := os.WriteFile(fileName, []byte(d), 0755)

	if err != nil {
		ui.PrintfErr("[-] Cannot write to file %s", fileName)
	}
}
