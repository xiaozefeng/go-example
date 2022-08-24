package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
)

// X0001 MagicNum
const (
	X0001 = 0x114301
)

// Version
const (
	_ byte = iota
	V1
)

// Alg
const (
	_ byte = iota
	AlgJSOM
)

// Order
const (
	_ byte = iota
	Login
)

type Package struct {
	MagicNum int32  // 魔数
	Version  byte   // 版本号
	Alg      byte   // 序列化算法
	Order    byte   // 指令
	Len      int32  // 数据长度
	Data     []byte // 数据
}

func (p *Package) String() string {
	return fmt.Sprintf("magicNum:%d, version:%d, alg:%d, order:%d, len:%d, data:%s",
		p.MagicNum, p.Version, p.Alg, p.Order, p.Len, p.Data)
}

type errWrapper struct {
	io.Reader
	err error
}

func (e *errWrapper) Read(p []byte) (int, error) {
	if e.err != nil {
		return 0, e.err
	}
	var n int
	n, e.err = e.Reader.Read(p)
	return n, e.err
}

func (p *Package) Unpack(reader io.Reader) error {
	r := &errWrapper{Reader: reader}
	r.err = binary.Read(r, binary.BigEndian, &p.MagicNum)
	if p.MagicNum != X0001 {
		return errors.New("not my magic number")
	}
	r.err = binary.Read(r, binary.BigEndian, &p.Version)
	r.err = binary.Read(r, binary.BigEndian, &p.Alg)
	r.err = binary.Read(r, binary.BigEndian, &p.Order)
	r.err = binary.Read(r, binary.BigEndian, &p.Len)
	p.Data = make([]byte, p.Len)
	r.err = binary.Read(r, binary.BigEndian, &p.Data)
	return r.err
}

func (p *Package) Pack(writer io.Writer) error {
	var err error
	err = binary.Write(writer, binary.BigEndian, &p.MagicNum)
	err = binary.Write(writer, binary.BigEndian, &p.Version)
	err = binary.Write(writer, binary.BigEndian, &p.Alg)
	err = binary.Write(writer, binary.BigEndian, &p.Order)
	err = binary.Write(writer, binary.BigEndian, &p.Len)
	err = binary.Write(writer, binary.BigEndian, &p.Data)
	return err
}

var addr string

func main() {
	flag.StringVar(&addr, "a", ":9090", "tcp server addr")
	flag.Parse()
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("listen tcp server on", addr)
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("conn request error:", err)
			continue
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	for {

		pkg, err := DecodePackage(conn)
		if err != nil {
			log.Println("decode pkg error:", err)
			return
		}
		log.Printf("%+v", pkg)
		switch pkg.Order {
		case Login:
			resp, err := login(pkg.Alg, pkg.Data)
			if err != nil {
				log.Println("process login err", err)
				return
			}
			b, err := json.Marshal(resp)
			if err != nil {
				log.Println("process login resp err", err)
				return
			}
			pkg.Len = int32(len(b))
			pkg.Data = b
			buf := new(bytes.Buffer)
			err = pkg.Pack(buf)
			if err != nil {
				log.Println("package login resp error", err)
				return
			}
			conn.Write(buf.Bytes())
		default:
			log.Println("unknown order: ", pkg.Order)
		}
	}
}

type LoginParam struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type LoginResp struct {
	Id    string `json:"id"`
	Token string `json:"token"`
}

func login(alg byte, data []byte) (*LoginResp, error) {
	var p LoginParam
	switch alg {
	case AlgJSOM:
		err := json.Unmarshal(data, &p)
		if err != nil {
			return nil, err
		}
		log.Printf("login param: %+v", p)
		return &LoginResp{Id: "U1134", Token: "token114312"}, err
	default:
		log.Println("un supported alg")
		return nil, errors.New("unsupported serialize alg")
	}
}

func DecodePackage(conn net.Conn) (*Package, error) {
	var p = &Package{}
	err := p.Unpack(conn)
	return p, err
}
