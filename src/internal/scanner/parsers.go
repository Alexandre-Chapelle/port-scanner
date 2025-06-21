package scanner

import (
	"regexp"
	"strconv"
	"strings"
)

func ParsePortsArg(portsArg string) (start int, end int) {
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
