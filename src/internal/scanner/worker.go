package scanner

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/Alexandre-Chapelle/port-scanner/src/internal/ui"
)

func Worker(v bool, d bool, proc string, t string, pp chan int, ports chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	for p := range ports {
		start := time.Now()
		address := net.JoinHostPort(t, fmt.Sprintf("%d", p))
		conn, err := net.Dial(proc, address)

		if err != nil {

			if d {
				ui.PrintfErr(err.Error())
			}

			pp <- 0

			if v {
				tElapsed := time.Since(start)
				ui.PrintfErr("[-] Port %d is closed, %dms", p, tElapsed.Milliseconds())
			}

			continue
		}

		if v {
			tElapsed := time.Since(start)
			ui.PrintfSuc("[+] Port %d is open, %dms", p, tElapsed.Milliseconds())
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
