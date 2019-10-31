package protocols

import (
	"fmt"
	"strconv"
	"strings"
)

// IPv6 represents an IPv6 encoded in bytes
type IPv6 [8]int

// IPv6LocalHost is the localhost ipv6 address for a device
var IPv6LocalHost IPv6 = IPv6{0, 0, 0, 0, 0, 0, 0, 1}

func StringToIPv6(str string) (IPv6, error) {
	strs := strings.Split(str, ":")
	ipv6 := new(IPv6)
	if len(strs) == 8 {
		for i, str := range strs {
			integer, _ := strconv.ParseInt(str, 16, 0)
			ipv6[i] = int(integer)
		}
		return *ipv6, nil
	}
	return *ipv6, fmt.Errorf("can't parse ipv6: %v", str)
}
