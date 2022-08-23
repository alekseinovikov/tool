package cmd

import (
	"fmt"
	"github.com/rivo/tview"
	"net"
	"sync"
)

const (
	maxPort         uint16 = 65535
	goroutinesCount        = 2000
)

type attemptResult struct {
	port uint16
	open bool
}

type portScan struct {
	canceled   bool
	cancelLock sync.RWMutex
}

func (p *portScan) Run(app *tview.Application, main tview.Primitive) {
}

func NewPortScanCommand() Command {
	return &portScan{canceled: false}
}

func (p *portScan) scan(url string) (res <-chan attemptResult, cancel chan<- interface{}) {
	var waitGroup sync.WaitGroup
	results := make(chan attemptResult)
	can := make(chan interface{})
	var port uint16 = 1

	go func() {
		waitGroup.Wait()
		close(results)
	}()

	go func() {
		_, ok := <-can
		if ok {
			p.markCanceled()
		}
	}()

	waitGroup.Add(int(maxPort))
	scanFunction := func(port uint16) {
		defer waitGroup.Done()

		if p.isCanceled() {
			return
		}

		hostPort := fmt.Sprintf("%s:%d", url, port)
		conn, err := net.Dial("tcp", hostPort)
		var result attemptResult
		if err != nil {
			result = attemptResult{port: port, open: false}
		} else {
			conn.Close()
			result = attemptResult{port: port, open: true}
		}

		if p.isCanceled() {
			return
		}

		results <- result
	}

	workChan := make(chan uint16, maxPort)
	for i := 0; i < goroutinesCount; i++ {
		go func() {
			for port := range workChan {
				scanFunction(port)
			}
		}()
	}

	for port != 0 { //overflow goes to zero on 65535
		workChan <- port
		port++
	}
	close(workChan)

	return results, can
}

func (p *portScan) isCanceled() bool {
	p.cancelLock.RLock()
	defer p.cancelLock.RUnlock()

	return p.canceled
}

func (p *portScan) markCanceled() {
	p.cancelLock.Lock()
	defer p.cancelLock.Unlock()

	p.canceled = true
}
