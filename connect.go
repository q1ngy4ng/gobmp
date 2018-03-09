package main

import (
	"fmt"
	"net"
)

// TODO: make this a non-blocking connection?

func connectBmp(address string, port uint) (net.Conn, error) {
	// Return channel which is used to send command
	// strings to the connection routine
	//c := make(chan string)

	addr := fmt.Sprintf("%s:%d", address, port)

	conn, err := net.Dial("tcp", addr)

	if err != nil {
		return nil, err
	}

	return conn, nil
}

func serviceBmpConnection(conn net.Conn, c chan string) {
}
