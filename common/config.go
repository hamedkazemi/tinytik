package common

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
	"strconv"
	"sync"
)

var Config config

var onceConfig sync.Once

func init() {
	onceConfig.Do(func() {
		if _, err := toml.DecodeFile("./config.toml", &Config); err != nil {
			fmt.Println(err)
		}

		// Replace OS env with toml config file if available

		// Database
		if len(os.Getenv("ENV_DB_SERVER")) > 0 {
			Config.Database.Server = os.Getenv("ENV_DB_SERVER")
		}
		if len(os.Getenv("ENV_DB_DATABASE")) > 0 {
			Config.Database.Database = os.Getenv("ENV_DB_DATABASE")
		}
		if len(os.Getenv("ENV_DB_USERNAME")) > 0 {
			Config.Database.User = os.Getenv("ENV_DB_USERNAME")
		}
		if len(os.Getenv("ENV_DB_PASSWORD")) > 0 {
			Config.Database.Password = os.Getenv("ENV_DB_PASSWORD")
		}
		if len(os.Getenv("ENV_DB_PORT")) > 0 {
			Config.Database.Port = os.Getenv("ENV_DB_PORT")
		}
		if len(os.Getenv("ENV_DB_DEBUG")) > 0 {
			Config.Database.Debug, _ = strconv.ParseBool(os.Getenv("ENV_DB_DEBUG"))
		}
		// Kafka
		if len(os.Getenv("ENV_KAFKA_IP")) > 0 {
			Config.Kafka.Ip = os.Getenv("ENV_KAFKA_IP")
		}
		if len(os.Getenv("ENV_KAFKA_PORT")) > 0 {
			Config.Kafka.Port = os.Getenv("ENV_KAFKA_PORT")
		}
		if len(os.Getenv("ENV_KAFKA_TOPIC")) > 0 {
			Config.Kafka.Topic = os.Getenv("ENV_KAFKA_TOPIC")
		}

		if len(os.Getenv("ENV_REDIS_CONNECTION_STRING")) > 0 {
			Config.Redis.ConnectionString = os.Getenv("ENV_REDIS_CONNECTION_STRING")
		}

		// APP
		if len(os.Getenv("ENV_APP_HOST")) > 0 {
			Config.App.Host = os.Getenv("ENV_APP_HOST")
		}
		if len(os.Getenv("ENV_APP_PROXY")) > 0 {
			Config.App.Proxy = os.Getenv("ENV_APP_PROXY")
		}
		if len(os.Getenv("ENV_APP_ENVIRONMENT")) > 0 {
			Config.App.Environment = os.Getenv("ENV_APP_ENVIRONMENT")
		}
	})
}
