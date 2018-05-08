package metal

import (
	"bytes"
	"net"
)

// ActiveAddress : 有効なInterfaceのMACアドレスとIPアドレスを取得する
func ActiveAddress() (mac string, ip string) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return
	}
	for _, i := range interfaces {
		if i.Flags&net.FlagUp != 0 && bytes.Compare(i.HardwareAddr, nil) != 0 {
			mac := i.HardwareAddr.String()

			addrs, err := i.Addrs()
			if err != nil {
				continue
			}
			for _, addr := range addrs {
				var ipValue net.IP
				switch v := addr.(type) {
				case *net.IPAddr:
					ipValue = v.IP
				case *net.IPNet:
					ipValue = v.IP
				}
				if ipValue == nil || ipValue.IsLoopback() {
					continue
				}
				vv := ipValue.To4()
				if vv == nil {
					continue
				}
				ip := vv.String()
				return mac, ip
			}
		}
	}
	return
}
