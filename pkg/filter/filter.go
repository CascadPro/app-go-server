package filter

import (
	"maps"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ResponseBody struct {
	Status    int     `json:"status"`
	Message   string  `json:"message"`
	Timestamp string  `json:"timestamp"`
	Cause     *string `json:"cause"`
}

func Success(c *gin.Context, message string, additional ...gin.H) {
	body := gin.H{
		"status":    http.StatusOK,
		"message":   message,
		"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
	}

	for _, add := range additional {
		maps.Copy(body, add)
	}

	c.JSON(http.StatusOK, body)
}

type ErrorParams struct {
	Status  int
	Message string
	Cause   *string
}

func Error(c *gin.Context, p ErrorParams) {
	message := p.Message

	if p.Status == http.StatusInternalServerError {
		message = "Что-то пошло не так! Попробуйте позже"
	}

	body := ResponseBody{
		Status:    p.Status,
		Message:   message,
		Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
		Cause:     p.Cause,
	}

	c.JSON(p.Status, body)
}
