package config

import (
	"log"
	"os"
	"strconv"
	"github.com/joho/godotenv"
)

func GetConfig()(*Config,error){
	err:=godotenv.Load()
	if err!=nil{
		log.Println("no env file found")
		return nil,err
	}
	config:=&Config{
		ServerConfig:ServerConfig{
			Port:getEnv("PORT",""),
			JwtSecret:getEnv("JWT_SECRET",""),
		},
		MySqlConfig: MySqlConfig{
			Port: getEnvInt("MYSQL_PORT",3306),
			User: getEnv("MYSQL_USER",""),
			Password: getEnv("MYSQL_PASS",""),
			Database: getEnv("MYSQL_DB",""),
			Host: getEnv("MYSQL_HOST",""),
		},
		GrpcServerConfig: GrpcServerConfig{
			Port: getEnv("POST_SERVICE_PORT",""),
			Host: getEnv("POST_SERVICE_HOST",""),
		},
	}
	return config,nil
}

func getEnvInt(key string,defaultVal int)int{
	if value,exists := os.LookupEnv(key) ; exists{
		result,err:=strconv.Atoi(value)
		if err==nil{
			return result;
		}
	}
	return defaultVal;
}

func getEnv(key string,defaultVal string)string{
	if value,exists := os.LookupEnv(key);exists {
		return value
	}
	return defaultVal;
}
