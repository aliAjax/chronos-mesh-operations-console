package platform

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func SignalContext(parent context.Context) (context.Context, context.CancelFunc) {
	return signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
}
