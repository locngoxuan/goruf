package tcp

import (
	"bufio"
	"bytes"
	"testing"
)

func TestRead(t *testing.T) {
	type Object struct {
		Cmd string `json:"cmd"`
		Msg string `json:"message"`
	}

	data, err := Pack(Join(), 1024)
	if err != nil {
		t.Logf("failed to Pack(interface{}, size): %v", err)
		t.FailNow()
	}
	if len(data) != 1 {
		t.Logf("size of data should be 1")
		t.FailNow()
	}
	r := bufio.NewReader(bytes.NewBuffer(data[0]))
	msg, err := Read(r)
	if err != nil {
		t.Logf("failed to Read(r bufio.Reader): %v", err)
		t.FailNow()
	}
	if msg.Version != Version {
		t.Logf("Versions should be 1, actual = %v", msg.Version)
		t.FailNow()
	}
}

func TestReadMultiple(t *testing.T) {
	data, err := Pack(Join(
		TlvString(0, "This is a command"),
		TlvString(1, "This is a message"),
	), 10)
	if err != nil {
		t.Logf("failed to Pack(interface{}, size): %v", err)
		t.FailNow()
	}
	totalMsg := 5
	if len(data) != totalMsg {
		t.Logf("size of data should be 6, actual = %d", len(data))
		t.FailNow()
	}
	msgs := make([]Msg, 0)
	for i := 0; i < totalMsg; i++ {
		r := bufio.NewReader(bytes.NewBuffer(data[i]))
		msg, err := Read(r)
		if err != nil {
			t.Logf("failed to Read(r bufio.Reader): %v", err)
			t.FailNow()
		}
		msgs = append(msgs, msg)
	}

	payload := make([]byte, 0)
	for _, msg := range msgs {
		payload = append(payload, msg.Payload...)
	}
	if v, err := GetTlv(0, payload); err != nil || v.GetString() != "This is a command" {
		t.Logf("command is not correct, expected = This is a command, actual = %v", v.GetString())
		t.FailNow()
	}
	if v, err := GetTlv(1, payload); err != nil || v.GetString() != "This is a message" {
		t.Logf("message is not correct, expected = This is a message, actual = %v", v.GetString())
		t.FailNow()
	}
}
