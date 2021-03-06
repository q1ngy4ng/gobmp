// Copyright (c) 2018 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.

package bmpconnect

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"time"
)

type BmpMsg struct {
	// A message Id assigned by the connection
	msgId uint
	// BMP message type
	msgType uint8
	// Length of entire BMP message
	msgLen uint32
	// Entire message (including common header)
	msgData []byte
}

func (bmpMsg *BmpMsg) MessageData() []byte {
	return bmpMsg.msgData
}

type BmpConnection struct {
	// Connection created in connectBmp (listener TBD)
	conn net.Conn
	// Counter for generating msg Id's
	msgIdGen uint
	// Messages indexed by msg Id
	msgs map[uint]*BmpMsg
}

func (bmpConn *BmpConnection) Message(index uint) (*BmpMsg, bool) {
	msg, ok := bmpConn.msgs[index]
	return msg, ok
}

//
// TODO: make this a non-blocking connection?
//
func ConnectBmp(address string, port uint) (*BmpConnection, error) {
	addr := fmt.Sprintf("%s:%d", address, port)

	conn, err := net.Dial("tcp", addr)

	if err != nil {
		return nil, err
	}

	bmpConn := new(BmpConnection)
	bmpConn.conn = conn
	bmpConn.msgIdGen = 0
	bmpConn.msgs = make(map[uint]*BmpMsg)

	return bmpConn, nil
}

// Communication with goroutine servicing Bmp connection
// - The service function reads channel for instruction
// - Instructions:
//   - read N messages
//   - disconnect from Bmp
//   - terminate
//
// When the instruction is completed the status is written back

const (
	// Send ReadMsg on the channel followed by the number of messages to read
	ReadMsg    = 1
	Disconnect = 2
	Terminate  = 3
)

// Status returned through channel to communicate result to requester
const (
	Ok      = 0
	Error   = 1
	Timeout = 2
)

//
// Read numMsgs on connection. The timeout argument is in seconds, or -1 to wait forever
//
func (bmpConn *BmpConnection) readBmpMsgs(numMsgs int, timeout int) (int, error) {
	msgCount := 0
	tmp := make([]byte, 6)
	var err error
	err = nil
	for msgCount < numMsgs {
		if timeout > 0 {
			curTime := time.Now()
			deadline := curTime.Add(time.Duration(timeout) * time.Second)
			bmpConn.conn.SetReadDeadline(deadline)
		}
		_, err := io.ReadFull(bmpConn.conn, tmp)
		bmpConn.conn.SetReadDeadline(time.Time{})
		if err != nil {
			if err != io.EOF {
				fmt.Println("read error:", err)
			}
			if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
				fmt.Println("read timeout", neterr)
			}
			break
		}
		//fmt.Println("Read", n, "bytes")
		version := tmp[0]
		if version != 3 {
			// TODO fix bogus error
			err = errors.New(fmt.Sprintf("Invalid BMP Version %d", version))
			break
		}
		//fmt.Println("Version:", version)
		len := binary.BigEndian.Uint32(tmp[1:5])
		//fmt.Println("Length:", len)
		msgType := tmp[5]
		//fmt.Println("Type:", msgType)

		// Read len bytes into buffer
		msgData := make([]byte, len)
		copy(msgData[0:], tmp[:])
		if timeout > 0 {
			curTime := time.Now()
			deadline := curTime.Add(time.Duration(timeout) * time.Second)
			bmpConn.conn.SetReadDeadline(deadline)
		}
		_, err = io.ReadFull(bmpConn.conn, msgData[6:])
		bmpConn.conn.SetReadDeadline(time.Time{})
		if err != nil {
			if err != io.EOF {
				fmt.Println("read error:", err)
			}
			if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
				fmt.Println("read timeout", neterr)
			}
			break
		}
		//fmt.Println("Read", n, "bytes")
		//fmt.Println("add message index", bmpConn.msgIdGen)
		bmpMsg := new(BmpMsg)
		bmpMsg.msgId = bmpConn.msgIdGen
		bmpMsg.msgType = msgType
		bmpMsg.msgLen = len
		bmpMsg.msgData = msgData
		bmpConn.msgs[bmpMsg.msgId] = bmpMsg

		bmpConn.msgIdGen = bmpConn.msgIdGen + 1
		msgCount++
	}
	return msgCount, err
}

func (bmpConn *BmpConnection) ServiceBmpConnection(c chan int) {
	for {
		cmd := <-c
		//fmt.Println("ServiceBmpConnection cmd:", cmd)
		switch cmd {
		case ReadMsg:
			numMsgs := <-c
			timeout := <-c
			msgCount, err := bmpConn.readBmpMsgs(numMsgs, timeout)
			if err == nil {
				c <- 0
				c <- msgCount
			} else {
				if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
					c <- Timeout
					c <- msgCount
				} else {
					c <- Error
					c <- msgCount
				}
			}
		case Disconnect:
			bmpConn.conn.Close()
			c <- 0
		case Terminate:
			c <- 0
			break
		default:
			fmt.Println("serviceBmpConnection: Invalid cmd")
		}
	}
	//fmt.Println("exit ServiceBmpConnection")
}
