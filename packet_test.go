package protocols_test

import (
	"bytes"
	"testing"

	"github.com/martinlehoux/protocols"
)

func TestL3toL2toL3(t *testing.T) {
	var packetL3, packetL2 []byte
	var err error
	from := []byte{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}
	to := []byte{0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa}
	packetL3 = make([]byte, 5)
	packetL2, err = protocols.L3toL2(packetL3, from, to)
	if err == nil {
		t.Error("should not send one byte packet")
	}
	packetL3 = make([]byte, 1800)
	packetL2, err = protocols.L3toL2(packetL3, from, to)
	if err == nil {
		t.Error("should not send more than 1518 bytes packet")
	}
	packetL3 = make([]byte, 64)
	packetL2, err = protocols.L3toL2(packetL3, from, to)
	if err != nil {
		t.Error("can't encapsulate package of 64 bytes")
	}
	PacketL3, From, To := protocols.L2toL3(packetL2)
	if !bytes.Equal(packetL3, PacketL3) {
		t.Error("packet L3 has changed")
	}
	if !bytes.Equal(from, From) {
		t.Error("from address has changed")
	}
	if !bytes.Equal(to, To) {
		t.Error("to address has changed")
	}
}
