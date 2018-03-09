package main

// to run the test run "nc -l 10000 < bmp_messages.bin" in bash and then run this program
// go run test_connect.go connect.go

import "os"
import "fmt"

func main() {
	fmt.Println("Connect to 127.0.0.1 port 10000")
	bmpConn, err := connectBmp("127.0.0.1", 10000)
	if err != nil {
		fmt.Printf("Error connecting to Bmp")
		os.Exit(-1)
	}
	fmt.Println("Connected to Bmp speaker")
	c := make(chan int)
	go bmpConn.serviceBmpConnection(c)
	c <- ReadMsg
	c <- 5
	status := <-c
	msgCount := <-c
	if status == 0 {
		fmt.Println("Read", msgCount, "msgs")
	}
	// Check first message is type 4 (Initiation)
	if bmpConn.msgs[0].msgType == 4 {
		fmt.Println("Found Initiation Message")
	}
	// Now just read until Termination message
	for {
		c <- ReadMsg
		c <- 1
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
