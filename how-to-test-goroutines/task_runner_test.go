package how_to_test_goroutines

import (
	"bytes"
	"context"
	"log/slog"
	"math/rand"
	"testing"
	"time"

	"github.com/josestg/gotips/how-to-test-goroutines/await"
)

func NewTask(l *slog.Logger, name string) Task {
	return func(ctx context.Context, args []string) {
		l.InfoContext(ctx, "Task started", "name", name)
		delay := time.Duration(rand.Intn(5)*100) * time.Millisecond
		select {
		case <-ctx.Done():
			l.InfoContext(ctx, "Task canceled", "name", name)
		case <-time.After(delay): // simulate some work.
			l.InfoContext(ctx, "Task finished", "name", name, "args", args)
		}
	}
}
func TestTaskRunner_Run(t *testing.T) {
	var logHistory bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&logHistory, &slog.HandlerOptions{}))
	defer func() { t.Log(logHistory.String()) }()

	task1 := NewTask(logger, "task1")
	task2 := NewTask(logger, "task2")
	task3 := NewTask(logger, "task3")

	runner := NewTaskRunner(logger, task1, task2, task3)

	ctx := context.Background()
	ctx = await.Context(ctx)
	args := []string{"a", "b", "c"}
	runner.Run(ctx, args)
}
