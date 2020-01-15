package service

import (
	"github.com/influxdata/influxdb-client-go"
	"github.com/sirupsen/logrus"
	"github.com/tranndc/benchmark/configs"
	"github.com/tranndc/benchmark/repository"
	"github.com/tranndc/benchmark/utils"
	"io/ioutil"
	"log"
	"net/http"
)

var dbClient *influxdb.Client

func GetResult(config *configs.ServiceConfig, id []byte) {


	distributedResult, err := getDistributedFile(config)
	if err != nil {
		logrus.Fatal(err)
	}
	requestResult, err := getRequestsFile(config)
	if err != nil {
		logrus.Fatal(err)
	}
	if err = repository.Save(dbClient, string(id), config, distributedResult, requestResult); err !=nil{
		logrus.Fatal(err)
	}
}


func getDistributedFile(config *configs.ServiceConfig) (map[string]string, error) {
	des := config.Locust+"/stats/distribution/csv"
	req, err := http.NewRequest("GET",des,nil)
	if err!=nil{
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err!=nil{
		return nil,err
	}
	return utils.ParseMessage(body),nil
}


func getRequestsFile(config *configs.ServiceConfig) (map[string]string, error) {
	des := config.Locust + "/stats/requests/csv"
	req, err := http.NewRequest("GET", des, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return utils.ParseMessage(body), nil
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