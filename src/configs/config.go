package configs

import (
	"github.com/spf13/viper"
	"path/filepath"
	"strings"
)

const DefaultSlaveConfiguration = "./configs/default-slave-config.yaml"
const DefaultCannonConfiguration = "./configs/default-cannon-config.yaml"

//SlaveConfig struct
type SlaveConfig struct {
	GRPCPort      int
	GRPCHost      string
	Method        string
	Proto         string
	LocustWebPort string
	LocustHost    string
	LocustPort    int
}

//CannonConfig struct
type CannonConfig struct {
	NoWorkers     int
	HatchRate     int
	LocustWebPort string
	LocustHost    string
	LocustPort    int
	DatabaseAddr  string
	Token         string
	Measurement   string
	Bucket        string
	Origin        string
}


func LoadDefaultSlaveConfig() error{
	return LoadMyConfig(DefaultSlaveConfiguration)
}

func LoadDefaultCannonConfig() error{
	return LoadMyConfig(DefaultCannonConfiguration)
}

func LoadMyConfig(path string) error {
	dir := filepath.Dir(path)
	name := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	ext := filepath.Ext(path)
	ext = ext[1:]
	viper.SetConfigName(name)
	viper.AddConfigPath(dir)
	viper.SetConfigType(ext)

	return viper.ReadInConfig()
}