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
	//"scanme.nmap.org"
	scanCommand := NewPortScan("localhost")
	results, _ := scanCommand.Scan()
	foundPort := false
	for result := range results {
		if result.Open {
			fmt.Printf("Port: %d is open %t\n", result.Port, result.Open)
		}

		if 8080 == result.Port {
			foundPort = true
		}
	}

	if !foundPort {
		t.Errorf("Port 8080 is not found but is open!")
	}
}
