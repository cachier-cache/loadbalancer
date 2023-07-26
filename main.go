package main

import (
	"bufio"
	"fmt"
	"io"
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
		return false
	}
	defer conn.Close()
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
	defer conn.Close()

	reader := bufio.NewReader(conn)

	for {
		// incoming request
		message, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				// client closed the connection
				break
			}
			log.Fatal(err)
		}

		message = strings.TrimSuffix(message, "\n") // remove newline
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
		serverConnection, err := net.Dial("tcp", fmt.Sprintf(":%v", port))
		if err != nil {
			errorMessage = fmt.Sprintf(`{"status": "error", "message": "loadbalancer server error: could not connect to port %v"}`, port)
			conn.Write([]byte(errorMessage))
			return
		}
		defer serverConnection.Close()

		_, err = serverConnection.Write([]byte(request))
		if err != nil {
			errorMessage = fmt.Sprintf(`{"status": "error", "message": "loadbalancer server error: could not write to port %v"}`, port)
			conn.Write([]byte(errorMessage))
			return
		}

		// get response
		buffer := make([]byte, 1024)
		n, err := serverConnection.Read(buffer)
		if err != nil {
			errorMessage = fmt.Sprintf(`{"status": "error", "message": "loadbalancer server error: could not read from port %v"}`, port)
			conn.Write([]byte(errorMessage))
			return
		}

		// add port to response
		buffer = addPortToBuffer(buffer[:n], port)

		// add newline to response
		buffer = append(buffer, '\n')

		conn.Write(buffer)
	}
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
