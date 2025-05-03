package redis

type Config struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	Db       int    `json:"db"`
	PoolSize int    `json:"pool_size"`
}
