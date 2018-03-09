package main

import "os"
import "fmt"

func main() {
	fmt.Println("Connect to 127.0.0.1 port 10000")
	conn, err := connectBmp("127.0.0.1", 10000)
	if err != nil {
		fmt.Printf("Error connecting to Bmp")
		os.Exit(-1)
	}
	c := make(chan string)
	go serviceBmpConnection(conn, c)
}
