package cmd

import (
	"fmt"
	"net"
	"sync"
)

const (
	maxPort         uint16 = 65535
	goroutinesCount        = 5000
)

type AttemptResult struct {
	Port uint16
	Open bool
}

type PortScan struct {
	host       string
	canceled   bool
	cancelLock sync.RWMutex
}

func NewPortScan(host string) *PortScan {
	return &PortScan{host: host}
}

func (p *PortScan) Scan() (res <-chan AttemptResult, cancel chan<- interface{}) { //cancel channel must be closed at the end
	var waitGroup sync.WaitGroup
	results := make(chan AttemptResult, maxPort) //Buf channel for results
	can := make(chan interface{}, 1)             //To be able to cancel the process gracefully
	var port uint16 = 1

	go func() {
		_, ok := <-can //We wait cancel message or closing of the channel
		if ok {        //If it's not closed - cancel flag
			p.markCanceled()
		}
	}()

	waitGroup.Add(int(maxPort)) //We are going to wait for all to be processed or cancelled
	go func() {
		waitGroup.Wait() //We wait all the results to be sent
		close(results)   //And only then we close the channel
	}()

	scanFunction := func(port uint16) { // Function to check the ports
		defer waitGroup.Done() //Anyway we release one counter in waitGroup

		if p.isCanceled() { //Check is canceled before
			return
		}

		hostPort := fmt.Sprintf("%s:%d", p.host, port)
		conn, err := net.Dial("tcp", hostPort)
		var result AttemptResult
		if err != nil { //If no errors -> the port is reachable
			result = AttemptResult{Port: port, Open: false}
		} else {
			conn.Close() //If port is reachable, we have to close the connection
			result = AttemptResult{Port: port, Open: true}
		}

		if p.isCanceled() { //Check is cancelled after
			return
		}

		results <- result
	}

	workChan := make(chan uint16, maxPort) //We dispatch all ports to this chan for processing
	for i := 0; i < goroutinesCount; i++ { //We are going to use limited amount of goroutine to reduce memory consumption
		go func() {
			for port := range workChan {
				scanFunction(port)
			}
		}()
	}

	for port != 0 { //overflow goes to zero on 65535
		workChan <- port //We send every possible port to this channel
		port++
	}
	close(workChan) //and we are ready to close the channel

	return results, can
}

// Goroutine safe, guarded by RWLock
func (p *PortScan) isCanceled() bool {
	p.cancelLock.RLock()
	defer p.cancelLock.RUnlock()

	return p.canceled
}

func (p *PortScan) markCanceled() {
	p.cancelLock.Lock()
	defer p.cancelLock.Unlock()

	p.canceled = true
}
