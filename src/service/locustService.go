package service

import (
	"github.com/influxdata/influxdb-client-go"
	"github.com/sirupsen/logrus"
	"github.com/tranndc/benchmark/configs"
	"github.com/tranndc/benchmark/repository"
	"github.com/tranndc/benchmark/utils"
	"log"
)

var dbClient *influxdb.Client

func GetResult(config *configs.ServiceConfig, id []byte) {


	distributedResult, err := utils.GetDistributedFile(config)
	if err != nil {
		logrus.Fatal(err)
	}
	requestResult, err := utils.GetRequestsFile(config)
	if err != nil {
		logrus.Fatal(err)
	}
	if err = repository.Save(dbClient, string(id), config, distributedResult, requestResult); err !=nil{
		logrus.Fatal(err)
	}
}

func Start(config *configs.ServiceConfig) {
	err := utils.StartLocust(config)
	if err != nil {
		log.Fatal("Start Locust ", err)
	}

	dbClient, _ = influxdb.New(config.DatabaseAddr, config.Token)
}
func Stop(config *configs.ServiceConfig){
	utils.CloseLocust(config)
	dbClient.Close()
}