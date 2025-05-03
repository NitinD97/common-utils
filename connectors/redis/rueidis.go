package redis

import (
	"github.com/redis/rueidis"
	"strconv"
)

func NewRueidisClient(cfg Config) (rueidis.Client, error) {
	client, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress: []string{cfg.Host + ":" + strconv.Itoa(cfg.Port)},
		Username:    "",
		Password:    cfg.Password,
		SelectDB:    cfg.Db,
	})
	if err != nil {
		return nil, err
	}
	return client, nil
}
