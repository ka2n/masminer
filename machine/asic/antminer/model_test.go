package antminer

import (
	"reflect"
	"testing"

	"github.com/ka2n/masminer/machine"
)

func TestMinerTypeFromString(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    machine.Model
		wantErr bool
	}{
		{"Z9-Mini", args{
			`Sat May 26 20:42:30 CST 2018
Antminer Z9-Mini`,
		}, ModelZ9Mini, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MinerTypeFromString(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("MinerTypeFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MinerTypeFromString() = %v, want %v", got, tt.want)
			}
		})
	}
}
