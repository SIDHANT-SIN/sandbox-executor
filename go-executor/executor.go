package main

import (
	"os"
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
			Error:  out,
			Status: "runtime_error",
		}
	}

	return Resp{
		Output: out,
		Status: "success",
	}
}