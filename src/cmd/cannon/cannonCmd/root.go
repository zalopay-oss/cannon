package cannonCmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tranndc/benchmark/configs"
	"github.com/tranndc/benchmark/controller"
	"os"
)

var DefaultConfigFile = "./configs/config.yaml"

var noUsers int
var hatchRate int

var configFile string
var config *configs.CannonConfig

var rootCmd = &cobra.Command{
	Use:   "run",
	Short: "Run Cannon",
	Long: `Command run Cannon`,
	Run: func(cmd *cobra.Command, args []string) {
		if configFile != DefaultConfigFile {
			config = &configs.CannonConfig{}
			if err := configs.LoadMyConfig(configFile); err != nil {
				logrus.Fatal("Load config: ", err)
				os.Exit(1)
			}
			if err := viper.Unmarshal(config); err != nil {
				logrus.Fatal("Load config: ", err)
				os.Exit(1)
			}
		}
		config.HatchRate = hatchRate
		config.NoWorkers = noUsers
		controller.Run(config)
	},
}


func Execute() {
	config = &configs.CannonConfig{}
	if err := configs.LoadDefaultCannonConfig(); err != nil {
		logrus.Fatal("Load config: ", err)
		os.Exit(1)
	}
	if err := viper.Unmarshal(config); err != nil {
		logrus.Fatal("Load config: ", err)
		os.Exit(1)
	}
	rootCmd.PersistentFlags().IntVarP(&hatchRate, "hatchRate","r", config.HatchRate , "config Hatch rate (users spawned/second)")
	rootCmd.PersistentFlags().IntVarP(&noUsers, "no-workers", "w", config.NoWorkers, "Number of workers to simulate")
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", DefaultConfigFile, "Config file")
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
