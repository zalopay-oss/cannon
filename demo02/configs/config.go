package configs


import (
	"github.com/spf13/viper"
)

const configFilePath = "./configs"
const configFileName = "config"

//ServiceConfig struct
type ServiceConfig struct {
	GRPCPort     int
	GRPCHost     string
	NoConns      int
	NoWorkers    int
	Service      string
	Proto        string
	Data         string
	HatchRate    int
	Locust       string
	DatabaseAddr string
	Token        string
	Measurement  string
	Bucket       string
	Origin		 string
}

//LoadConfig load configs
func LoadConfig() error {
	viper.SetConfigName(configFileName)
	viper.AddConfigPath(configFilePath)
	viper.SetConfigType("yaml")

	return viper.ReadInConfig()
}