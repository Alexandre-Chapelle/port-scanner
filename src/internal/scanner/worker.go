package scanner

import (
	"fmt"
	"net"
	"sync"

	"github.com/Alexandre-Chapelle/port-scanner/src/internal/ui"
)

func Worker(v bool, d bool, proc string, t string, pp chan int, ports chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	for p := range ports {

		address := net.JoinHostPort(t, fmt.Sprintf("%d", p))
		conn, err := net.Dial(proc, address)

		if err != nil {
			if d {
				ui.PrintfErr(err.Error())
			}

			pp <- 0

			if v {
				ui.PrintfErr("[-] Port %d is closed", p)
			}

			continue
		}

		if v {
			ui.PrintfSuc("[+] Port %d is open", p)
		}

		pp <- p

		conn.Close()
	}
}

func PoolPorts(ports chan int, sPort int, ePort int, d bool) {
	for i := sPort; i <= ePort; i++ {
		ports <- i

		if d {
			ui.PrintfInfo("[~] Added port %d to queue", i)
		}
	}
}
