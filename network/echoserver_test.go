package network

import (
	"log"
	"testing"
)

type echoServer struct {
	server *TCPServer
}

func newEchoServer(addr string) *echoServer {
	srv := &echoServer{
		server: NewTCPServer(addr),
	}
	return srv
}

func (srv *echoServer) ListenAndServe() {
	if err := srv.server.ListenAndServe(srv, nil); err != nil {
		log.Println(err)
	}
}

func (srv *echoServer) Close() {
	srv.server.Close()
}

func (srv *echoServer) Connect(connection *TCPConnection, connected bool) {
	if connected {
		log.Printf("echo server: %v connected.\n", connection.RemoteAddr())
	} else {
		log.Printf("echo server: %v disconnected.\n", connection.RemoteAddr())
	}
}

func (srv *echoServer) Receive(connection *TCPConnection, b []byte) {
	message := string(b)
	log.Printf("echo server: %v receive %v\n", connection.RemoteAddr(), message)
	log.Println("echo server: send", message)
	connection.Send(b)
	connection.Shutdown()
}

type echoClient struct {
	client *TCPClient
}

func newEchoClient(addr string) *echoClient {
	echoClient := &echoClient{
		client: NewTCPClient(addr),
	}
	return echoClient
}

func (c *echoClient) ConnectAndServe() {
	c.client.EnableRetry() // 启用retry
	if err := c.client.DialAndServe(c, nil); err != nil {
		log.Println(err)
	}
}

func (c *echoClient) Connect(connection *TCPConnection, connected bool) {
	const message string = "hello"
	if connected {
		log.Printf("echo client: %v connected.\n", connection.RemoteAddr())
		log.Println("echo client: send", message)
		connection.Send([]byte(message))
	} else {
		log.Printf("echo client: %v disconnected.\n", connection.RemoteAddr())
		c.client.Close()
	}
}

func (c *echoClient) Receive(connection *TCPConnection, b []byte) {
	log.Printf("echo client: %v receive %v\n", connection.RemoteAddr(), string(b))
}

func TestEcho(t *testing.T) {
	srv := newEchoServer("localhost:8000")
	go func() {
		c := newEchoClient("localhost:8000")
		c.ConnectAndServe()
		srv.Close()
	}()
	srv.ListenAndServe()
}
