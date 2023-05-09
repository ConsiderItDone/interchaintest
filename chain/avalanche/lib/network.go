package lib

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/ava-labs/avalanchego/api/info"
)

func IsOpened(host string, port string) bool {
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), time.Second)
	if err != nil {
		return false
	}

	if conn != nil {
		conn.Close()
		return true
	}

	return false
}

func WaitPort(ctx context.Context, host, port string) error {
	var err error
	for done := false; !done && err == nil; {
		select {
		case <-ctx.Done():
			err = fmt.Errorf("WaitPort(%s, %s) context closed", host, port)
		default:
			done = IsOpened(host, port)
		}
	}
	return err
}

func WaitNode(ctx context.Context, host, port string) error {
	err := WaitPort(ctx, host, port)
	if err != nil {
		return err
	}

	time.Sleep(10 * time.Second)

	client := info.NewClient(fmt.Sprintf("http://%s:%s", host, port))
	for done := false; !done && err == nil; {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context closed")
		default:
			xdone, xerr := client.IsBootstrapped(ctx, "X")
			if errors.Is(err, io.EOF) {
				err = nil
			}
			pdone, perr := client.IsBootstrapped(ctx, "P")
			if errors.Is(err, io.EOF) {
				err = nil
			}
			cdone, cerr := client.IsBootstrapped(ctx, "C")
			if errors.Is(err, io.EOF) {
				err = nil
			}
			done = xdone && pdone && cdone
			err = errors.Join(xerr, perr, cerr)
			time.Sleep(500 * time.Millisecond)
		}
	}

	return err
}
