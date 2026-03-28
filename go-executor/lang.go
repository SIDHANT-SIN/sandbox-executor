package main

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