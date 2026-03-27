package main

import (
	"fmt"
	"net"
	"strconv"
)

var i int = 0

func main() {

	//for tcp
	channel := make(chan int)
	go startTcpServer()
	<-channel

}
func startTcpServer() {
	listener, err := net.Listen("tcp", "localhost:9000")
	if err != nil {
		panic(err)
	}
	fmt.Println("Server listening on :9000")
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Accept error:", err)
			continue
		}
		fmt.Println("Client connected:", conn.RemoteAddr())
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 1024)

	for {
		n, err := conn.Read(buffer)
		i++
		if err != nil {
			fmt.Println("Client disconnected")
			return
		}

		fmt.Println("Received data:", string(buffer[:n]))

		conn.Write([]byte("Message received \n" + strconv.Itoa(i)))
	}
}
