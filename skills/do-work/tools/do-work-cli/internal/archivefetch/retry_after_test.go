package archivefetch

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestDownloadAtomicWaitsForRetryAfterBeforeRetrying(t *testing.T) {
	setAtomicRetryTiming(t, time.Millisecond, 3*time.Second)
	readyAt := time.Now().Add(time.Second)
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if time.Now().Before(readyAt) {
			response.Header().Set("Retry-After", "1")
			response.WriteHeader(http.StatusTooManyRequests)
			return
		}
		_, _ = response.Write([]byte("downloaded"))
	}))
	defer server.Close()
	result := DownloadAtomic(context.Background(), server.URL, filepath.Join(t.TempDir(), "download"))
	if result.Err != nil || result.Attempts != 2 || result.BytesWritten != 10 {
		t.Fatalf("server-directed delay was ignored: %+v", result)
	}
}

func TestDownloadAtomicRetryAfterCannotOutrunBudgetOrCancellation(t *testing.T) {
	for _, cancel := range []bool{false, true} {
		t.Run(fmt.Sprintf("cancellation=%t", cancel), func(t *testing.T) {
			setAtomicRetryTiming(t, time.Millisecond, 2*time.Second)
			var calls atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
				calls.Add(1)
				response.Header().Set("Retry-After", "1")
				if !cancel {
					response.Header().Set("Retry-After", "8")
				}
				response.WriteHeader(http.StatusServiceUnavailable)
			}))
			defer server.Close()
			ctx := context.Background()
			if cancel {
				var stop context.CancelFunc
				ctx, stop = context.WithTimeout(ctx, 50*time.Millisecond)
				defer stop()
			}
			directory := t.TempDir()
			result := DownloadAtomic(ctx, server.URL, filepath.Join(directory, "download"))
			if result.Err == nil || result.Attempts != 1 || calls.Load() != 1 {
				t.Fatalf("early retries bypassed delay: %+v calls=%d", result, calls.Load())
			}
			if cancel && !errors.Is(result.Err, context.DeadlineExceeded) {
				t.Fatalf("cancellation lost: %v", result.Err)
			}
			if !cancel && !strings.Contains(result.Err.Error(), "retry budget exhausted") {
				t.Fatalf("budget failure lost: %v", result.Err)
			}
			entries, err := os.ReadDir(directory)
			if err != nil || len(entries) != 0 {
				t.Fatalf("failed download left files: %v, %v", entries, err)
			}
		})
	}
}

func TestRetryAfterDelayParsesDatesAndRetainsFallback(t *testing.T) {
	setAtomicRetryTiming(t, 2*time.Second, time.Minute)
	now := time.Date(2026, 9, 16, 12, 0, 0, 0, time.UTC)
	for _, test := range []struct {
		header string
		want   time.Duration
	}{
		{"8", 8 * time.Second},
		{now.Add(8 * time.Second).Format(http.TimeFormat), 8 * time.Second},
		{now.Add(-8 * time.Second).Format(http.TimeFormat), 2 * time.Second},
		{"1", 2 * time.Second},
		{"0", 2 * time.Second},
		{"", 2 * time.Second},
		{"invalid", 2 * time.Second},
		{"-3", 2 * time.Second},
		{"1.5", 2 * time.Second},
		{"9223372036854775807", time.Duration(1<<63 - 1)},
		{"9999999999999999999999999", time.Duration(1<<63 - 1)},
	} {
		t.Run(test.header, func(t *testing.T) {
			if got := retryAfterDelay(test.header, now); got != test.want {
				t.Fatalf("Retry-After %q: got %s, want %s", test.header, got, test.want)
			}
		})
	}
}
