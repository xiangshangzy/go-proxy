package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

type ProxyServer struct {
	network string
	address string
}

func NewProxyServer(network string, address string) *ProxyServer {
	return &ProxyServer{network: network, address: address}
}
func (this *ProxyServer) Start() {
	listenr, err := net.Listen(this.network, this.address)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	for {
		conn, err := listenr.Accept()
		if err != nil {
			continue
		}
		go HandConn(conn)
	}
}
func HandConn(client net.Conn) {
	err := ClientHandshake(client)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	server, err := ServerConnect(client)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	go io.Copy(server, client)
	go io.Copy(client, server)
}
func ClientHandshake(client net.Conn) error {
	head := make([]byte, 2)
	_, err := io.ReadFull(client, head)
	if err != nil {
		return err
	}
	version := head[0]
	n := head[1]
	if version != 5 {
		return fmt.Errorf("unexpected version: %v", version)
	}
	if n > 0 {
		methods := make([]byte, n)
		_, err := io.ReadFull(client, methods)
		if err != nil {
			fmt.Printf("err: %v\n", err)
		}
	}
	client.Write([]byte{5, 0})
	return nil
}
func ServerConnect(client net.Conn) (net.Conn, error) {
	head := make([]byte, 4)
	io.ReadFull(client, head)
	addrType := head[3]
	var host string
	switch addrType {
	case 1:
		hostBuf := make([]byte, net.IPv4len+2)
		io.ReadFull(client, hostBuf)
		hostname := net.IP(hostBuf[:net.IPv4len]).String()
		port := binary.BigEndian.Uint16(hostBuf[net.IPv4len:])
		host = net.JoinHostPort(hostname, strconv.Itoa(int(port)))
	case 3:
		lenBuf := make([]byte, 1)
		io.ReadFull(client, lenBuf)
		hostBuf := make([]byte, lenBuf[0])
		io.ReadFull(client, hostBuf)
		host = string(hostBuf)
		if !strings.Contains(host, ":") {
			host += ":80"
		}
	case 4:
		hostBuf := make([]byte, net.IPv6len+2)
		io.ReadFull(client, hostBuf)
		hostname := net.IP(hostBuf[:net.IPv6len]).String()
		port := binary.BigEndian.Uint16(hostBuf[net.IPv6len:])
		host = net.JoinHostPort(hostname, strconv.Itoa(int(port)))
	}
	resp := make([]byte, 10)
	resp[0] = 5
	resp[3] = 1
	server, err := net.Dial("tcp", host)
	if err != nil {
		return nil, err
	}
	client.Write(resp)

	return server, nil
}
func main() {
	server := NewProxyServer("tcp", ":1080")
	server.Start()
}
