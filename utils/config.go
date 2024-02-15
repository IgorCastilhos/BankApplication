package utils

import (
    "github.com/spf13/viper"
    "time"
)

// Config armazenará todas as configurações da aplicação.
// Os valores serão lidos pelo viper a partir de um arquivo ou das variáveis de ambiente.
type Config struct {
    DBSource             string        `mapstructure:"DB_SOURCE"`
    MigrationURL         string        `mapstructure:"MIGRATION_URL"`
    HTTPServerAddress    string        `mapstructure:"HTTP_SERVER_ADDRESS"`
    GRPCServerAddress    string        `mapstructure:"GRPC_SERVER_ADDRESS"`
    TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
    AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
    RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
    EmailSenderName      string        `mapstructure:"EMAIL_SENDER_NAME"`
    EmailSenderAddress   string        `mapstructure:"EMAIL_SENDER_ADDRESS"`
    EmailSenderPassword  string        `mapstructure:"EMAIL_SENDER_PASSWORD"`
}

// LoadConfig lê as configurações de um arquivo de configurações, a partir do path, se ele existir.
// Ou sobrescreve os valores com as variáveis de ambiente, caso elas sejam passadas.
func LoadConfig(path string) (config Config, err error) {
    viper.AddConfigPath(path)
    viper.SetConfigName("app")
    viper.SetConfigType("env")
    
    viper.AutomaticEnv()
    
    err = viper.ReadInConfig()
    if err != nil {
        return
    }
    
    err = viper.Unmarshal(&config)
    return
}
