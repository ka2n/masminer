package base

import (
	"context"
	"net"
	"strings"
	"testing"
	"time"
)

type dialer struct {
}

func (d dialer) Dial(network, address string) (net.Conn, error) {
	return net.Dial(network, address)
}

func TestOutputMinerRPC(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var l net.Listener
	var err error

	connected := make(chan struct{})
	go func() {
		l, err = net.Listen("tcp", ":4028")
		if err != nil {
			t.Fatal(err)
		}
		connected <- struct{}{}

		for {
			c, err := l.Accept()
			if err != nil {
				return
			}

			go func() {
				defer c.Close()

				buf := make([]byte, 1024)
				c.Read(buf)

				input := string(buf)
				t.Log("req: ", input)

				if strings.Contains(input, "timeout") {
					t.Log("Make timeout")
					time.Sleep(time.Millisecond * 500)
				}

				t.Log("Write response")
				c.Write([]byte("OK\x00"))
				t.Log("Writen response")
			}()
		}
	}()

	<-connected
	defer l.Close()

	t.Run("requests", func(t *testing.T) {
		t.Run("nromal", func(t *testing.T) {
			ctx, cancel := context.WithTimeout(ctx, time.Millisecond*10000)
			defer cancel()

			ret, err := OutputMinerRPC(ctx, dialer{}, "test", "")
			t.Log(string(ret))
			if err != nil {
				t.Fatal(err)
			}
		})

		t.Run("timeout", func(t *testing.T) {
			ctx, cancel := context.WithTimeout(ctx, time.Millisecond*100)
			defer cancel()

			_, err = OutputMinerRPC(ctx, dialer{}, "timeout", "")
			t.Log(err)
			if err == nil {
				t.Fatal("Expected to timeout error")
			}
		})
	})
}
