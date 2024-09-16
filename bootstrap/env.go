package bootstrap

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Env struct {
	ServerAddress    string `mapstructure:"SERVER_ADDRESS"`
	ServerPort       string `mapstructure:"SERVER_PORT"`
	PostgresConn     string `mapstructure:"POSTGRES_CONN"`
	PostgresJDBCUrl  string `mapstructure:"POSTGRES_JDBC_URL"`
	PostgresUsername string `mapstructure:"POSTGRES_USERNAME"`
	PostgresPasswrod string `mapstructure:"POSTGRES_PASSWORD"`
	PostgresHost     string `mapstructure:"POSTGRES_HOST"`
	PostgresPort     string `mapstructure:"POSTGRES_PORT"`
	PostgresDatabase string `mapstructure:"POSTGRES_DATABASE"`
}

func NewEnv() *Env {
	env := Env{}
	viper.SetConfigFile(".env")

	if err := viper.ReadInConfig(); err != nil {
		logrus.Fatalf(err.Error())
	}

	if err := viper.Unmarshal(&env); err != nil {
		logrus.Fatalf(err.Error())
	}
	return &env
}
