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

type Reader struct {
	r   io.Reader
	err error
}

func (r *Reader) read(data any) {
	if r.err == nil {
		r.err = binary.Read(r.r, binary.BigEndian, data)
	}
}

func (p *Package) Unpack(input io.Reader) error {
	r := &Reader{r: input}
	r.read(&p.MagicNum)
	r.read(&p.Version)
	r.read(&p.Alg)
	r.read(&p.Order)
	r.read(&p.Len)
	p.Data = make([]byte, p.Len)
	r.read(&p.Data)
	if r.err != nil {
		return r.err
	}
	if p.MagicNum != X0001 {
		return errors.New("not my magic number")
	}
	return r.err
}

type Writer struct {
	w   io.Writer
	err error
}

func (w *Writer) write(data any) {
	if w.err == nil {
		w.err = binary.Write(w.w, binary.BigEndian, data)
	}
}

func (p *Package) Pack(writer io.Writer) error {
	w := Writer{w: writer}
	w.write(&p.MagicNum)
	w.write(&p.Version)
	w.write(&p.Alg)
	w.write(&p.Order)
	w.write(&p.Len)
	w.write(&p.Data)
	return w.err
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
			if err == io.EOF || err == net.ErrClosed {
				return
			}
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
