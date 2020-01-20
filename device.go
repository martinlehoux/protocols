package protocols

// BUG(kagamino) Connections created after device.Run() don't work

import (
	"bytes"
	"crypto/rand"
	"fmt"
)

// Ifnet represent a network interface of a device
type Ifnet struct {
	channel   chan []byte
	IPAddress IPv6
}

// Device represent any hardware device in real life
type Device struct {
	nickname      string
	MAC           []byte
	sendIfnets    [24]Ifnet
	receiveIfnets [24]Ifnet
	macCache      map[int][]byte
}

// Log enable Logging for a device
func (device *Device) Log(str string, a ...interface{}) {
	str = fmt.Sprintf(str, a...)
	fmt.Printf("[%v:%x] %s\n", device.nickname, device.MAC, str)
}

// CreateDevice creates a new device with a nickname and generate its MAC and localhost ifnet
func CreateDevice(nickname string) Device {
	// Manufacturer MAC prefix
	prefix := []byte{0xff, 0xff, 0xff}
	MAC := make([]byte, 3)
	rand.Read(MAC)
	MAC = append(prefix, MAC...)
	device := Device{nickname: nickname, MAC: MAC, macCache: make(map[int][]byte)}
	// Connect device to itself
	Connect(&device, &device)
	// Connected ifnet is index 0
	device.sendIfnets[0].IPAddress = IPv6LocalHost
	device.receiveIfnets[0].IPAddress = IPv6LocalHost
	return device
}

func (device *Device) PrintIfnetsTable() {
	// Send and receive should be the same
	// TODO: one array but two channels
	fmt.Printf("INTERFACES TABLES: %v\t%x\n", device.nickname, device.MAC)
	for i, ifnet := range device.sendIfnets {
		if ifnet.channel != nil {
			fmt.Printf("%v\t%v\t\n", i, ifnet.IPAddress)
		}
	}
}

func (device *Device) PrintMACCacheTable() {
	fmt.Printf("MAC CACHE TABLES: %v\n", device.nickname)
	for index, mac := range device.macCache {
		fmt.Printf("%v\t%x\n", index, mac)
	}
}

// Connect two devices together
func Connect(device1 *Device, device2 *Device) error {
	var chan12, chan21 = make(chan []byte), make(chan []byte)
	var newIfnet12, newIfnet21 = Ifnet{channel: chan12}, Ifnet{channel: chan21}
	// There should be available ifnets on both devices
	var index1, index2 int
	var ifnet Ifnet
	for index1, ifnet = range device1.sendIfnets {
		if ifnet.channel == nil {
			break
		}

	}
	for index2, ifnet = range device2.sendIfnets {
		if ifnet.channel == nil {
			break
		}
	}
	if index1 == len(device1.sendIfnets)-1 {
		return fmt.Errorf("no interface available for %v (%v interfaces)", device1.nickname, index1)
	}
	if index2 == len(device2.sendIfnets)-1 {
		return fmt.Errorf("no interface available for %v (%v interfaces)", device2.nickname, index2)
	}
	device1.sendIfnets[index1] = newIfnet12
	device2.receiveIfnets[index2] = newIfnet12
	device2.sendIfnets[index2] = newIfnet21
	device1.receiveIfnets[index1] = newIfnet21
	if !bytes.Equal(device1.MAC, device2.MAC) {
		fmt.Printf("%v:%v <-> %v:%v\n", device1.nickname, index1, device2.nickname, index2)
	}
	return nil
}

func (device *Device) findMACCache(MAC []byte) int {
	for index, mac := range device.macCache {
		if bytes.Equal(mac, MAC) {
			return index
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
	if index := device.findMACCache(to); index == -1 {
		device.Log("no cache found for %x", to)
		// If not found, send to all in
		for index, ifnet := range device.sendIfnets {
			if ifnet.channel != nil {
				device.Log("sending packet on interface %v", index)
				ifnet.channel <- packetL2
			}
		}

	} else {
		// If found, send only on interface
		device.Log("cache found for MAC:%x", to)
		device.Log("sending packet on interface %v", index)
		device.sendIfnets[index].channel <- packetL2
	}
	return nil
}

// ReceivePacket is the routine for a device receiving a packet
func (device *Device) ReceivePacket(index int, packetL2 []byte) error {
	// Register sender MAC in cache
	packetL3, from, to := L2toL3(packetL2)
	// Check if MAC is own or broadcast
	if bytes.Equal(device.MAC, to) {
		device.Log("packet received from %x: %v bytes", from, len(packetL3))
	} else {
		device.Log("packet dropped from %x", from)
	}
	device.Log("updating cache for interface %v: %x", index, from)
	device.macCache[index] = from
	return nil
}

func (device *Device) runInterface(index int) {
	ifnet := device.receiveIfnets[index]
	for {
		packet := <-ifnet.channel
		device.ReceivePacket(index, packet)
	}

}

// Run device loop
func (device *Device) Run() {
	for index := range device.receiveIfnets {
		go device.runInterface(index)
	}
}
