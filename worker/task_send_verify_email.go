package worker

import (
    "context"
    "encoding/json"
    "fmt"
    "github.com/hibiken/asynq"
    "github.com/rs/zerolog/log"
)

const TaskSendVerifyEmail = "task:send_verify_email"

type PayloadSendVerifyEmail struct {
    Username string `json:"username"`
}

func (distributor *RedisTaskDistributor) DistributeTaskSendVerifyEmail(
  ctx context.Context,
  payload *PayloadSendVerifyEmail,
  opts ...asynq.Option) error {
    jsonPayload, err := json.Marshal(payload)
    if err != nil {
        return fmt.Errorf("falha ao fazer o marshal do payload da tarefa: %w", err)
    }
    
    task := asynq.NewTask(TaskSendVerifyEmail, jsonPayload, opts...)
    info, err := distributor.client.EnqueueContext(ctx, task)
    if err != nil {
        return fmt.Errorf("falha ao enfileirar tarefa: %w", err)
    }
    
    log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).Str("queue", info.Queue).Int("max_retry", info.MaxRetry).Msg("tarefa enfileirada")
    return nil
}

func (processor *RedisTaskProcessor) ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
    var payload PayloadSendVerifyEmail
    if err := json.Unmarshal(task.Payload(), &payload); err != nil {
        return fmt.Errorf("falha ao fazer unmarshal do payload: %w", asynq.SkipRetry)
    }
    
    user, err := processor.store.GetUser(ctx, payload.Username)
    if err != nil {
        return fmt.Errorf("falha ao buscar usuário: %w", err)
    }
    
    log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).Str("email", user.Email).Msg("tarefa processada")
    return nil
}
