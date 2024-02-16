package grpcApi

import (
    "context"
    "github.com/rs/zerolog/log"
    "google.golang.org/grpc"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
    "net/http"
    "time"
)

func GrpcLogger(
  ctx context.Context,
  req any,
  info *grpc.UnaryServerInfo,
  handler grpc.UnaryHandler,
) (resp any, err error) {
    startTime := time.Now()
    result, err := handler(ctx, req)
    duration := time.Since(startTime)
    
    statusCode := codes.Unknown
    if st, ok := status.FromError(err); ok {
        statusCode = st.Code()
    }
    
    logger := log.Info()
    if err != nil {
        logger = log.Error().Err(err)
    }
    
    logger.Str("protocol", "grpc").
        Str("method", info.FullMethod).
        Int("status_code", int(statusCode)).
        Str("status_text", statusCode.String()).
        Dur("duration", duration).
        Msg("recebida uma requisição gRPC")
    return result, err
}

type ResponseRecorder struct {
    http.ResponseWriter
    StatusCode int
    Body       []byte
}

func (rec *ResponseRecorder) WriteHeader(statusCode int) {
    rec.StatusCode = statusCode
    rec.ResponseWriter.WriteHeader(statusCode)
}

func (rec *ResponseRecorder) Write(body []byte) (int, error) {
    rec.Body = body
    return rec.ResponseWriter.Write(body)
}

func HTTPLogger(handler http.Handler) http.Handler {
    return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
        startTime := time.Now()
        rec := &ResponseRecorder{
            ResponseWriter: writer,
            StatusCode:     http.StatusOK,
        }
        handler.ServeHTTP(rec, request)
        duration := time.Since(startTime)
        
        logger := log.Info()
        if rec.StatusCode != http.StatusOK {
            logger = log.Error().Bytes("body", rec.Body)
        }
        
        logger.Str("protocol", "http").
            Str("method", request.Method).
            Str("path", request.RequestURI).
            Int("status_code", rec.StatusCode).
            Str("status_text", http.StatusText(rec.StatusCode)).
            Dur("duration", duration).
            Msg("recebida uma requisição HTTP")
    })
}
