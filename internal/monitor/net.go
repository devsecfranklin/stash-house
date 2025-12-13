package net

import (
    "errors"
    "net"
    "strings"
)

// Server defines the minimum contract our
// TCP and UDP server implementations must satisfy.
type Server interface {
    Run() error
    Close() error
}

// NewServer creates a new Server using given protocol
// and addr.
func NewServer(protocol, addr string) (Server, error) {
    switch strings.ToLower(protocol) {
    case "tcp":
        return &TCPServer{
            addr: addr,
        }, nil
    case "udp":
    }
    return nil, errors.New("Invalid protocol given")
}

// TCPServer holds the structure of our TCP
// implementation.
type TCPServer struct {
    addr   string
    server net.Listener
}

// Run starts the TCP Server.
func (t *TCPServer) Run() (err error) {
    t.server, err = net.Listen("tcp", t.addr)
    if err != nil {
        return
    }
    for {
        conn, err := t.server.Accept()
        if err != nil {
            err = errors.New("could not accept connection")
            break
        }
        if conn == nil {
            err = errors.New("could not create connection")
            break
        }
        conn.Close()
    }
    return
}

// Close shuts down the TCP Server
func (t *TCPServer) Close() (err error) {
    return t.server.Close()
}