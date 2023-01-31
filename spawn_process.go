package launch_editor

import (
	"context"
	"os/exec"

	"github.com/google/shlex"
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

func execCmd(cmd string) (output string, err error) {
	shellArgs, err := shlex.Split(cmd)
	if err != nil {
		return
	}

	cp := spawn(shellArgs[0], shellArgs[1:]...)
	buf, err := cp.CombinedOutput()
	if err != nil {
		return
	}
	output = string(buf)
	return
}
