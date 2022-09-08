package https

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/url"
)

type ProxyServer struct {
	network string
	address string
}
type Connection struct {
	conn net.Conn
}

func NewConnection(conn net.Conn) *Connection {
	return &Connection{
		conn: conn,
	}
}
func NewServer(network string, address string) *ProxyServer {
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
	buf := make([]byte, 1024)
	n, err := client.Read(buf)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	var method, rawURL string
	end := bytes.IndexByte(buf[:n], '\n')
	if end == -1 {
		return
	}
	fmt.Sscanf(string(buf[:n]), "%s%s", &method, &rawURL)
	var server net.Conn
	switch method {
	case "CONNECT":
		server, err = net.Dial("tcp", rawURL)
		if err != nil {
			fmt.Printf("https err: %v\n", err)
			return
		}
		client.Write([]byte("HTTP/1.1 200 Connection established\r\n\r\n"))
	default:
		URL, _ := url.Parse(rawURL)
		if URL.Port() == "" {
			URL.Host += ":80"
		}
		server, err = net.Dial("tcp", URL.Host)
		if err != nil {
			fmt.Printf("http err: %v\n", err)
			return
		}
		server.Write(buf[:n])
	}
	go io.Copy(server, client)
	go io.Copy(client, server)
}
func main() {
	server := NewServer("tcp", ":443")
	server.Start()
}
