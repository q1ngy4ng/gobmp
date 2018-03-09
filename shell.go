package main

import (
	"bufio"
	"fmt"
	"os"
	"io"
	"strings"
	"strconv"
	"net"
	"gobmp/bmpconnect"
)

/*
gobmp> connect 1.1.1.1:3000
gobmp> connect 2.2.2.2:3000
gobmp> connections
1.1.1.1:3000
2.2.2.2:3000
gobmp> read-messages 1.1.1.1:3000 5 3
gobmp> read-messages 2.2.2.3:3000 10 5
gobmp> disconnect 1.1.1.1:3000
gobmp> connections
2.2.2.2:3000
*/

func main() {
	connections := make(map[string]chan int)
	err := evalCommand("connect " + os.Args[1], connections)
	if err != nil {
		fmt.Println(err)
	}
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("gobmp> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				fmt.Errorf("error")
			}
			break
		}
		err = evalCommand(line, connections)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func evalCommand(line string, connections map[string]chan int) error {
	cmdParts := strings.Fields(line)
	if len(cmdParts) == 0 {
		return nil
	}
	switch cmdParts[0] {
	case "connect":
		return evalCmdConnect(cmdParts, connections)
	case "disconnect":
		return evalCmdDisconnect(cmdParts, connections)
	case "read-messages":
		return evalCmdReadMessages(cmdParts, connections)
	case "show-connections":
		return evalCmdShowConnections(cmdParts, connections)
	default:
		return fmt.Errorf("invalid command")
	}
}

func evalCmdConnect(cmdParts []string, connections map[string]chan int) error {
	if !isValidIpAndPort(cmdParts[1]) {
		return fmt.Errorf("malformed ip addresss and port")
	}
	parts := strings.Split(cmdParts[1], ":")
	port, _ := strconv.ParseUint(parts[1], 10, 32)
	bmpConnection, err := bmpconnect.ConnectBmp(parts[0], uint(port))
	if err != nil {
		return fmt.Errorf("connection refused")
	}
	c := make(chan int)
	go bmpConnection.ServiceBmpConnection(c)
	connections[cmdParts[1]] = c
	return nil
}

func isValidIpAndPort(ipAndPort string) bool {
	parts := strings.Split(ipAndPort, ":")
	if len(parts) != 2 {
		return false
	}
	netIP := net.ParseIP(parts[0])
	if netIP == nil {
		return false
	}
	port, err := strconv.ParseUint(parts[1], 10, 32)
	if err != nil || port < 0 || port > 65535 {
		return false
	}
	return true
}

func evalCmdDisconnect(cmdParts []string, connections map[string]chan int) error {
	k := cmdParts[1]
	c, ok := connections[k]
	if !ok {
		return fmt.Errorf("no connection to %s", cmdParts[1])
	}
	c <- bmpconnect.Terminate
	delete(connections, k)
	return nil
}

func evalCmdReadMessages(cmdParts []string, connections map[string]chan int) error {
	if len(cmdParts) != 4 {
		return fmt.Errorf("invalid command")
	}
	k := cmdParts[1]
	c, ok := connections[k]
	if !ok {
		return fmt.Errorf("no connection to %s", cmdParts[1])
	}
	c <- bmpconnect.ReadMsg
	numMsgs, _ := strconv.ParseInt(cmdParts[2], 10, 32)
	timeout, _ := strconv.ParseInt(cmdParts[3], 10, 32)
	c <- int(numMsgs)
	c <- int(timeout)
	return nil
}

func evalCmdShowConnections(cmdParts []string, connections map[string]chan int) error {
	if len(cmdParts) != 1 {
		return fmt.Errorf("invalid command")
	}
	for k, _ := range connections {
		fmt.Println(k)
	}
	return nil
}
