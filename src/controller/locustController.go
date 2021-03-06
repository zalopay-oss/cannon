package controller

import (
	"github.com/sirupsen/logrus"
	"github.com/zalopay-oss/benchmark/configs"
	"github.com/zalopay-oss/benchmark/service"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

func Run(config *configs.CannonConfig) {
	var stop = make(chan bool)
	service.Start(config, stop)
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		logrus.Info("Shutting Down Cannon....")
		service.Stop(config)
		stop <- true
		os.Exit(1)
	}()
	GetMetric(config, stop)
}

func GetMetric(config *configs.CannonConfig, stop chan bool) {
	id, err := exec.Command("uuidgen").Output()
	if err != nil {
		log.Fatal(err)
	}
	for true {
		time.Sleep(2 * time.Second)
		service.GetResult(config, id)
		select {
		case _, _ = <-stop:
			return
		default:
			continue
		}
	}
}
