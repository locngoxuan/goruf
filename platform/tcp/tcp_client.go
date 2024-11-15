package tcp

import "net"

type TransferData func(conn net.Conn) error

func ConnectAndTransferData(addr string, f TransferData) error {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return err
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return err
	}
	defer conn.Close()
	return f(conn)
}
