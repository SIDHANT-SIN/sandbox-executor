package main

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "os"
)

var sem = make(chan struct{}, 3) 

type Req struct {
    Lang  string `json:"lang"`
    Code  string `json:"code"`
    Input string `json:"input"`
}

type Resp struct {
    Output string `json:"output"`
    Error  string `json:"error"`
    Status string `json:"status"`
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
        Img:  "python",
        File: "code.py",
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
        c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
        return
    }

    sem <- struct{}{}

    res := execute(req)

    <-sem

    c.JSON(http.StatusOK, res)
}

