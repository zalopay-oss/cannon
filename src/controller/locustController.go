package controller

import (
	"github.com/spf13/viper"
	. "github.com/tranndc/benchmark/configs"
	"github.com/tranndc/benchmark/service"
	"github.com/tranndc/benchmark/slave"
	"log"
	"os/exec"
	"time"
)

func Run(){
	var done = make(chan bool)
	var stop = make(chan bool)

	config := &ServiceConfig{}
	if err := LoadConfig(); err != nil {
		log.Fatal("Load config: ", err)
	}

	if err := viper.Unmarshal(config); err != nil {
		log.Fatal("Load config: ", err)
	}

	managerSlave,err := slave.NewSlaveManager(config)
	if err!=nil{
		log.Fatal("Create Slave ", err)
	}

	for i:=0;i<config.NoConns;i++{
		go slave.RunTask(managerSlave.GetSlave(i),config,done)
	}

	go func(){
		time.Sleep(2 * time.Second)
		service.Start(config)
		GetMetric(config, stop)
	}()

	for i:=0;i<config.NoConns;i++{
		<-done
	}

	stop<-true

	defer service.Stop(config)
}

func GetMetric(config *ServiceConfig, stop chan bool) {
	id, err := exec.Command("uuidgen").Output()
	if err != nil {
		log.Fatal(err)
	}

	for true {
		time.Sleep(2 * time.Second)
		service.GetResult(config,id)
		select {
		case _,_= <-stop:
			return
		default:
			continue
		}
	}
}
