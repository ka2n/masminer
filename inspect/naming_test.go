package inspect

import "testing"

func TestShortName(t *testing.T) {
	type args struct {
		macAddr string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "OK", args: args{macAddr: "00:0a:95:9d:68:16"}, want: "9d6816"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ShortName(tt.args.macAddr); got != tt.want {
				t.Errorf("ShortName() = %v, want %v", got, tt.want)
			}
		})
	}
}
