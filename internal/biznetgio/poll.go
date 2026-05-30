package biznetgio

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const retryAttempts = 3

func retryable(err error) bool {
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		return false
	}
	switch apiErr.StatusCode {
	case http.StatusTooManyRequests, http.StatusInternalServerError,
		http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		return true
	}
	return false
}

func retryBackoff(attempt int) time.Duration {
	switch attempt {
	case 1:
		return 500 * time.Millisecond
	case 2:
		return time.Second
	default:
		return 2 * time.Second
	}
}

func sleepWithCtx(ctx context.Context, d time.Duration) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(d):
		return nil
	}
}

// WaitForStatus polls get until status(t) matches a terminal value.
// interval: how often to poll. If ctx has no deadline, a 30m timeout is applied.
func WaitForStatus[T any](
	ctx context.Context,
	interval time.Duration,
	get func(ctx context.Context) (T, error),
	status func(T) string,
	ready, failed []string,
) (T, error) {
	var zero T
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 30*time.Minute)
		defer cancel()
	}
	if interval <= 0 {
		interval = 5 * time.Second
	}
	for {
		v, err := get(ctx)
		if err != nil {
			return zero, err
		}
		s := status(v)
		if contains(ready, s) {
			return v, nil
		}
		if contains(failed, s) {
			return zero, fmt.Errorf("resource reached failed status %q", s)
		}
		select {
		case <-ctx.Done():
			return zero, ctx.Err()
		case <-time.After(interval):
		}
	}
}

func contains(xs []string, v string) bool {
	for _, x := range xs {
		if x == v {
			return true
		}
	}
	return false
}
// wip 436
