package main

import (
	"log"

	"github.com/izhw/gnet/tcp/tcpserver"
)

func main() {
	s := tcpserver.NewServer("0.0.0.0:7777", NewServerHandler())
	log.Fatal("Exit:", s.Serve())
}
