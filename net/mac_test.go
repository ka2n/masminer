package net

import "testing"

func TestValidateMAC(t *testing.T) {
	type args struct {
		mac string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"Empty", args{""}, false},
		{"Zero", args{"00:00:00:00:00:00"}, false},
		{"Broadcast", args{"FF:FF:FF:FF:FF:FF"}, false},
		{"IPv4 multicast", args{"01:00:5E:00:00:00"}, false},
		{"IPv6 multicast", args{"33:33:00:00:00:00"}, false},
		{"Valid example mac address", args{"32:61:3C:4E:B6:05"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateMAC(tt.args.mac); got != tt.want {
				t.Errorf("ValidateMAC() = %v, want %v", got, tt.want)
			}
		})
	}
}
