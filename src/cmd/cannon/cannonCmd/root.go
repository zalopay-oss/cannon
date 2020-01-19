package cannonCmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zalopay-oss/benchmark/configs"
	"github.com/zalopay-oss/benchmark/controller"
	"os"
)

var noUsers int
var hatchRate int

var configFile string
var config *configs.CannonConfig

var masterCmd = &cobra.Command{
	Use:   "run",
	Short: "Run Cannon",
	Long:  `Command run Cannon`,
	Run: func(cmd *cobra.Command, args []string) {
		if configFile != configs.DefaultCannonConfiguration {
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

func loadConfigFile() {
	config = &configs.CannonConfig{}
	if err := configs.LoadDefaultCannonConfig(); err != nil {
		logrus.Fatal("Load config: ", err)
		os.Exit(1)
	}
	if err := viper.Unmarshal(config); err != nil {
		logrus.Fatal("Load config: ", err)
		os.Exit(1)
	}
}

func initFlagsMasterCannon() {
	masterCmd.PersistentFlags().IntVarP(&hatchRate, "hatchRate", "r", config.HatchRate, "config Hatch rate (users spawned/second)")
	masterCmd.PersistentFlags().IntVarP(&noUsers, "no-workers", "w", config.NoWorkers, "Number of workers to simulate")
	masterCmd.PersistentFlags().StringVarP(&configFile, "config", "c", configs.DefaultCannonConfiguration, "Config file")
}

func Execute() {
	loadConfigFile()
	initFlagsMasterCannon()

	var rootCmd = &cobra.Command{Use: "cannon"}
	rootCmd.AddCommand(masterCmd)
	rootCmd.SetVersionTemplate("0.1.0")
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
