package checks

import (
	"context"
	"net"
	"time"

	"github.com/pkg/errors"
)

// PingContext verifies a resource is still alive.
// This would normally be a TCP dial check, a db.PingContext() or something similar.
type Pinger interface {
	PingContext(ctx context.Context) error
}

// PingContextFunc type is an adapter to allow the use of ordinary functions as Pingers.
type PingContextFunc func(ctx context.Context) error

// PingContext calls f(ctx).
func (f PingContextFunc) PingContext(ctx context.Context) error {
	return f(ctx)
}

// NewPingCheck returns a Check that pings using the specified Pinger and fails on timeout or ping failure
func NewPingCheck(name string, pinger Pinger, timeout time.Duration) (Check, error) {
	if pinger == nil {
		return nil, errors.New("Pinger must not be nil")
	}

	return &CustomCheck{
		CheckName: name,
		CheckFunc: func() (details interface{}, err error) {
			pingCtx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()
			return nil, pinger.PingContext(pingCtx)
		},
	}, nil
}

// NewDialPinger returns a Pinger that pings the specified address
func NewDialPinger(network, address string) PingContextFunc {
	var d net.Dialer
	return func(ctx context.Context) error {
		conn, err := d.DialContext(ctx, network, address)
		if err == nil {
			_ = conn.Close()
		}

		return err
	}
}
