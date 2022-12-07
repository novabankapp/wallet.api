package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Max-Age", "600")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
		} else {
			fmt.Println("In CORSMiddleware")
			fmt.Println(c)
			c.Next()
		}
	}
}

func ErrorHandler(c *gin.Context) {
	c.Next()

	for _, err := range c.Errors {
		err.Error()
	}

	c.JSON(http.StatusInternalServerError, "")
}

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestId := c.GetString("x-request-id")
		clientIp := c.ClientIP()
		userAgent := c.Request.UserAgent()
		method := c.Request.Method
		path := c.Request.URL.Path

		t := time.Now()
		fmt.Println("In Middleware")
		fmt.Println(c)
		c.Next()

		latency := float32(time.Since(t).Seconds())

		status := c.Writer.Status()
		log.Info().Str("request_id", requestId).Str("client_ip", clientIp).
			Str("user_agent", userAgent).Str("method", method).Str("path", path).
			Float32("latency", latency).Int("status", status).Msg("")

	}

}

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate UUID
		id := uuid.New().String()
		// Set context variable
		c.Set("x-request-id", id)
		// Set header
		c.Header("x-request-id", id)
		c.Next()
	}
}
