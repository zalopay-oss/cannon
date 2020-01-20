package repository

import (
	"context"
	"github.com/influxdata/influxdb-client-go"
	"github.com/sirupsen/logrus"
	"github.com/zalopay-oss/benchmark/configs"
	"github.com/zalopay-oss/benchmark/model"
	"time"
)

func Save(cli *influxdb.Client, resId string, config *configs.CannonConfig, distributedData map[string]string, requestData map[string]string) error {
	if !config.IsPersistent {
		return nil
	}

	logrus.Info(model.GetFields(config, distributedData, requestData))
	myMetrics := []influxdb.Metric{
		influxdb.NewRowMetric(
			model.GetFields(config, distributedData, requestData),
			"benchmark-results",
			model.GetTags(resId),
			time.Now(),
		)}
	if _, err := cli.Write(context.Background(), config.Bucket, config.Origin, myMetrics...); err != nil {
		return err
	}
	return nil
}
