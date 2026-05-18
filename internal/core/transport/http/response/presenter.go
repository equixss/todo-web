package core_http_response_presenter

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	core_errors "github.com/equixss/todo-web/internal/core/errors"
	core_logger "github.com/equixss/todo-web/internal/core/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type HTTPResponsePresenter struct {
}

func NewHTTPResponsePresenter() *HTTPResponsePresenter {
	return &HTTPResponsePresenter{}
}

func (p *HTTPResponsePresenter) ResponseNoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func (p *HTTPResponsePresenter) JSONResponse(
	c *gin.Context,
	responseBody any,
	statusCode int,
) {
	if statusCode == http.StatusNoContent {
		p.ResponseNoContent(c)
		return
	}
	//ловим ошибки jsone encode
	jsonBytes, err := json.Marshal(responseBody)
	if err != nil {
		log := core_logger.FromContext(c.Request.Context())
		log.Error("failed to encode response body to json", zap.Error(err))
		p.ErrorResponse(c, err, "Error generating server response")
		return
	}
	c.Header("Content-Type", "application/json; charset=utf-8")
	c.Status(statusCode)
	if _, writeErr := c.Writer.Write(jsonBytes); writeErr != nil {
		log := core_logger.FromContext(c.Request.Context())
		log.Warn("failed to write json bytes to connection", zap.Error(writeErr))
	}
}

func (p *HTTPResponsePresenter) ErrorResponse(c *gin.Context, err error, msg string) {
	var (
		statusCode int
		logFunc    func(string, ...zap.Field)
	)
	log := core_logger.FromContext(c.Request.Context())

	switch {
	case errors.Is(err, core_errors.ErrInvalidArgument):
		statusCode = http.StatusBadRequest
		logFunc = log.Warn
	case errors.Is(err, core_errors.ErrNotFound):
		statusCode = http.StatusNotFound
		logFunc = log.Debug
	case errors.Is(err, core_errors.ErrConflict):
		statusCode = http.StatusConflict
		logFunc = log.Warn
	case errors.Is(err, core_errors.ErrUnauthorized):
		statusCode = http.StatusUnauthorized
		logFunc = log.Warn
	case errors.Is(err, core_errors.ErrInvalidPassword):
		statusCode = http.StatusUnauthorized
		logFunc = log.Warn
	default:
		statusCode = http.StatusInternalServerError
		logFunc = log.Error
	}

	logFunc(msg, zap.Error(err))
	p.errorResponse(c, statusCode, err, msg)
}

func (h *HTTPResponsePresenter) PanicResponse(c *gin.Context, p any, msg string) {
	err := fmt.Errorf("unexpected panic:%v", p)
	logger := core_logger.FromContext(c.Request.Context())
	logger.Error(msg, zap.Any("panic_info", p), zap.Stack("stack"))
	h.errorResponse(c, http.StatusInternalServerError, err, msg)
}

func (p *HTTPResponsePresenter) errorResponse(
	c *gin.Context,
	statusCode int,
	err error,
	msg string,
) {
	c.AbortWithStatusJSON(statusCode, gin.H{
		"message": msg,
		"error":   err.Error(),
	})
}
