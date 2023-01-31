package launch_editor

import (
	"context"
	"os/exec"
)

type crossSpawn struct {
	*exec.Cmd
	cancel context.CancelFunc
}

func spawn(name string, args ...string) *crossSpawn {
	ctx, cancel := context.WithCancel(context.Background())
	var cp *exec.Cmd
	if isWindows {
		args = append([]string{"/C", name}, args...)
		cp = exec.CommandContext(ctx, "cmd.exe", args...)
	} else {
		cp = exec.CommandContext(ctx, name, args...)
	}
	return &crossSpawn{
		Cmd:    cp,
		cancel: cancel,
	}
}
