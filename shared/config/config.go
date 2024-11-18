package config

type Config struct {
	OrderServiceServer struct {
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
		DSN  string `mapstructure:"dsn"`
	} `mapstructure:"order_service"`

	UserServiceServer struct {
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
		DSN  string `mapstructure:"dsn"`
	} `mapstructure:"user_service"`

	LogLevel string `mapstructure:"loglevel"`
}

var defaultUserSvcConfig = struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	DSN  string `mapstructure:"dsn"`
}{
	Host: "0.0.0.0",
	Port: 5050,
	DSN:  "postgres://ordersvc:password@postgres:5432/order-svc",
}

var defaultOrderSvcConfig = struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	DSN  string `mapstructure:"dsn"`
}{
	Host: "0.0.0.0",
	Port: 5051,
	DSN:  "postgres://usersvc:password@postgres:5432/user-svc",
}

// default config
var defaultConfig = Config{
	OrderServiceServer: defaultOrderSvcConfig,
	UserServiceServer:  defaultUserSvcConfig,
	LogLevel:           "debug",
}

// TODO: setup config to load from env files or external key-value store
func setupConfig() {}

// global config object
var GlobalConfig = &defaultConfig
