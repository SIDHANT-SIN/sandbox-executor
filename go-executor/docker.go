package main

import (
	"context"
	"os/exec"
	"time"
)

func runDocker(dir, img string, cmd []string, t time.Duration) (string, error, string) {

	args := []string{
		"run", "--rm",
		"--memory=256m",
		"--cpus=0.5",
		"--network=none",
		"--pids-limit=64",
		"-v", dir + ":/code",
		img,
	}

	args = append(args, cmd...)

	var ctx context.Context
	var cancel context.CancelFunc

	if t > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), t)
	} else {
		ctx, cancel = context.WithCancel(context.Background())
	}
	defer cancel()

	c := exec.CommandContext(ctx, "docker", args...)

	out, err := c.CombinedOutput()

	// compile error
	if ctx.Err() == context.DeadlineExceeded {
		return "", ctx.Err(), "infra_timeout"
	}

	if err != nil {

		// TLE , MLE, segfault
		if exitErr, ok := err.(*exec.ExitError); ok {
			if exitErr.ExitCode() == 124 {
				return "", err, "timeout"
			}
			    if exitErr.ExitCode() == 137 {
        return "", err, "memory_limit_exceeded"
    }

    if exitErr.ExitCode() == 139 {
        return "", err, "segmentation_fault"
    }
		}
          // runtime error
		return string(out), err, "runtime_error"
	}

	return string(out), nil, "success"
}