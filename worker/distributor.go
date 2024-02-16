package worker

import (
    "context"
    "github.com/hibiken/asynq"
)

type TaskDistributor interface {
    DistributeTaskSendVerifyEmail(
      ctx context.Context,
      payload *PayloadSendVerifyEmail,
      opts ...asynq.Option) error
}

type RedisTaskDistributor struct {
    client *asynq.Client
}

// NewRedisTaskDistributor retorna uma interface, pois é preciso forçar a struct RedisTaskDistributor a implementar a interface TaskDistributor.
// Se não forem implementadas todas as funções requisitadas pela interface, o compilador irá alertar um erro.
func NewRedisTaskDistributor(redisOpt asynq.RedisClientOpt) TaskDistributor {
    client := asynq.NewClient(redisOpt)
    return &RedisTaskDistributor{
        client: client,
    }
}
