package baikalver

import (
	"strconv"
)

// VersionFromFWV returns Firmware version from 'FWV' value
func VersionFromFWV(s string) (string, error) {
	fwv, err := strconv.Atoi(s)
	if err != nil {
		return "", err
	}
	fwvn := uint8(fwv)
	return strconv.Itoa(int(fwvn>>4&255)) + "." + strconv.Itoa(int(fwvn&15)), nil
}
