package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"golang.org/x/exp/slices"
)

const initialPort = 8080

var availablePorts = []int{}
var nextPortIndex = 0

func testConnection(port int) bool {
	conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{IP: []byte{127, 0, 0, 1}, Port: port})
	if err != nil {
		println("Dial failed:", err.Error())
		return false
	}
	conn.Close()
	return true
}

func initializeAvailablePorts() {
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

	fmt.Println("Available ports:", availablePorts)
}

func getAvailablePort() int {
	currentAvailablePort := availablePorts[nextPortIndex]
	nextPortIndex = (nextPortIndex + 1) % len(availablePorts)
	return currentAvailablePort
}


func handleRequest(conn net.Conn) {
	// incoming request
	buffer := make([]byte, 1024)
	_, err := conn.Read(buffer)
	if err != nil {
		log.Fatal(err)
	}

	message := string(buffer)
	var errorMessage string

	portStr, request, succeeded := strings.Cut(message, " ")
	if !succeeded {
		errorMessage = `{"status": "error", "message": "Invalid request. Expected format: <port> <request>"}`
		conn.Write([]byte(errorMessage))
		return
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		errorMessage = fmt.Sprintf(`{"status": "error", "message": "Invalid port %v"}`, portStr)
		conn.Write([]byte(errorMessage))
		return
	}

	if port == -1 {
		port = getAvailablePort()
	} else if !slices.Contains(availablePorts, port) {
		errorMessage = fmt.Sprintf(`{"status": "error", "message": "Port %v is not available"}`, port)
		conn.Write([]byte(errorMessage))
		return
	}

	// send request to port
	// TODO: leave these connections open on the testConnection method with a global map
	serverConnection, err := net.Dial("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		errorMessage = fmt.Sprintf(`{"status": "error", "message": "loadbalancer server error: could not connect to port %v"}`, port)
		conn.Write([]byte(errorMessage))
		return
	}

	_, err = serverConnection.Write([]byte(request))
	if err != nil {
		errorMessage = fmt.Sprintf(`{"status": "error", "message": "loadbalancer server error: could not write to port %v"}`, port)
		conn.Write([]byte(errorMessage))
		return
	}

	// get response
	buffer = make([]byte, 1024)
	_, err = serverConnection.Read(buffer)
	if err != nil {
		errorMessage = fmt.Sprintf(`{"status": "error", "message": "loadbalancer server error: could not read from port %v"}`, port)
		conn.Write([]byte(errorMessage))
		return
	}

	// add port to response
	buffer = addPortToBuffer(buffer, port)

	// add newline to response
	buffer = append(buffer, '\n')

	conn.Write(buffer)
}

func addPortToBuffer(buffer []byte, port int) []byte {
	portStr := strconv.Itoa(port)

	// convert buffer to string
	bufferStr := string(buffer)

	// convert bufferStr to json
	bufferStr = strings.Replace(bufferStr, "{", `{"port": "`+portStr+`", `, 1)

	// convert bufferStr to []byte
	buffer = []byte(bufferStr)

	return buffer
}

func listenForRequests() {
	listen, err := net.Listen("tcp", ":4000")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	// close listener
	defer listen.Close()
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		log.Println("Accepted connection")
		go handleRequest(conn)
	}
}

func main() {
	initializeAvailablePorts()
	listenForRequests()
}
