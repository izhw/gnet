package main

import (
	"github.com/izhw/gnet/gcore"
	"github.com/izhw/gnet/logger"
)

// server handler
type ServerHandler struct {
	logger logger.Logger
}

func NewServerHandler() *ServerHandler {
	return &ServerHandler{
		logger: logger.GlobalSimpleLogger(),
	}
}

func (h *ServerHandler) OnOpened(c gcore.Conn) {
	h.logger.Info(c.RemoteAddr(), "opened")
}

func (h *ServerHandler) OnClosed(c gcore.Conn) {
	h.logger.Info(c.RemoteAddr(), "closed")
}

func (h *ServerHandler) OnReadMsg(c gcore.Conn, data []byte) error {
	h.logger.Info(c.RemoteAddr(), "Server read msg:", string(data))
	if err := c.Write(data); err != nil {
		h.logger.Error("write error:", err)
		return err
	}
	return nil
}

func (h *ServerHandler) OnWriteError(c gcore.Conn, data []byte, err error) {
	h.logger.Warn(c.RemoteAddr(), "data:", string(data), "write error:", err)
}

// client handler
type AsyncHandler struct {
	*gcore.NetEventHandler
	logger logger.Logger
}

func NewAsyncHandler() *AsyncHandler {
	return &AsyncHandler{
		logger: logger.GlobalSimpleLogger(),
	}
}

func (h *AsyncHandler) OnReadMsg(c gcore.Conn, data []byte) error {
	h.logger.Info(c.GetTag(), "AsyncClient read msg:", string(data))
	return nil
}

func (h *AsyncHandler) OnWriteError(c gcore.Conn, data []byte, err error) {
	h.logger.Warn(c.GetTag(), "AsyncClient write error:", err)
}
