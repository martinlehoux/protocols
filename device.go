package protocols

import (
	"bytes"
	"crypto/rand"
	"fmt"
)

// Device represent any hardware device in real life
type Device struct {
	nickname      string
	MAC           []byte
	sendIfnets    [24]chan []byte
	receiveIfnets [24]chan []byte
	macCache      map[int][]byte
}

// Log enable Logging for a device
func (device *Device) Log(str string, a ...interface{}) {
	str = fmt.Sprintf(str, a...)
	fmt.Printf("[%v:%x] %s\n", device.nickname, device.MAC, str)
}

// CreateDevice creates a new device with a nickname and generate its MAC
func CreateDevice(nickname string) Device {
	// Manufacturer MAC prefix
	prefix := []byte{0xff, 0xff, 0xff}
	MAC := make([]byte, 3)
	rand.Read(MAC)
	MAC = append(prefix, MAC...)
	return Device{nickname: nickname, MAC: MAC, macCache: make(map[int][]byte)}
}

// Connect two devices together
func Connect(device1 *Device, device2 *Device) error {
	newIfnet12 := make(chan []byte)
	newIfnet21 := make(chan []byte)
	// There should be available ifnets on both devices
	var port1, port2 int
	var ifnet chan []byte
	for port1, ifnet = range device1.sendIfnets {
		if ifnet == nil {
			break
		}

	}
	for port2, ifnet = range device2.sendIfnets {
		if ifnet == nil {
			break
		}
	}
	if port1 == len(device1.sendIfnets)-1 {
		return fmt.Errorf("no port available for %v (%v ports)", device1.nickname, port1)
	}
	if port2 == len(device2.sendIfnets)-1 {
		return fmt.Errorf("no port available for %v (%v ports)", device2.nickname, port2)
	}
	device1.sendIfnets[port1] = newIfnet12
	device2.receiveIfnets[port2] = newIfnet12
	device2.sendIfnets[port2] = newIfnet21
	device1.receiveIfnets[port1] = newIfnet21
	fmt.Printf("%v:%v <-> %v:%v\n", device1.nickname, port1, device2.nickname, port2)
	return nil
}

func (device *Device) findMACCache(MAC []byte) int {
	for port, mac := range device.macCache {
		if bytes.Equal(mac, MAC) {
			return port
		}
	}
	return -1
}

// SendPacket to a MAC
func (device *Device) SendPacket(to []byte, packetL3 []byte) error {
	device.Log("sending packet to %x", to)
	packetL2, err := L3toL2(packetL3, device.MAC, to)
	if err != nil {
		return err
	}
	// Try to find MAC in MAC table
	if port := device.findMACCache(to); port == -1 {
		device.Log("no cache found for %x", to)
		// If not found, send to all ports
		for port, ifnet := range device.sendIfnets {
			if ifnet != nil {
				device.Log("sending packet on port %v", port)
				ifnet <- packetL2
			}
		}

	} else {
		// If found, send only on port
		device.Log("cache found for MAC:%x", to)
		device.Log("sending packet on port %v", port)
		device.sendIfnets[port] <- packetL2
	}
	return nil
}

// ReceivePacket is the routine for a device receiving a packet
func (device *Device) ReceivePacket(port int, packetL2 []byte) error {
	// Register sender MAC in cache
	packetL3, from, to := L2toL3(packetL2)
	// Check if MAC is own or broadcast
	if bytes.Equal(device.MAC, to) {
		device.Log("packet received from %x: %v bytes", from, len(packetL3))
	} else {
		device.Log("packet dropped from %x", from)
	}
	device.Log("updating cache for port %v: %x", port, from)
	device.macCache[port] = from
	return nil
}

func (device *Device) runPort(port int) {
	ifnet := device.receiveIfnets[port]
	for {
		packet := <-ifnet
		device.ReceivePacket(port, packet)
	}

}

// Run device loop
func (device *Device) Run() {
	for port := range device.receiveIfnets {
		go device.runPort(port)
	}
}
