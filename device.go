package protocols

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"time"
)

// Device represent any hardware device in real life
type Device struct {
	nickname string
	MAC      []byte
	links    [24]chan []byte
	macCache map[int][]byte
}

// Log enable Logging for a device
func (device *Device) Log(str string, a ...interface{}) {
	str = fmt.Sprintf(str, a...)
	fmt.Printf("[%v] %s\n", device.nickname, str)
}

// CreateDevice creates a new device with a nickname and generate its MAC
func CreateDevice(nickname string) Device {
	// Manufacturer MAC prefix
	prefix := []byte{0xDA, 0xDA, 0xBE}
	MAC := make([]byte, 3)
	rand.Read(MAC)
	MAC = append(prefix, MAC...)
	return Device{nickname: nickname, MAC: MAC, macCache: make(map[int][]byte)}
}

// Connect two devices together
func Connect(device1 *Device, device2 *Device) error {
	newLink := make(chan []byte)
	// There should be available links on both devices
	var port1, port2 int
	for port1, link := range device1.links {
		if link == nil {
			break
		}
		return fmt.Errorf("no port available for %v (%v ports)", device1.nickname, port1)
	}
	for port2, link := range device2.links {
		if link == nil {
			break
		}
		return fmt.Errorf("no port available for %v (%v ports)", device2.nickname, port2)
	}
	device1.links[port1] = newLink
	device2.links[port2] = newLink
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
	packetL2, err := L3toL2(packetL3, device.MAC, to)
	if err != nil {
		panic(err)
	}
	// Try to find MAC in MAC table
	if port := device.findMACCache(to); port == -1 {
		device.Log("no cache found for %x", to)
		// If not found, send to all ports
		for _, link := range device.links {
			link <- packetL2
		}

	} else {
		// If found, send only on port
		device.Log("cache found for MAC:%x : %v", to, port)
		device.links[port] <- packetL2
	}
	return nil
}

// ReceivePacket is the routine for a device receiving a packet
func (device *Device) ReceivePacket(port int, packetL2 []byte) error {
	// Register sender MAC in cache
	packetL3, from, to := L2toL3(packetL2)
	device.macCache[port] = from
	fmt.Println(from)
	// Check if MAC is own or broadcast
	if bytes.Equal(device.MAC, to) {
		device.Log("packet received from %x : %x", from, packetL3)
	} else {
		device.Log("packet dropped from %x", from)
	}
	return nil
}

func (device *Device) isPacketForMe(packet []byte) bool {
	return bytes.Equal(packet[0:6], device.MAC)
}

// Run device loop
func (device *Device) Run() {
	for {
		time.Sleep(10 * time.Millisecond)
		for port, link := range device.links {
			go func(link chan []byte, port int) {
				packet := <-link
				device.ReceivePacket(port, packet)
			}(link, port)
		}
	}
}
