package main

import (
	"fmt"
	"goruf/platform/core"
	"goruf/platform/tcp"
	"sort"
)

type ServerMessageHandler struct {
	queue []tcp.Msg
}

func NewServerMessageHandler() tcp.MessageHandler {
	return &ServerMessageHandler{
		queue: make([]tcp.Msg, 0),
	}
}

func (s *ServerMessageHandler) Handle(msg tcp.Msg) ([]byte, error) {
	if msg.TotalPage == 1 {
		return s.handle(msg.Payload)
	}
	s.queue = append(s.queue, msg)
	if len(s.queue) < int(msg.TotalPage) {
		return nil, nil
	}
	payload := make([]byte, 0)
	sort.Slice(s.queue, func(i, j int) bool {
		return s.queue[i].Page < s.queue[j].Page
	})
	for _, msg := range s.queue {
		payload = append(payload, msg.Payload...)
	}
	s.queue = make([]tcp.Msg, 0)
	return s.handle(payload)
}

func (s *ServerMessageHandler) handle(b []byte) ([]byte, error) {
	cmd, err := tcp.GetTlv(tcp.TypeCmd, b)
	if err != nil {
		return nil, err
	}
	switch cmd.GetUInt32() {
	case core.CmdConnectReq:
		{
			return tcp.Join(tcp.TlvUInt32(tcp.TypeCmd, core.CmdConnectRep)), nil
		}
	case core.CmdUploadJsReq:
		{
			return tcp.Join(tcp.TlvUInt32(tcp.TypeCmd, core.CmdUploadJsRep)), nil
		}
	case core.CmdUploadCssReq:
		{
			return tcp.Join(tcp.TlvUInt32(tcp.TypeCmd, core.CmdUploadCssRep)), nil
		}
	default:
		return nil, fmt.Errorf("")
	}
}
