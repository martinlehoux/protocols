package protocols

import "fmt"

// L3toL2 encapsulates a packet from the Network(3) to the Datalink(2) layer
func L3toL2(packetL3 []byte, from []byte, to []byte) ([]byte, error) {
	var packetL2 []byte
	if len(packetL3) > 1518 {
		return nil, fmt.Errorf("packet size exceeds 1518: %d", len(packetL3))
	}
	if len(packetL3) < 64 {
		return nil, fmt.Errorf("packet size is less than 64: %d", len(packetL3))
	}
	packetL2 = append([]byte{8, 0}, packetL3...)
	packetL2 = append(from, packetL3...)
	packetL2 = append(to, packetL2...)
	return packetL2, nil
}

// L2toL3 decapsulate a packet from the Datalink(2) to the Network(3) layer
func L2toL3(packetL2 []byte) (packetL3 []byte, from []byte, to []byte) {
	to = packetL2[0:6]
	from = packetL2[6:12]
	packetL3 = packetL2[14:]
	return
}
