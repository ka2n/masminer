package net

import (
	"testing"
)

func Test_ParseIPAddr(t *testing.T) {
	type args struct {
		in string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"ipv4", args{"192.168.10.5"}, "192.168.10.5", false},
		{"w/ CIDR", args{"192.168.10.5/24"}, "192.168.10.5", false},
		{"ipv6", args{"fe80::3735:6e27:473b:509f"}, "fe80::3735:6e27:473b:509f", false},
		{"ipv6 w/ CIDR", args{"fe80::3735:6e27:473b:509f/64"}, "fe80::3735:6e27:473b:509f", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseIPAddr(tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseIPAddr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseIPAddr() = %v, want %v", got, tt.want)
			}
		})
	}
}
