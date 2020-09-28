package main

import (
	"github.com/izhw/gnet"
	"github.com/izhw/gnet/gcore"
	"github.com/izhw/gnet/logger"
)

func main() {
	log := logger.GlobalSimpleLogger()
	svc := gnet.NewService(
		gcore.WithServiceType(gcore.ServiceTCPServer),
		gcore.WithAddr("0.0.0.0:7777"),
		gcore.WithEventHandler(NewServerHandler()),
		gcore.WithLogger(log),
	)

	s := svc.Server()
	if err := s.Init(); err != nil {
		log.Fatal("server init error:", err)
	}
	log.Fatal("Exit:", s.Serve())
}
