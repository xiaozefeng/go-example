package main

import (
	"log"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:9090")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	_, err = conn.Write([]byte("Hello"))
	if err != nil {
		log.Fatal(err)
	}
	b := make([]byte, 1024)
	_, err = conn.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("resp:", string(b))

}
