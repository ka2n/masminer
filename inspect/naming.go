package inspect

import "strings"

// ShortName : Get short name from macAddr
func ShortName(macAddr string) string {
	segs := strings.Split(macAddr, ":")
	return strings.Join(segs[len(segs)-3:len(segs)], "")
}
