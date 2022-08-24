package main

import (
	"bytes"
	"fmt"
	"net"
	"testing"
)

func TestClient(t *testing.T) {
	conn, err := net.Dial("tcp", ":9090")
	if err != nil {
		t.Error(err)
	}
	defer conn.Close()

	data := []byte("{\"username\":\"mickey\", \"password\":\"112w35235\"}")
	p := &Package{
		MagicNum: X0001,
		Version:  V1,
		Alg:      AlgJSOM,
		Order:    Login,
		Len:      int32(len(data)),
		Data:     data,
	}
	buf := new(bytes.Buffer)
	err = p.Pack(buf)
	if err != nil {
		t.Error(err)
	}
	conn.Write(buf.Bytes())

	r := &Package{}
	err = r.Unpack(conn)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("resp:", r)
}
