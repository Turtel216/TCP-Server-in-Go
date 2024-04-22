package main

import (
	"fmt"
	"log"
	"net"
)

type Server struct {
  listenAdrr  string
  ln          net.Listener
  quitch      chan struct{}
}

func NewServer(listenAdrr string) *Server {
  return &Server{
    listenAdrr: listenAdrr,
    quitch:     make(chan struct{}),
  }
}

func (s *Server) Start() error {
  ln, err := net.Listen("tcp", s.listenAdrr)
  if err != nil {
    return err
  }
  defer ln.Close()
  s.ln = ln

  go s.acceptLoop()

  <-s.quitch

  return nil
}

func (s *Server) acceptLoop() {
  for {
    conn, err := s.ln.Accept()
    if err != nil {
      fmt.Println("accept error: ", err)
      continue
    }

    fmt.Println("new connection to the server: ", conn.RemoteAddr())

    go s.readLoop(conn)
  }
}

func (s *Server) readLoop(conn net.Conn) {
  defer conn.Close()

  buf := make([]byte, 2048)
  for{
    n, err := conn.Read(buf)
    if err != nil {
      fmt.Println("read error: ", err)
      continue
    }

    msg := buf[:n]
    fmt.Println(string(msg))
  }
}

func main() {
  server := NewServer(":3000")
  log.Fatal(server.Start())

}
