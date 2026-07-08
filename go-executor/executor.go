package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func execute(req Req) Resp {
	fmt.Println("===== NEW EXECUTION REQUEST =====")
	fmt.Println("Lang:", req.Lang)
	fmt.Println("ProblemID:", req.ProblemID)

	lang, ok := langs[req.Lang]
	if !ok {
		fmt.Println("ERROR: unsupported language:", req.Lang)
		return Resp{Error: "unsupported language", Status: "error"}
	}

	testInput, err := readBlob(req.ProblemID, os.Getenv("TEST_FILE"))
	if err != nil {
		fmt.Println("ERROR: failed to read test file for problem", req.ProblemID, "-", err)
		return Resp{Error: "failed to load test data", Status: "error"}
	}

	expectedOutput, err := readBlob(req.ProblemID, os.Getenv("SOL_FILE"))
	if err != nil {
		fmt.Println("ERROR: failed to read solution file for problem", req.ProblemID, "-", err)
		return Resp{Error: "failed to load solution data", Status: "error"}
	}

	fmt.Println("Test input loaded, length:", len(testInput))
	fmt.Println("Expected output loaded, length:", len(expectedOutput))

	baseDir := os.Getenv("CODE_SHARED_DIR")
	if baseDir == "" {
		baseDir = os.TempDir() 
	}

	dir, err := os.MkdirTemp(baseDir, "exec-*")
	if err != nil {
		fmt.Println("ERROR: failed to create temp dir:", err)
		return Resp{Error: "infrastructure error", Status: "error"}
	}
	defer os.RemoveAll(dir)

	err = writeCode(dir, lang.File, req.Code)
	if err != nil {
		fmt.Println("ERROR: file write error:", err)
		return Resp{Error: "file write error", Status: "error"}
	}

	files, _ := os.ReadDir(dir)
	fmt.Println("Files written in temp dir:", dir)
	for _, f := range files {
		fmt.Println("-", f.Name())
	}

	if lang.CompileCmd != nil {
		out, err, status := runDocker(dir, lang.Img, lang.CompileCmd, 10*time.Second, "")
		if status == "infra_timeout" {
			fmt.Println("ERROR: compilation timed out")
			return Resp{
				Error:  "Compilation Timeout",
				Status: "compile_timeout",
			}
		}
		if err != nil {
			fmt.Println("ERROR: compile error:", out)
			return Resp{
				Error:  out,
				Status: "compile_error",
			}
		}
	}

	out, err, status := runDocker(dir, lang.Img, lang.RunCmd, 0, testInput)

	if status == "timeout" {
		fmt.Println("ERROR: time limit exceeded")
		return Resp{
			Error:  "Time Limit Exceeded",
			Status: "timeout",
		}
	}
	if status == "memory_limit_exceeded" {
		fmt.Println("ERROR: memory limit exceeded")
		return Resp{
			Error:  "Memory Limit Exceeded",
			Status: "exceeded memory",
		}
	}
	if status == "segmentation_fault" {
		fmt.Println("ERROR: segmentation fault")
		return Resp{
			Error:  "Segmentation Fault",
			Status: "segmentation_fault",
		}
	}
	if status == "infra_timeout" {
		fmt.Println("ERROR: server/infra timeout")
		return Resp{
			Error:  "Server Timeout",
			Status: "server_error",
		}
	}
	if err != nil {
		fmt.Println("ERROR: runtime error:", err.Error())
		fmt.Println("RAW OUTPUT ON ERROR:", out)
		return Resp{
			Error:  err.Error(),
			Status: "runtime_error",
		}
	}

	fmt.Println("===== COMPARING OUTPUT =====")
	fmt.Println("ACTUAL OUTPUT:")
	fmt.Println(out)
	fmt.Println("EXPECTED OUTPUT:")
	fmt.Println(expectedOutput)

	correct := compareOutput(out, expectedOutput)
	fmt.Println("MATCH RESULT:", correct)

	return Resp{
		Output:  out,
		Status:  "success",
		Correct: correct,
	}
}

func compareOutput(actual, expected string) bool {
	actualLines := splitLines(actual)
	expectedLines := splitLines(expected)

	if len(actualLines) != len(expectedLines) {
		fmt.Println("LINE COUNT MISMATCH: got", len(actualLines), "lines, want", len(expectedLines), "lines")
		return false
	}

	for i := range actualLines {
		a := strings.TrimSpace(actualLines[i])
		e := strings.TrimSpace(expectedLines[i])
		if a != e {
			fmt.Printf("LINE %d MISMATCH: got %q, want %q\n", i, a, e)
			return false
		}
	}

	return true
}

func splitLines(s string) []string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	lines := strings.Split(s, "\n")
	for len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) == "" {
		lines = lines[:len(lines)-1]
	}
	return lines
}

func runDocker(dir, img string, cmd []string, t time.Duration, stdin string) (string, error, string) {
	fmt.Println("START runDocker : - ")

	//  DinD check
	fmt.Println("DOCKER AVAILABILITY CHECK : ")
	cmdCheck := exec.Command("docker", "ps")
	outCheck, errCheck := cmdCheck.CombinedOutput()
	fmt.Println("DOCKER PS OUTPUT:\n", string(outCheck))
	fmt.Println("DOCKER PS ERROR:", errCheck)

	args := []string{
		"run", "--rm", "-i",
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
	fmt.Println("STDIN LENGTH:", len(stdin))
	fmt.Println("FULL ARGS:", args)

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

	fmt.Println("=== EXECUTING DOCKER COMMAND ===")
	c := exec.CommandContext(ctx, "docker", args...)
	c.Stdin = strings.NewReader(stdin)

	out, err := c.CombinedOutput()

	fmt.Println("=== DOCKER EXECUTION RESULT ===")
	fmt.Println("OUTPUT:\n", string(out))
	fmt.Println("ERROR:", err)

	if ctx.Err() == context.DeadlineExceeded {
		fmt.Println("TIMEOUT HIT")
		return "", ctx.Err(), "infra_timeout"
	}

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

		fmt.Println("RUNTIME ERROR DETECTED")
		return string(out), err, "runtime_error"
	}

	fmt.Println("SUCCESS")
	return string(out), nil, "success"
}