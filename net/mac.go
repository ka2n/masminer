package net

import (
	"net"
)

// ValidateMAC validates a address is hardware address
func ValidateMAC(mac string) bool {
	addr, err := net.ParseMAC(mac)
	if err != nil {
		return false
	}

	isBroadcast := (addr[0] >> 0 & 1) == 1
	if isBroadcast {
		return false
	}

	// Check if all zero: 00:00:00:00:00:00
	for _, b := range addr {
		if b != 0 {
			return true
		}
	}
	return false
}
