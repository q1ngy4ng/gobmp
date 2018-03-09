package bmpconnect

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
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

type BmpConnection struct {
	// Connection created in connectBmp (listener TBD)
	conn net.Conn
	// Counter for generating msg Id's
	msgIdGen uint
	// Messages indexed by msg Id
	msgs map[uint]*BmpMsg
}

//
// TODO: make this a non-blocking connection?
//
func connectBmp(address string, port uint) (*BmpConnection, error) {
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

func (bmpConn *BmpConnection) readBmpMsgs(numMsgs int) (int, error) {
	msgCount := 0
	tmp := make([]byte, 6)
	var err error
	err = nil
	for msgCount < numMsgs {
		_, err := io.ReadFull(bmpConn.conn, tmp)
		if err != nil {
			if err != io.EOF {
				fmt.Println("read error:", err)
			}
			break
		}
		//fmt.Println("Read", n, "bytes")
		//version := tmp[0]
		//fmt.Println("Version:", version)
		len := binary.BigEndian.Uint32(tmp[1:5])
		//fmt.Println("Length:", len)
		msgType := tmp[5]
		//fmt.Println("Type:", msgType)

		// Read len bytes into buffer
		msgData := make([]byte, len)
		copy(msgData[0:], tmp[:])
		_, err = io.ReadFull(bmpConn.conn, msgData[6:])
		if err != nil {
			if err != io.EOF {
				fmt.Println("read error:", err)
			}
			break
		}
		//fmt.Println("Read", n, "bytes")
		fmt.Println("add message index", bmpConn.msgIdGen)
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

func (bmpConn *BmpConnection) serviceBmpConnection(c chan int) {
	for {
		cmd := <-c
		switch cmd {
		case ReadMsg:
			numMsgs := <-c
			msgCount, err := bmpConn.readBmpMsgs(numMsgs)
			if err == nil {
				c <- 0
				c <- msgCount
			}
		case Disconnect:
			bmpConn.conn.Close()
			c <- 0
		case Terminate:
			c <- 0
			break
		}
	}
}
