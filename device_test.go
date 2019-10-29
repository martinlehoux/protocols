package protocols_test

import (
	"protocols"
	"testing"
)

func TestCreateDevice(t *testing.T) {
	protocols.CreateDevice("test device")
}

func Test(t *testing.T) {
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
