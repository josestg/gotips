package slog_with_request_id

import (
	"context"
	"log/slog"
)

type requestIDKey struct{}

var _requestIDKey = new(requestIDKey)

// WithRequestID creates a new context with request ID injected.
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, _requestIDKey, requestID)
}

// SlogWithRequestIDHandler is a custom slog.Handler that adds request ID to log records, if present in the context.
type SlogWithRequestIDHandler struct {
	// embed the slog.Handler interface to inherit the slog.Handler methods.
	// by using embedding, we can override the specific methods we want to change, in this case, the Handle method.
	slog.Handler
}

// ensure SlogWithRequestIDHandler implements slog.Handler.
var _ slog.Handler = (*SlogWithRequestIDHandler)(nil)

// NewSlogWithRequestIDHandler creates a new SlogWithRequestIDHandler.
func NewSlogWithRequestIDHandler(h slog.Handler) *SlogWithRequestIDHandler {
	return &SlogWithRequestIDHandler{
		Handler: h,
	}
}

func (h *SlogWithRequestIDHandler) Handle(ctx context.Context, rec slog.Record) error {
	if requestID, ok := ctx.Value(_requestIDKey).(string); ok {
		rec.AddAttrs(slog.String("request_id", requestID))
	}
	return h.Handler.Handle(ctx, rec)
}
