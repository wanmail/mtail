// Copyright 2020 Google Inc. All Rights Reserved.
// This file is available under the Apache license.

package logstream_test

import (
	"context"
	"net"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/google/mtail/internal/logline"
	"github.com/google/mtail/internal/tailer/logstream"
	"github.com/google/mtail/internal/testutil"
	"github.com/google/mtail/internal/waker"
)

const dgramTimeout = 1 * time.Second

func TestDgramStreamReadCompletedBecauseSocketClosed(t *testing.T) {
	testutil.TimeoutTest(dgramTimeout, func(t *testing.T) {
		var wg sync.WaitGroup

		tmpDir := testutil.TestTempDir(t)

		name := filepath.Join(tmpDir, "sock")

		lines := make(chan *logline.LogLine, 1)
		ctx, cancel := context.WithCancel(context.Background())
		waker, awaken := waker.NewTest(ctx, 1)

		sockName := "unixgram://" + name
		ss, err := logstream.New(ctx, &wg, waker, sockName, lines, false)
		testutil.FatalIfErr(t, err)

		s, err := net.DialUnix("unixgram", nil, &net.UnixAddr{name, "unixgram"})
		testutil.FatalIfErr(t, err)

		_, err = s.Write([]byte("1\n"))
		testutil.FatalIfErr(t, err)

		awaken(0) // sync past read

		ss.Stop()

		// "Close" the socket by sending zero bytes, which after Stop tells the stream to act as if we're done.
		_, err = s.Write([]byte{})
		testutil.FatalIfErr(t, err)

		wg.Wait()
		close(lines)

		received := testutil.LinesReceived(lines)
		expected := []*logline.LogLine{
			{context.TODO(), name, "1"},
		}
		testutil.ExpectNoDiff(t, expected, received, testutil.IgnoreFields(logline.LogLine{}, "Context"))

		cancel()
		wg.Wait()

		if !ss.IsComplete() {
			t.Errorf("expecting dgramstream to be complete because socket closed")
		}
	})(t)
}

func TestDgramStreamReadCompletedBecauseCancel(t *testing.T) {
	testutil.TimeoutTest(dgramTimeout, func(t *testing.T) {
		var wg sync.WaitGroup

		tmpDir := testutil.TestTempDir(t)

		name := filepath.Join(tmpDir, "sock")

		lines := make(chan *logline.LogLine, 1)
		ctx, cancel := context.WithCancel(context.Background())
		waker, awaken := waker.NewTest(ctx, 1)

		sockName := "unixgram://" + name
		ss, err := logstream.New(ctx, &wg, waker, sockName, lines, false)
		testutil.FatalIfErr(t, err)

		s, err := net.DialUnix("unixgram", nil, &net.UnixAddr{name, "unixgram"})
		testutil.FatalIfErr(t, err)

		_, err = s.Write([]byte("1\n"))
		testutil.FatalIfErr(t, err)

		awaken(0) // Synchronise past read.

		cancel() // This cancellation should cause the stream to shut down.

		wg.Wait()
		close(lines)

		received := testutil.LinesReceived(lines)
		expected := []*logline.LogLine{
			{context.TODO(), name, "1"},
		}
		testutil.ExpectNoDiff(t, expected, received, testutil.IgnoreFields(logline.LogLine{}, "Context"))

		if !ss.IsComplete() {
			t.Errorf("expecting dgramstream to be complete because cancel")
		}
	})(t)
}