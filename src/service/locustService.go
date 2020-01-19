package service

import (
	"github.com/influxdata/influxdb-client-go"
	"github.com/sirupsen/logrus"
	"github.com/zalopay-oss/benchmark/configs"
	"github.com/zalopay-oss/benchmark/repository"
	"github.com/zalopay-oss/benchmark/utils"
)

var dbClient *influxdb.Client

func GetResult(config *configs.CannonConfig, id []byte) {
	distributedResult, err := utils.GetDistributedFile(config)
	if err != nil {
		logrus.Fatal(err)
	}
	requestResult, err := utils.GetRequestsFile(config)
	if err != nil {
		logrus.Fatal(err)
	}
	if err = repository.Save(dbClient, string(id), config, distributedResult, requestResult); err != nil {
		logrus.Fatal(utils.WrapError(err))
	}
}

func Start(config *configs.CannonConfig) {
	// TODO: get locust status
	msg, err := utils.GetLocustStatus(config)
	if err != nil {
		utils.Log(logrus.FatalLevel, err, "Fail GetLocustStatus")
	}

	if msg.State != "stopping" {
		err = utils.StartLocust(config)
		if err != nil {
			utils.Log(logrus.FatalLevel, err, "Fail StartLocust")
		}
	}

	if config.IsPersistent {
		dbClient, _ = influxdb.New(config.DatabaseAddr, config.Token)
	}

	utils.Log(logrus.InfoLevel, nil, "Start Cannon service success")
}

func Stop(config *configs.CannonConfig) {
	utils.CloseLocust(config)

	if dbClient != nil {
		dbClient.Close()
	}
}
