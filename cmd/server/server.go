package main

import (
	"fmt"
	"net"
	"net/http"
)

var clientConn net.Conn

func main() {

	//for tcp: this will exchange data between browser and tunnel-client
	go startTcpServer()
	http.HandleFunc("/", handleHttp)
	fmt.Println("Server listening on :3000")
	//will handle browser requests: turnnel-server
	http.ListenAndServe(":3000", nil)
}
func startTcpServer() {
	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		panic(err)
	}
	fmt.Println(" Tcp Server listening on :9000")
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		fmt.Println("Client connected:", conn.RemoteAddr())
		clientConn = conn
	}
}

func handleHttp(w http.ResponseWriter, r *http.Request) {
	if clientConn == nil {
		w.Write([]byte("no connections"))
	}
	reqData := fmt.Sprintf("%s %s", r.Method, r.URL.String())
	clientConn.Write([]byte(reqData))

	buffer := make([]byte, 4096)
	n, _ := clientConn.Read(buffer)

	w.Write(buffer[:n])
}
