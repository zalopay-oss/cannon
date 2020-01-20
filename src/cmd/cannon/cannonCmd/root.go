package cannonCmd

import (
	"os"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zalopay-oss/benchmark/configs"
	"github.com/zalopay-oss/benchmark/controller"
	"github.com/zalopay-oss/benchmark/slave"
	"github.com/zalopay-oss/benchmark/utils"
)

var configFile string
var cannonConfig *configs.CannonConfig

// command
var masterCmd = &cobra.Command{
	Use:   "run",
	Short: "Run Cannon",
	Long:  `Command run Cannon`,
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Infof("%s", configFile)

		if configFile != "" {
			cannonConfig = &configs.CannonConfig{}
			if err := configs.LoadMyConfig(configFile); err != nil {
				utils.Log(logrus.FatalLevel, err, "Load config")
			}
			if err := viper.Unmarshal(cannonConfig); err != nil {
				utils.Log(logrus.FatalLevel, err, "Parse config")
			}
		}

		waitRun := &sync.WaitGroup{}
		waitRun.Add(1)
		go func() {
			mSlave, err := slave.CreateSlave(cannonConfig)
			if err != nil {
				logrus.Fatal("Create Slave ", err)
			}
			mSlave.RunTask(waitRun)
			waitRun.Done()
		}()

		// start master
		waitRun.Wait()
		waitRun.Add(2)

		go func() {
			controller.Run(cannonConfig)
			waitRun.Done()
		}()

		waitRun.Wait()
	},
}

func initFlags() {
	cannonConfig = configs.NewDefaultCannonConfig()

	// init flags
	// master
	masterCmd.PersistentFlags().IntVarP(&cannonConfig.HatchRate, "hatchRate", "r", cannonConfig.HatchRate, "config Hatch rate (users spawned/second)")
	masterCmd.PersistentFlags().IntVarP(&cannonConfig.NoWorkers, "no-workers", "w", cannonConfig.NoWorkers, "Number of workers to simulate")
	masterCmd.PersistentFlags().StringVarP(&cannonConfig.ConfigFile, "config", "c", "", "Config file")

	// slave
	masterCmd.PersistentFlags().StringVarP(&cannonConfig.Method, "method", "m", cannonConfig.Method, "Method name")
	masterCmd.PersistentFlags().StringVarP(&cannonConfig.Proto, "proto", "p", cannonConfig.Proto, "Proto File")
	masterCmd.PersistentFlags().StringVarP(&cannonConfig.GRPCHost, "host", "H", cannonConfig.GRPCHost, "Target gRPC host")
	masterCmd.PersistentFlags().IntVarP(&cannonConfig.GRPCPort, "port", "P", cannonConfig.GRPCPort, "Target gRPC port")

	// locust
	masterCmd.PersistentFlags().StringVar(&cannonConfig.LocustHost, "locust-host", cannonConfig.LocustHost, "Host of Locust Master")
	masterCmd.PersistentFlags().IntVar(&cannonConfig.LocustPort, "locust-port", cannonConfig.LocustPort, "Port of Locust Master")
	masterCmd.PersistentFlags().StringVar(&cannonConfig.LocustWebTarget, "locust-web", cannonConfig.LocustWebTarget, "Locust Web Target")
}

func Execute() {
	initFlags()

	var rootCmd = &cobra.Command{Use: "cannon"}
	rootCmd.AddCommand(masterCmd)
	rootCmd.SetVersionTemplate("0.1.0")
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
