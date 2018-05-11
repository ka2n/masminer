package machine

import "strings"

// ShortName : Get short name from macAddr
func ShortName(macAddr string) string {
	segs := strings.Split(macAddr, ":")
	if len(segs) < 4 {
		return ""
	}
	return strings.Join(segs[len(segs)-3:len(segs)], "")
}
