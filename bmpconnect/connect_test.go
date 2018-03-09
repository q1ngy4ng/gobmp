// Copyright (c) 2018 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.

package bmpconnect

// go test

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

func TestBmpConnect(t *testing.T) {
	t.Log("Start BMP Speaker")
	fmt.Println("Start BMP Speaker")
	command := exec.Command("/usr/bin/nc", "-l", "10000")
	f, err := os.Open("bmp_messages.bin")
	if err != nil {
		fmt.Printf("error opening file: %v", err)
		t.Fail()
	}
	defer f.Close()
	// On this line you're going to redirect the output to a file
	command.Stdin = f
	if err := command.Start(); err != nil {
		fmt.Fprintln(os.Stderr, "Command failed.", err)
		t.Fail()
	}

	fmt.Println("Connect to 127.0.0.1 port 10000")
	bmpConn, err := ConnectBmp("127.0.0.1", 10000)
	if err != nil {
		fmt.Printf("Error connecting to Bmp")
		t.Fail()
	}
	fmt.Println("Connected to Bmp speaker")
	c := make(chan int)
	go bmpConn.ServiceBmpConnection(c)
	c <- ReadMsg
	c <- 5
	c <- 3
	status := <-c
	msgCount := <-c
	if status == 0 {
		fmt.Println("Read", msgCount, "msgs")
	}
	// Check first message is type 4 (Initiation)
	if bmpConn.msgs[0].msgType == 4 {
		fmt.Println("Found Initiation Message")
	} else {
		t.Log("Failed to find Initiation Message")
		t.Fail()
	}
	// Now just read until Termination message
	for {
		c <- ReadMsg
		c <- 1
		c <- 3
		status = <-c
		msgCount = <-c
		index := uint(len(bmpConn.msgs) - 1)
		if bmpConn.msgs[index].msgType == 5 {
			fmt.Println("Found Termination Message")
			break
		}
	}
	c <- Terminate
	status = <-c
	if status == 0 {
		fmt.Println("Terminated successfully")
	}

	// Print out all message types found
	nMsgsRcvd := len(bmpConn.msgs)
	for i := 0; i < nMsgsRcvd; i++ {
		fmt.Println(bmpConn.msgs[uint(i)])
	}
}
