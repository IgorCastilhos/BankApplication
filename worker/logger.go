package worker

import (
    "fmt"
    "github.com/rs/zerolog"
    "github.com/rs/zerolog/log"
)

type Logger struct{}

func NewLogger() *Logger {
    return &Logger{}
}

// Print retorna o nível do log e uma string concatenando os argumentos subjacentes
func (logger *Logger) Print(level zerolog.Level, args ...interface{}) {
    log.WithLevel(level).Msg(fmt.Sprint(args...))
}

// Debug logs a message at Debug level.
func (logger *Logger) Debug(args ...interface{}) { logger.Print(zerolog.DebugLevel, args...) }

// Info logs a message at Info level.
func (logger *Logger) Info(args ...interface{}) { logger.Print(zerolog.InfoLevel, args...) }

// Warn logs a message at Warning level.
func (logger *Logger) Warn(args ...interface{}) { logger.Print(zerolog.WarnLevel, args...) }

// Error logs a message at Error level.
func (logger *Logger) Error(args ...interface{}) { logger.Print(zerolog.ErrorLevel, args...) }

// Fatal loga uma mensagem nível Fatal e o processo encerrará com status definido para 1
func (logger *Logger) Fatal(args ...interface{}) { logger.Print(zerolog.FatalLevel, args...) }
