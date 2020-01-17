package slaveCmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zalopay-oss/benchmark/configs"
	"github.com/zalopay-oss/benchmark/slave"
	"os"
)

var proto string
var method string
var grpcHost string
var grpcPort int

var configFile string
var config *configs.SlaveConfig

var rootCmd = &cobra.Command{
	Use:   "run",
	Short: "Run Cannon",
	Long: `Command run Cannon`,
	Run: func(cmd *cobra.Command, args []string) {
		if configFile != configs.DefaultSlaveConfiguration {
			config = &configs.SlaveConfig{}
			if err := configs.LoadMyConfig(configFile); err != nil {
				logrus.Fatal("Load config: ", err)
				os.Exit(1)
			}
			if err := viper.Unmarshal(config); err != nil {
				logrus.Fatal("Load config: ", err)
				os.Exit(1)
			}
		}
		config.Proto = proto
		config.Method = method
		config.GRPCHost = grpcHost
		config.GRPCPort = grpcPort
		mSlave,err := slave.CreateSlave(config)
		if err!=nil{
			logrus.Fatal("Create Slave ", err)
		}
		mSlave.RunTask()
	},
}


func Execute() {
	config = &configs.SlaveConfig{}
	if err := configs.LoadDefaultSlaveConfig(); err != nil {
		logrus.Fatal("Load config: ", err)
		os.Exit(1)
	}
	if err := viper.Unmarshal(config); err != nil {
		logrus.Fatal("Load config: ", err)
		os.Exit(1)
	}
	rootCmd.PersistentFlags().StringVarP(&method, "method","m", config.Method , "Method name")
	rootCmd.PersistentFlags().StringVarP(&proto, "proto", "p", config.Proto, "Proto File")
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", configs.DefaultSlaveConfiguration, "Config file")
	rootCmd.PersistentFlags().StringVar(&grpcHost, "host", config.GRPCHost, "Config gRPC host")
	rootCmd.PersistentFlags().IntVar(&grpcPort, "port", config.LocustPort, "Config gRPC port")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
