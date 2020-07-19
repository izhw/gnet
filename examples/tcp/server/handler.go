package main

import (
	"github.com/izhw/gnet"
	"github.com/izhw/gnet/logger"
)

type ServerHandler struct {
	logger logger.Logger
}

func NewServerHandler() *ServerHandler {
	return &ServerHandler{
		logger: logger.NewSimpleLogger(),
	}
}

func (h *ServerHandler) OnOpened(c gnet.Conn) {
	h.logger.Info(c.RemoteAddr(), "opened")
}

func (h *ServerHandler) OnClosed(c gnet.Conn) {
	h.logger.Info(c.RemoteAddr(), "closed")
}

func (h *ServerHandler) OnReadMsg(c gnet.Conn, data []byte) error {
	h.logger.Info(c.RemoteAddr(), "Server read msg:", string(data))
	if err := c.Write(data); err != nil {
		h.logger.Error("write error:", err)
		return err
	}
	return nil
}

func (h *ServerHandler) OnWriteError(c gnet.Conn, data []byte, err error) {
	h.logger.Warn(c.RemoteAddr(), "data:", string(data), "write error:", err)
}
