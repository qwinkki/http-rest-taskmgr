package core_http_response

import (
	"encoding/json"
	"fmt"
	"net/http"

	core_logger "github.com/qwinkki/http-rest-taskmgr/internal/core/logger"
	"go.uber.org/zap"
)

type HTTPRestponseHandler struct {
	log core_logger.Logger
	rw  http.ResponseWriter
}

func NewHTTPResponseHandler(log core_logger.Logger, rw http.ResponseWriter) *HTTPRestponseHandler {
	return &HTTPRestponseHandler{
		log: log,
		rw:  rw,
	}
}

func (h *HTTPRestponseHandler) Handle(p any, msg string) {
	statusCode := http.StatusInternalServerError
	err := fmt.Errorf("unexpected panic: %v", p)

	h.log.Error(msg, zap.Error(err))
	h.rw.WriteHeader(statusCode)

	response := map[string]string{
		"message": msg,
		"error":   err.Error(),
	}
	if err := json.NewEncoder(h.rw).Encode(response); err != nil {
		h.log.Error("failed to write response", zap.Error(err))
	}
}
