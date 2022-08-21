package main

import (
	"flag"
	"io"
	"log"
	"net"
)

// usage
// nc localhost 9090
// input hi
var (
	addr string
)

func init() {
	flag.StringVar(&addr, "addr", ":9090", "tcp server addr")
	flag.Parse()
}

func main() {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("handle incoming request err:", err)
		}
		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	buffer := make([]byte, 1024)
	_, err := conn.Read(buffer)
	if err != nil {
		if err != io.EOF {
			log.Println("handle conn ,err:", err)
		}
	}
	log.Printf("request: %s\n", buffer)
	//now := time.Now().Format("2006-01-02 15:04:05")
	conn.Write([]byte("HTTP/1.1 200 OK\n"))
	conn.Close()
}
