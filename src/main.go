package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

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

func worker(v bool, d bool, proc string, t string, pp chan int, ports chan int) {
	for p := range ports {

		address := net.JoinHostPort(t, fmt.Sprintf("%d", p))
		conn, err := net.Dial(proc, address)

		if err != nil {
			if d {
				PrintfErr(err.Error())
			}

			pp <- 0

			if v {
				PrintfErr("[-] Port %d is closed", p)
			}

			continue
		}

		if v {
			PrintfSuc("[+] Port %d is open", p)
		}

		pp <- p

		conn.Close()
	}
}

func poolPorts(ports chan int, sPort int, ePort int, d bool) {
	for i := sPort; i <= ePort; i++ {
		ports <- i

		if d {
			PrintfInfo("[~] Added port %d to queue", i)
		}
	}
}

func parsePortsArg(portsArg string) (start int, end int) {
	re := regexp.MustCompile(`\d+-\d+`)
	unformattedPorts := re.FindString(portsArg)

	if unformattedPorts == "" {
		return 0, 0
	}

	ports := strings.Split(unformattedPorts, "-")

	startPort, err := strconv.Atoi(ports[0])

	if err != nil {
		return 0, 0
	}

	endPort, err := strconv.Atoi(ports[1])

	if err != nil {
		return 0, 0
	}

	return startPort, endPort
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

func outputToFile(fileName string, d string, of string) {
	err := os.WriteFile(fileName, []byte(d), 0755)

	if err != nil {
		PrintfErr("[-] Cannot write to file %s", fileName)
	}
}

func main() {
	var resultPorts []int
	processedPorts := make(chan int)

	target, portRange, threads, verbose, debug, outputFile, outputFormat, protocol := initFlags()

	sPort, ePort := parsePortsArg(portRange)
	ports := make(chan int, threads)

	for i := 0; i < threads; i++ {
		go worker(verbose, debug, protocol, target, processedPorts, ports)
	}

	go poolPorts(ports, sPort, ePort, debug)

	for range ePort + 1 {
		port := <-processedPorts

		if port != 0 {
			resultPorts = append(resultPorts, port)
		}
	}

	close(ports)
	close(processedPorts)

	PrintfSuc("========================= OPEN PORTS ===========================")

	sort.Ints(resultPorts)
	for _, p := range resultPorts {
		PrintfSuc("[+] %d", p)
	}

	if outputFile != "" {
		strP := make([]string, len(resultPorts))
		for i, v := range resultPorts {
			strP[i] = strconv.Itoa(v)
		}

		d := strings.Join(strP, "\n")

		outputToFile(outputFile, d, outputFormat)
	}
}
