package main

import (
	"flag"
	"sort"
	"strconv"
	"strings"

	"github.com/Alexandre-Chapelle/port-scanner/src/internal/output"
	"github.com/Alexandre-Chapelle/port-scanner/src/internal/scanner"
	"github.com/Alexandre-Chapelle/port-scanner/src/internal/ui"
)

func initFlags() (t string, pR string, th int, v bool, d bool, o string, of string, proc string) {
	var target string
	var portRange string
	var threads int
	var verbose bool
	var debug bool
	var outputFile string
	var outputFormat string
	var protocol string

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

	flag.Parse()

	return target, portRange, threads, verbose, debug, outputFile, outputFormat, protocol
}

func main() {
	target, portRange, threads, verbose, debug, outputFile, outputFormat, protocol := initFlags()

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

	ui.PrintfSuc("========================= OPEN PORTS ===========================")

	sort.Ints(resultPorts)
	for _, p := range resultPorts {
		ui.PrintfSuc("[+] %d", p)
	}

	if outputFile != "" {
		strP := make([]string, len(resultPorts))
		for i, v := range resultPorts {
			strP[i] = strconv.Itoa(v)
		}

		d := strings.Join(strP, "\n")

		output.OutputToFile(outputFile, d, outputFormat)
	}
}
