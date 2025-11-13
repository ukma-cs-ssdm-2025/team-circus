package httpx

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ResponseSpec describes how to translate a matched error into a HTTP response.
type ResponseSpec struct {
	Target     error
	Status     int
	Message    string
	LogMessage string
	LogLevel   zapcore.Level
}

// RespondError sends a JSON error payload with a consistent schema.
func RespondError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"error": message})
}

// logWithLevel routes logs through the provided logger using the desired level.
func logWithLevel(logger *zap.Logger, level zapcore.Level, message string, fields ...zap.Field) {
	switch level {
	case zapcore.DebugLevel:
		logger.Debug(message, fields...)
	case zapcore.InfoLevel:
		logger.Info(message, fields...)
	case zapcore.WarnLevel:
		logger.Warn(message, fields...)
	default:
		logger.Error(message, fields...)
	}
}

// HandleError applies the provided response specs to the given error. The first matching
// response spec wins. If none match, the fallback spec is used. Returns true when a response
// was written to the client.
func HandleError(
	c *gin.Context,
	logger *zap.Logger,
	err error,
	fallback ResponseSpec,
	fields []zap.Field,
	mappings ...ResponseSpec,
) bool {
	if err == nil {
		return false
	}

	for _, mapping := range mappings {
		if mapping.Target != nil && errors.Is(err, mapping.Target) {
			if mapping.LogMessage != "" {
				logWithLevel(logger, mapping.LogLevel, mapping.LogMessage, append(fields, zap.Error(err))...)
			}
			RespondError(c, mapping.Status, mapping.Message)
			return true
		}
	}

	if fallback.LogMessage != "" {
		logWithLevel(logger, fallback.LogLevel, fallback.LogMessage, append(fields, zap.Error(err))...)
	}
	RespondError(c, fallback.Status, fallback.Message)
	return true
}

// RequestContextFields returns common logging fields derived from the request context.
func RequestContextFields(c *gin.Context, extra ...zap.Field) []zap.Field {
	fields := []zap.Field{
		zap.String("request_id", c.GetString("request_id")),
		zap.String("client_ip", c.ClientIP()),
	}
	return append(fields, extra...)
}

// ResolveUserUUID extracts the authenticated user UUID from context, writing an error response
// if the value is missing or has an unexpected type.
func ResolveUserUUID(c *gin.Context) (uuid.UUID, bool) {
	userUUIDValue, exists := c.Get("user_uid")
	if !exists {
		RespondError(c, http.StatusUnauthorized, "user context missing")
		return uuid.Nil, false
	}

	userUUID, ok := userUUIDValue.(uuid.UUID)
	if !ok {
		RespondError(c, http.StatusUnauthorized, "invalid user context")
		return uuid.Nil, false
	}

	return userUUID, true
}

// ParseUUIDParam validates that the provided URI parameter contains a valid UUID.
func ParseUUIDParam(
	c *gin.Context,
	logger *zap.Logger,
	paramName string,
	logMessage string,
	fields ...zap.Field,
) (uuid.UUID, bool) {
	value := c.Param(paramName)
	parsed, err := uuid.Parse(value)
	if err != nil {
		logFields := append([]zap.Field{
			zap.String("param", paramName),
			zap.String("value", value),
		}, fields...)
		logger.Error(logMessage, append(logFields, zap.Error(err))...)
		RespondError(c, http.StatusBadRequest, "invalid uuid format")
		return uuid.Nil, false
	}

	return parsed, true
}
