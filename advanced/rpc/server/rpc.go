package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"reflect"
)

type UserService struct {
	Login func(*LoginParam) (*LoginResp, error)
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
			makeFunc := reflect.MakeFunc(field.Type, func(args []reflect.Value) []reflect.Value {
				if len(args) > 0 {
					param := args[0].Elem().Interface()
					b, err := json.Marshal(param)
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
					for i := 0; i < numOut; i++ {
						outField := field.Type.Out(i)
						switch {
						case outField.Kind() == reflect.Ptr || outField.Kind() == reflect.Struct:
							val := reflect.New(field.Type.Out(i))
							err = json.Unmarshal(bs, val.Interface())
							if err != nil {
								fmt.Println("unmarshal error:", err)
								return funcOuts
							}
							funcOuts = append(funcOuts, val.Elem())
						case outField.Implements(reflect.TypeOf(new(error)).Elem()):
							//还没有找到怎么判断是否Error的方法
							//default:
							funcOuts = append(funcOuts, reflect.ValueOf(&err).Elem())
						}
					}
					return funcOuts
				}
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
