package scanner

import (
	"sync"
)

type Config struct {
	Target       string
	PortRange    string
	Threads      int
	Verbose      bool
	Debug        bool
	OutputFile   string
	OutputFormat string
	Protocol     string
}

type Scanner struct {
	config Config
}

func New(cfg Config) *Scanner {
	return &Scanner{config: cfg}
}

func (s *Scanner) Scan() (results []int, err error) {
	var wg sync.WaitGroup
	var resultPorts []int
	processedPorts := make(chan int)

	sPort, ePort := ParsePortsArg(s.config.PortRange)
	ports := make(chan int, s.config.Threads)

	for i := 0; i < s.config.Threads; i++ {
		wg.Add(1)
		go Worker(s.config.Verbose, s.config.Debug, s.config.Protocol, s.config.Target, processedPorts, ports, &wg)
	}

	go PoolPorts(ports, sPort, ePort, s.config.Debug)

	for range ePort + 1 {
		port := <-processedPorts

		if port != 0 {
			resultPorts = append(resultPorts, port)
		}
	}

	close(ports)
	close(processedPorts)

	wg.Wait()

	return resultPorts, nil
}
