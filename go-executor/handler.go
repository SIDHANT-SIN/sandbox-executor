
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

var sem = make(chan struct{}, 3)

type Req struct {
	Lang      string `json:"lang"`
	Code      string `json:"code"`
	ProblemID string `json:"problem_id"`
}

type Resp struct {
	Output  string `json:"output"`
	Error   string `json:"error"`
	Status  string `json:"status"`
	Correct bool   `json:"correct"`
}

type Lang struct {
	Img        string
	File       string
	CompileCmd []string
	RunCmd     []string
}

var langs = map[string]Lang{
	"cpp": {
		Img:  "gcc",
		File: "code.cpp",
		CompileCmd: []string{
			"bash", "-c", "g++ /code/code.cpp -o /code/out && chmod +x /code/out",
		},
		RunCmd: []string{
			"bash", "-c", "timeout 2s /code/out",
		},
	},
	"python": {
		Img:        "python",
		File:       "code.py",
		CompileCmd: nil,
		RunCmd: []string{
			"bash", "-c", "chmod +r /code/code.py && timeout 2s python /code/code.py",
		},
	},
	"java": {
		Img:  "eclipse-temurin",
		File: "Main.java",
		CompileCmd: []string{
			"bash", "-c", "javac /code/Main.java && chmod +x /code/*.class",
		},
		RunCmd: []string{
			"bash", "-c",
			"timeout 2s java -Xms32m -Xmx64m -cp /code Main",
		},
	},
}

func writeCode(dir, filename, code string) error {
	path := dir + "/" + filename
	return os.WriteFile(path, []byte(code), 0644)
}

func execHandler(c *gin.Context) {
	var req Req

	if err := c.BindJSON(&req); err != nil {
		log.Println("ERROR: bad request body:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	log.Printf("INCOMING REQUEST: lang=%s problem_id=%s\n", req.Lang, req.ProblemID)

	if req.ProblemID == "" {
		log.Println("ERROR: missing problem_id in request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "problem_id is required"})
		return
	}

	sem <- struct{}{}
	res := execute(req)
	<-sem

	log.Printf("RESPONSE: status=%s correct=%v error=%s\n", res.Status, res.Correct, res.Error)

	c.JSON(http.StatusOK, res)
}
