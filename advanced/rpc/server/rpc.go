package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"reflect"
	"unsafe"
)

type UserService struct {
	Login func(username string, password string) (LoginResp, error)
}

func NewProxy[T any](target T) (T, error) {
	var t = target
	typeOf := reflect.TypeOf(target)
	valueOf := reflect.ValueOf(target)
	if typeOf.Kind() == reflect.Ptr {
		typeOf = typeOf.Elem()
		valueOf = valueOf.Elem()
	}
	if typeOf.Kind() != reflect.Struct {
		return t, errors.New("type error")
	}

	numField := typeOf.NumField()
	for i := 0; i < numField; i++ {
		field := typeOf.Field(i)
		f := valueOf.Field(i)
		if f.Kind() == reflect.Func && f.IsValid() && f.CanSet() {
			numOut := field.Type.NumOut()
			funcOuts := make([]reflect.Value, 0, numOut)
			//for i := 0; i < numOut; i++ {
			//	funcOuts = append(funcOuts, reflect.Zero(field.Type.Out(i)))
			//}
			numIn := field.Type.NumIn()
			makeFunc := reflect.MakeFunc(field.Type, func(args []reflect.Value) []reflect.Value {
				ps := make(map[string]any)
				for i := 0; i < numIn; i++ {
					val := args[i]
					ps[val.String()] = val.Interface()
				}
				b, err := json.Marshal(ps)
				if err != nil {
					fmt.Println("marshal error:", err)
					return funcOuts
				}
				client := NewClient("tcp", ":9090")
				bs, err := client.Call(b)
				if err != nil {
					fmt.Println("call rpc error:", err)
					return funcOuts
				}
				var result any
				err = json.Unmarshal(bs, &result)
				if err != nil {
					fmt.Println("unmarshal error:", err)
					return funcOuts
				}
				fmt.Println("result:", result)
				out1 := reflect.NewAt(field.Type.Out(0), unsafe.Pointer(&result))
				funcOuts = append(funcOuts, out1)
				funcOuts = append(funcOuts, reflect.ValueOf(err))
				return funcOuts
			})
			valueOf.Field(i).Set(makeFunc)
		}
	}
	return target, nil
}

type Client struct {
	network string
	addr    string
}

func NewClient(network, addr string) *Client {
	return &Client{
		network: network,
		addr:    addr,
	}
}

func (c *Client) Call(data []byte) ([]byte, error) {
	fmt.Printf("request %s\n", data)
	conn, err := net.Dial(c.network, c.addr)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

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
		return nil, err
	}
	_, err = conn.Write(buf.Bytes())
	if err != nil {
		return nil, err
	}
	r := &Package{}
	err = r.Unpack(conn)
	if err != nil {
		return nil, err
	}
	fmt.Println("resp:", r)

	return r.Data, nil
}
