package core

import "mfe-platform/tcp"

const (
	CmdConnectReq uint32 = iota
	CmdConnectRep
	CmdUploadJsReq
	CmdUploadJsRep
	CmdUploadCssReq
	CmdUploadCssRep
)

type CmdConnect struct {
	Cmd     uint32
	Payload []byte
}

func (c CmdConnect) Pack() []byte {
	return tcp.Join(
		tcp.TlvUInt32(tcp.TypeCmd, c.Cmd),
		tcp.NewTlv(tcp.TypePayload, c.Payload),
	)
}
