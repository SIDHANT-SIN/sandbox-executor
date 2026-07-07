package main

import (
	"os"

	"fmt"
		"context"
	
	"os/exec"
	"time"
)

func execute(req Req) Resp {

	lang, ok := langs[req.Lang]
	if !ok {
		return Resp{Error: "unsupported language", Status: "error"}
	}

	dir, _ := os.MkdirTemp("", "exec-*")
	defer os.RemoveAll(dir)

	err := writeCode(dir, lang.File, req.Code)
	if err != nil {
		return Resp{Error: "file write error", Status: "error"}
	}

	files, _ := os.ReadDir(dir)
    fmt.Println("Files written in temp dir:", dir)
    for _, f := range files {
         fmt.Println("-", f.Name())
}

	// COMPILE (with safety timeout)
	if lang.CompileCmd != nil {

		out, err, status := runDocker(dir, lang.Img, lang.CompileCmd, 10*time.Second)

		if status == "infra_timeout" {
			return Resp{
				Error:  "Compilation Timeout",
				Status: "compile_timeout",
			}
		}

		if err != nil {
			return Resp{
				Error:  out,
				Status: "compile_error",
			}
		}
	}

	// RUN (timeout handled INSIDE container)
	out, err, status := runDocker(dir, lang.Img, lang.RunCmd, 0)

	if status == "timeout" {
		return Resp{
			Error:  "Time Limit Exceeded",
			Status: "timeout",
		}
	}
	if status == "memory_limit_exceeded" {
		return Resp{
			Error:  "Memory Limit Exceeded",
			Status: "exceeded memory",
		}
	}
		if status == "segmentation_fault" {
		return Resp{
			Error:  "Segmentation Fault",
			Status: "segmentation_fault",
		}
	}

	if status == "infra_timeout" {
		return Resp{
			Error:  "Server Timeout",
			Status: "server_error",
		}
	}

	if err != nil {
		return Resp{
			Error:  err.Error(),
			Status: "runtime_error",
		}
	}

	return Resp{
		Output: out,
		Status: "success",
	}
}

func runDocker(dir, img string, cmd []string, t time.Duration) (string, error, string) {

	fmt.Println("START runDocker : - ")

	//  DinD check
	fmt.Println("DOCKER AVAILABILITY CHECK : ")
	cmdCheck := exec.Command("docker", "ps")
	outCheck, errCheck := cmdCheck.CombinedOutput()
	fmt.Println("DOCKER PS OUTPUT:\n", string(outCheck))
	fmt.Println("DOCKER PS ERROR:", errCheck)

	
	args := []string{
		"run", "--rm",
		"--memory=256m",
		"--cpus=0.5",
		"--network=none",
		"--pids-limit=64",
		"-v", dir + ":/code",
		"-w", "/code", 
		img,
	}

	args = append(args, cmd...)

	fmt.Println("=== DOCKER RUN CONFIG ===")
	fmt.Println("DIR:", dir)
	fmt.Println("IMAGE:", img)
	fmt.Println("CMD:", cmd)
	fmt.Println("FULL ARGS:", args)

	// Context setup
	var ctx context.Context
	var cancel context.CancelFunc

	if t > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), t)
		fmt.Println("TIMEOUT SET:", t)
	} else {
		ctx, cancel = context.WithCancel(context.Background())
		fmt.Println("NO TIMEOUT")
	}
	defer cancel()

	//Execute docker
	fmt.Println("=== EXECUTING DOCKER COMMAND ===")
	c := exec.CommandContext(ctx, "docker", args...)

	out, err := c.CombinedOutput()

	fmt.Println("=== DOCKER EXECUTION RESULT ===")
	fmt.Println("OUTPUT:\n", string(out))
	fmt.Println("ERROR:", err)

	// Timeout check
	if ctx.Err() == context.DeadlineExceeded {
		fmt.Println("TIMEOUT HIT")
		return "", ctx.Err(), "infra_timeout"
	}

	// Error classification
	if err != nil {

		fmt.Println("=== ERROR ANALYSIS ===")
		fmt.Println("RAW ERROR:", err)

		if exitErr, ok := err.(*exec.ExitError); ok {
			code := exitErr.ExitCode()
			fmt.Println("EXIT CODE:", code)

			if code == 124 {
				fmt.Println("TIME LIMIT EXCEEDED")
				return "", err, "timeout"
			}
			if code == 137 {
				fmt.Println("MEMORY LIMIT EXCEEDED")
				return "", err, "memory_limit_exceeded"
			}
			if code == 139 {
				fmt.Println("SEGMENTATION FAULT")
				return "", err, "segmentation_fault"
			}
		}

		// show actual output
		fmt.Println("RUNTIME ERROR DETECTED")
		return string(out), err, "runtime_error"
	}

	fmt.Println("SUCCESS")
	return string(out), nil, "success"
}