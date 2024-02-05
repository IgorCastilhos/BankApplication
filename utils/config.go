package utils

import (
	"github.com/spf13/viper"
)

// Config armazenará todas as configurações da aplicação.
// Os valores serão lidos pelo viper a partir de um arquivo ou das variáveis de ambiente.
type Config struct {
	DBDriver          string `mapstructure:"DB_DRIVER"`
	DBSource          string `mapstructure:"DB_SOURCE"`
	HTTPServerAddress string `mapstructure:"HTTP_SERVER_ADDRESS"`
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