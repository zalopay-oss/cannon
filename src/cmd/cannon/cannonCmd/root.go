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

// master variable
var noUsers int
var hatchRate int

var configFile string
var cannonConfig *configs.CannonConfig

// slave variable

var proto string
var method string
var grpcHost string
var grpcPort int

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
		} else {
			cannonConfig.HatchRate = hatchRate
			cannonConfig.NoWorkers = noUsers
			cannonConfig.Proto = proto
			cannonConfig.Method = method
			cannonConfig.GRPCHost = grpcHost
			cannonConfig.GRPCPort = grpcPort
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
	masterCmd.PersistentFlags().IntVarP(&hatchRate, "hatchRate", "r", cannonConfig.HatchRate, "config Hatch rate (users spawned/second)")
	masterCmd.PersistentFlags().IntVarP(&noUsers, "no-workers", "w", cannonConfig.NoWorkers, "Number of workers to simulate")
	masterCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "Config file")

	// slave
	masterCmd.PersistentFlags().StringVarP(&method, "method", "m", cannonConfig.Method, "Method name")
	masterCmd.PersistentFlags().StringVarP(&proto, "proto", "p", cannonConfig.Proto, "Proto File")
	masterCmd.PersistentFlags().StringVar(&grpcHost, "host", cannonConfig.GRPCHost, "Config gRPC host")
	masterCmd.PersistentFlags().IntVar(&grpcPort, "port", cannonConfig.LocustPort, "Config gRPC port")
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
