# Config Package

This package provides a utility to load configuration using [Viper](https://github.com/spf13/viper). It reads two configuration files:
1. `global.json`: A global configuration file whose path is specified by the `GLOBAL_CONFIG` environment variable.
2. Environment-specific configuration files (`development.json`, `production.json`, `staging.json`): These files are located in the `config` directory of the project and are selected based on the `ENV` environment variable.

## Usage

1. Set the `GLOBAL_CONFIG` environment variable to the path of the `global.json` file.
2. Set the `ENV` environment variable to one of `development`, `production`, or `staging`. If not set, it defaults to `development`.
3. Call `config.InitConfig()` to initialize the configuration.