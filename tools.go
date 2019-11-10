package partybot

import (
	"context"
	"time"
)

func SleepCanBreak(ctx context.Context, sleep float64) (isBreak bool) {
	select {
	case <-ctx.Done():
		isBreak = true
	case <-time.After(time.Duration(sleep*1000) * time.Millisecond):
		isBreak = false
	}
	return
}
