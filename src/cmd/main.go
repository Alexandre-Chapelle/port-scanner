package main

import (
	"flag"
	"os"
	"sort"
	"time"

	"github.com/Alexandre-Chapelle/port-scanner/src/internal/output"
	"github.com/Alexandre-Chapelle/port-scanner/src/internal/scanner"
	"github.com/Alexandre-Chapelle/port-scanner/src/internal/ui"
)

func printHelp() {
	ui.PrintfInfo("Port Scanner - Network Port Discovery Tool")
	ui.PrintfInfo("==============================================")
	ui.PrintfInfo("")
	ui.PrintfInfo("USAGE:")
	ui.PrintfInfo("  ./port-scanner [OPTIONS]")
	ui.PrintfInfo("")
	ui.PrintfInfo("TARGET OPTIONS:")
	ui.PrintfInfo("  -target, -t <url>           Target host to scan (default: 127.0.0.1)")
	ui.PrintfInfo("  -port-range, -p <range>     Port range to scan (default: 1-65535)")
	ui.PrintfInfo("                              Examples: 80, 80-443, 1-1000")
	ui.PrintfInfo("")
	ui.PrintfInfo("SCAN OPTIONS:")
	ui.PrintfInfo("  -threads, -ts <number>      Number of concurrent threads (default: 100)")
	ui.PrintfInfo("  -protocol, -proc <proto>    Protocol to use (default: tcp)")
	ui.PrintfInfo("                              Options: tcp, tcp4, tcp6, udp, udp4, udp6")
	ui.PrintfInfo("")
	ui.PrintfInfo("OUTPUT OPTIONS:")
	ui.PrintfInfo("  -verbose, -v                Enable verbose output")
	ui.PrintfInfo("  -debug, -d                  Enable debug output")
	ui.PrintfInfo("  -output, -o <file>          Save results to file")
	ui.PrintfInfo("  -output-format, -of <fmt>   Output format (PLAIN | HTML)")
	ui.PrintfInfo("")
	ui.PrintfInfo("HELP:")
	ui.PrintfInfo("  -help, -h                   Show this help message")
	ui.PrintfInfo("")
	ui.PrintfInfo("EXAMPLES:")
	ui.PrintfInfo("  ./port-scanner -t 192.168.1.1 -p 1-1000 -v")
	ui.PrintfInfo("  ./port-scanner --target google.com --port-range 80,443 --threads 50")
	ui.PrintfInfo("  ./port-scanner -t 10.0.0.1 -p 20-25 -o results.html -of HTML")
	ui.PrintfInfo("")
}

func initFlags() (t string, pR string, th int, v bool, d bool, o string, of string, proc string) {
	var target string
	var portRange string
	var threads int
	var verbose bool
	var debug bool
	var outputFile string
	var outputFormat string
	var protocol string
	var help bool

	flag.StringVar(&target, "target", "127.0.0.1", "[--target <url>] Specifies the target to scan")
	flag.StringVar(&target, "t", "127.0.0.1", "[-t <url>] Specifies the target to scan (short)")

	flag.StringVar(&portRange, "port-range", "1-65535", "[--port-range 0-65535] Specifies the port range")
	flag.StringVar(&portRange, "p", "1-65535", "[-p 0-65535] Specifies the port range (short)")

	flag.IntVar(&threads, "threads", 100, "[--threads 50] Specifies amount of threads")
	flag.IntVar(&threads, "ts", 100, "[-ts 50] Specifies amount of threads (short)")

	flag.BoolVar(&verbose, "verbose", false, "[--verbose] Specifies if you want a verbose output")
	flag.BoolVar(&verbose, "v", false, "[-v] Specifies if you want a verbose output (short)")

	flag.BoolVar(&debug, "debug", false, "[--debug] Specifies if you want a debug output")
	flag.BoolVar(&debug, "d", false, "[-d] Specifies if you want a debug output (short)")

	flag.StringVar(&outputFile, "output", "", "[--output <file>] Specifies where you want to save the results")
	flag.StringVar(&outputFile, "o", "", "[-o <file>] Specifies where you want to save the results (short)")

	flag.StringVar(&outputFormat, "output-format", "", "[--output-format <PLAIN | HTML>] Specifies in which format you want to save the results")
	flag.StringVar(&outputFormat, "of", "", "[-o <PLAIN | HTML>] Specifies in which format you want to save the results (short)")

	flag.StringVar(&protocol, "protocol", "tcp", "[--protocol <tcp | tcp4 | tcp6 | udp | udp4 | udp6 | ip | ip4 | ip6 | unix | unixgram and unixpacket>] Specifies the protocol to use")
	flag.StringVar(&protocol, "proc", "tcp", "[-proc <tcp | tcp4 | tcp6 | udp | udp4 | udp6 | ip | ip4 | ip6 | unix | unixgram and unixpacket>] Specifies the protocol to use (short)")

	flag.BoolVar(&help, "help", false, "[--help] Get help")
	flag.BoolVar(&help, "h", false, "[-h] Get help (short)")

	flag.Parse()

	if help {
		printHelp()
		os.Exit(0)
	}

	return target, portRange, threads, verbose, debug, outputFile, outputFormat, protocol
}

func main() {
	start := time.Now()
	target, portRange, threads, verbose, debug, outputFile, outputFormat, protocol := initFlags()

	ui.PrintfInfo("[~] Started scan for target %s on ports %s", target, portRange)

	scanner := scanner.New(
		scanner.Config{
			Target:       target,
			PortRange:    portRange,
			Threads:      threads,
			Verbose:      verbose,
			Debug:        debug,
			OutputFile:   outputFile,
			OutputFormat: outputFormat,
			Protocol:     protocol,
		},
	)

	resultPorts, err := scanner.Scan()

	if err != nil {
		ui.PrintfErr("[-] Cannot retrieve result ports")

		if debug {
			ui.PrintfErr(err.Error())
		}
	}

	tElapsed := time.Since(start)

	if outputFile != "" {

		switch outputFormat {
		case "HTML":
			html, err := output.FormatHTML(resultPorts, target, portRange, tElapsed)
			if err != nil {
				ui.PrintfErr("[-] Error while creating HTML template")
			}
			output.OutputToFile(outputFile, html, "html")
			ui.PrintfInfo("[~] Generated HTML report in %s file", outputFile)

		case "PLAIN":
			plain := output.FormatPlain(resultPorts, target, tElapsed)
			output.OutputToFile(outputFile, plain, "txt")
			ui.PrintfInfo("[~] Generated PLAIN report in %s file", outputFile)

		default:
			ui.PrintfSuc("\n\n========================= OPEN PORTS ===========================\n")
			sort.Ints(resultPorts)
			for _, p := range resultPorts {
				ui.PrintfSuc("[+] %d", p)
			}
		}

	}
}
