package repository

import (
	"context"
	"github.com/influxdata/influxdb-client-go"
	"github.com/sirupsen/logrus"
	"github.com/tranndc/benchmark/configs"
	"github.com/tranndc/benchmark/model"
	"time"
)

func Save(cli *influxdb.Client, resId string, config *configs.ServiceConfig, distributedData map[string]string, requestData map[string]string) error {
	logrus.Info(model.GetFields(config, distributedData, requestData))
	myMetrics := []influxdb.Metric{
		influxdb.NewRowMetric(
			model.GetFields(config, distributedData, requestData),
			"benchmark-results",
			model.GetTags(resId),
			time.Now(),
		)}
	if _,err := cli.Write(context.Background(),config.Bucket,config.Origin,myMetrics...); err!=nil{
		return err
	}
	return nil
}