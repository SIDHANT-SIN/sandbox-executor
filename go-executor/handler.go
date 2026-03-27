package main

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

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