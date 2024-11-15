package tcp

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/rs/zerolog/log"
)

type MessageHandler interface {
	Handle(msg Msg) ([]byte, error)
}

type HandlerCreator func() MessageHandler

func OpenListener(port int, creator HandlerCreator) error {
	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		return err
	}
	defer l.Close()
	log.Info().Int("port", port).Msg("")
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Error().Err(err).Msg("failed to accept connection")
			continue
		}
		client := ClientConn{
			conn:    conn,
			handler: creator(),
			buffer:  make(chan []byte, 10),
		}
		go client.handleRequest()
	}
}

type ClientConn struct {
	sync.Mutex
	conn    net.Conn
	handler MessageHandler
	buffer  chan []byte
}

func (c *ClientConn) send(b []byte) {
	c.Lock()
	defer c.Unlock()
	c.buffer <- b
}

func (c *ClientConn) handleRequest() {
	r := bufio.NewReader(c.conn)
	w := bufio.NewWriter(c.conn)
	defer func() {
		c.conn.Close()
		close(c.buffer)
	}()
	for {
		msg, err := Read(r)
		if err != nil {
			if errors.Is(io.EOF, err) {
				log.Info().Msg("connection from client has been closed")
				break
			}
			log.Error().Err(err).Msg("failed to read message from connection")
			break
		}
		resp, err := c.handler.Handle(msg)
		if err != nil {
			log.Error().Err(err).Msg("failed to handle incoming message")
			break
		}
		if resp == nil {
			continue
		}
		_, err = w.Write(resp)
		if err != nil {
			log.Error().Err(err).Msg("failed to send msg over tcp")
			break
		}
		err = w.Flush()
		if err != nil {
			log.Error().Err(err).Msg("failed to flush msg over tcp")
			break
		}
	}
}
