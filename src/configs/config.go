package configs

import (
	"github.com/spf13/viper"
	"path/filepath"
	"strings"
)

//CannonConfig struct
type CannonConfig struct {
	NoWorkers int
	HatchRate int

	LocustWebTarget string
	LocustHost      string
	LocustPort      int

	IsPersistent bool
	DatabaseAddr string
	Token        string
	Measurement  string
	Bucket       string
	Origin       string

	GRPCPort int
	GRPCHost string

	Method string
	Proto  string

	ConfigFile string
}

func NewDefaultCannonConfig() *CannonConfig {
	return &CannonConfig{
		NoWorkers: 10,
		HatchRate: 10,

		LocustWebTarget: "http://0.0.0.0:8089/",
		LocustHost:      "localhost",
		LocustPort:      5557,

		IsPersistent: false,
		DatabaseAddr: "",
		Token:        "",
		Measurement:  "",
		Bucket:       "",
		Origin:       "",

		GRPCHost: "localhost",
		GRPCPort: 8000,
		Method:   "transRecord",
		Proto:    "transaction.proto",
	}
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
