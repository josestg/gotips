package slog_with_request_id

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"
)

func TestSlogWithRequestIDHandler_Handle(t *testing.T) {
	t.Run("without request id", func(t *testing.T) {
		var (
			history              bytes.Buffer
			handler              = slog.NewTextHandler(&history, &slog.HandlerOptions{})
			handlerWithRequestID = NewSlogWithRequestIDHandler(handler)
			logger               = slog.New(handlerWithRequestID)
		)

		logger.InfoContext(context.Background(), "Hello, world!")

		line := history.String()
		t.Log(line)
		if strings.Contains(line, "request_id") {
			t.Error("request_id should not be present")
		}
	})

	t.Run("with request id", func(t *testing.T) {
		var (
			history              bytes.Buffer
			handler              = slog.NewTextHandler(&history, &slog.HandlerOptions{})
			handlerWithRequestID = NewSlogWithRequestIDHandler(handler)
			logger               = slog.New(handlerWithRequestID)
		)

		// adding request id to log record through context.
		ctx := WithRequestID(context.Background(), "xyz")
		logger.InfoContext(ctx, "Hello, world!")

		line := history.String()
		t.Log(line)
		if !strings.Contains(line, "request_id=xyz") {
			t.Error("request_id should be present")
		}
	})
}
