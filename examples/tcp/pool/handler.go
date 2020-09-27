package main

import (
	"github.com/izhw/gnet/gcore"
	"github.com/izhw/gnet/logger"
)

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
	h.logger.Info(c.GetTag(), "AsyncPool read msg:", string(data))
	return nil
}

func (h *AsyncHandler) OnWriteError(c gcore.Conn, data []byte, err error) {
	h.logger.Warn(c.GetTag(), "AsyncPool write error:", err)
}
