package configuration

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

var cfg *viper.Viper

var environments = map[string]bool{
	"development": true,
	"staging":     true,
	"production":  true,
}

func InitConfig() *viper.Viper {
	Config := viper.New()
	Config.SetConfigType("json")

	// Load global configuration
	globalConfigPath := os.Getenv("GLOBAL_CONFIG")
	if globalConfigPath == "" {
		log.Fatal("GLOBAL_CONFIG environment variable is not set")
	}
	Config.SetConfigFile(globalConfigPath)
	if err := Config.MergeInConfig(); err != nil {
		log.Fatalf("Error loading global configuration: %v", err)
	}

	env := Config.GetString("environment")
	if !environments[env] {
		log.Fatal("environment key is not set in the global configuration")
	}
	envConfigPath := "./config/" + env + ".json"
	Config.SetConfigFile(envConfigPath)
	if err := Config.MergeInConfig(); err != nil {
		log.Fatalf("Error loading %s configuration: %v", env, err)
	}

	cfg = Config
	return Config
}

func GetConfig() *viper.Viper {
	return cfg
}
