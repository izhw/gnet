package main

import (
	"github.com/izhw/gnet"
	"github.com/izhw/gnet/logger"
)

type AsyncHandler struct {
	*gnet.NetEventHandler
	logger logger.Logger
}

func NewAsyncHandler() *AsyncHandler {
	return &AsyncHandler{
		logger: logger.NewSimpleLogger(),
	}
}

func (h *AsyncHandler) OnReadMsg(c gnet.Conn, data []byte) error {
	h.logger.Info(c.GetTag(), "read msg:", string(data))
	return nil
}

func (h *AsyncHandler) OnWriteError(c gnet.Conn, data []byte, err error) {
	h.logger.Warn(c.GetTag(), "write error:", err)
}
