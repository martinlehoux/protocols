package protocols_test

import (
	"testing"

	"github.com/martinlehoux/protocols"
)

func TestCreateDevice(t *testing.T) {
	protocols.CreateDevice("test device")
}

func TestConnect(t *testing.T) {
	var err error
	device1 := protocols.CreateDevice("device 1")
	device2 := protocols.CreateDevice("device 2")
	if err = protocols.Connect(&device1, &device2); err != nil {
		t.Error("can't connect two devices")
	}
	if err = protocols.Connect(&device1, &device1); err != nil {
		t.Error("can't connect same device")
	}
}

func TestSendPacket(t *testing.T) {
	var err error
	device1 := protocols.CreateDevice("device 1")
	device2 := protocols.CreateDevice("device 2")
	protocols.Connect(&device1, &device2)
	device2.Run()
	if err = device1.SendPacket(device2.MAC, make([]byte, 1)); err == nil {
		t.Error("should not send one byte packet")
	}
	if err = device1.SendPacket(device2.MAC, make([]byte, 1800)); err == nil {
		t.Error("should not send more than 1518 bytes packet")
	}
	if err = device1.SendPacket(device2.MAC, make([]byte, 64)); err != nil {
		t.Error("can't send 64 bytes packet")
	}
}
