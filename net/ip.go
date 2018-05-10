package net

import (
	"fmt"
	"net"
	"strings"
)

func ParseIPAddr(in string) (string, error) {
	if strings.Contains(in, "/") {
		ip, _, err := net.ParseCIDR(in)
		if err != nil {
			return "", err
		}
		if ip != nil {
			return ip.String(), nil
		}
	} else {
		ip := net.ParseIP(in)
		if ip != nil {
			return ip.String(), nil
		}
	}
	return "", fmt.Errorf("parse error ip: %s", in)
}
