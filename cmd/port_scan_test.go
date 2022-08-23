package cmd

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

func Test_scan(t *testing.T) {
	go func() { http.ListenAndServe(":8080", http.DefaultServeMux) }()
	time.Sleep(time.Second) //waiting for 8080 to be open
	scanCommand := portScan{}

	//results, _ := scanCommand.scan("scanme.nmap.org")
	results, _ := scanCommand.scan("localhost")
	foundPort := false
	for result := range results {
		if result.open {
			fmt.Printf("Port: %d is open %t\n", result.port, result.open)
		}

		if 8080 == result.port {
			foundPort = true
		}
	}

	if !foundPort {
		t.Errorf("Port 8080 is not found but is open!")
	}
}
