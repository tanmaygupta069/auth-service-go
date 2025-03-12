package config

type Config struct{
	ServerConfig ServerConfig
	MySqlConfig MySqlConfig
	GrpcServerConfig GrpcServerConfig
}

type ServerConfig struct{
	Port string
	JwtSecret string
}

type MySqlConfig struct{
	Port int
    User string
    Password string
    Database string
	Host string
}

type GrpcServerConfig struct{
	Port string
	Host string
}