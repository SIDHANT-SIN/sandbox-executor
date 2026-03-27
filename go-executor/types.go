package main

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