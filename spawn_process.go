package launch_editor

import (
	"context"
	"os/exec"
)

type spawnProcess struct {
	*exec.Cmd
	cancel context.CancelFunc
}

func spawn(name string, args ...string) *spawnProcess {
	ctx, cancel := context.WithCancel(context.Background())
	return &spawnProcess{
		Cmd:    exec.CommandContext(ctx, name, args...),
		cancel: cancel,
	}
}
