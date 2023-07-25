package main

import (
	"fmt"
	"log"
	"net"
)

const initialPort = 8080
var availablePorts = []int{}

func testConnection(port int) bool {
	conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{IP: []byte{127, 0, 0, 1}, Port: port})
	if err != nil {
		println("Dial failed:", err.Error())
		return false
	}
	conn.Close()
	return true
}

func main() {
	for i := initialPort; i < initialPort+100; i++ {
		if testConnection(i) {
			availablePorts = append(availablePorts, i)
			continue
		}

		break
	}

	if len(availablePorts) == 0 {
		log.Fatal("No available ports")
	}

	fmt.Printf("Available ports: %v", availablePorts)

	// TODO: receive requests and route then to a port via round robin
	// and then return the port back to the client

	// the client will need to keep track of the port
	// the client should have load balancing off by default
}
