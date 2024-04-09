package how_to_test_goroutines

import (
	"context"
	"log/slog"
)

type Task func(ctx context.Context, args []string)

type TaskRunner struct {
	log   *slog.Logger
	tasks []Task
}

func NewTaskRunner(l *slog.Logger, tasks ...Task) *TaskRunner {
	return &TaskRunner{log: l, tasks: tasks}
}

func (r *TaskRunner) Run(ctx context.Context, args []string) {
	r.log.InfoContext(ctx, "Run tasks", "args", args)
	ctx = context.WithoutCancel(ctx)
	for _, task := range r.tasks {
		go task(ctx, args)
	}
}
