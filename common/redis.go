package common

import "github.com/rueian/rueidis"

func ConnectRedis() rueidis.Client {
	c, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress: []string{Config.Redis.ConnectionString}, DisableCache: true,
	})
	if err != nil {
		panic(err)
	}
	return c
}
