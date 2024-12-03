package handler

import (
	"github.com/gin-gonic/gin"
)

func sendError(ctx *gin.Context, code int, msg string) {
	ctx.Header("Content-type", "application/json")
	ctx.JSON(code, gin.H{
		"message":   msg,
		"errorCode": code,
	})
}

func sendSuccess(ctx *gin.Context, code int, data interface{}) {
	ctx.Header("Content-type", "application/json")
	ctx.JSON(code, gin.H{
		"code": code,
		"data": data,
	})
}

type SuccessResponse struct {
	Code string      `json:"code" example:"201"`
	Data interface{} `json:"data"`
}

type ErrorResponse struct {
	ErrorCode string `json:"errorCode" example:"400"`
	Message   string `json:"message" example:"error on create"`
}
