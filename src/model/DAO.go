package model

import (
	"encoding/json"
	"github.com/tranndc/benchmark/configs"
	"strconv"
)

const  REQUEST  = "# requests"
const  RPS  = "Requests/s"
const  FAILS  = "# failures"
const  P90  = "90%"
const  P95  = "95%"
const  P99  = "99%"
const  MedianResTime = "Median response time"
const  MaxResTime = "Max response time"
const  MinResTime = "Min response time"
const  AvgResTime = "Average response time"
const  AvgResSize = "Average Content Size"

func GetTags(resId string) map[string]string {
	return map[string]string{
		"id":resId[:len(resId)-1],

	}
}

func toFloat(value string)float64 {
	res,_ := strconv.ParseFloat(value,64)
	return res
}


func GetFields(config *configs.ServiceConfig, distributedData map[string]string, requestData map[string]string) map[string]interface{} {
	return map[string]interface{}{
		"configs":         getConfigField(config),
		"request":         toFloat(distributedData[REQUEST]),
		"fails":           toFloat(requestData[FAILS]),
		"rps":             toFloat(requestData[RPS]),
		"p90":             toFloat(distributedData[P90]),
		"p95":             toFloat(distributedData[P95]),
		"p99":             toFloat(distributedData[P99]),
		"max_res_time":    toFloat(requestData[MaxResTime]),
		"min_res_time":    toFloat(requestData[MinResTime]),
		"avg_res_time":    toFloat(requestData[AvgResTime]),
		"median_res_time": toFloat(requestData[MedianResTime]),
		"avg_res_size":    toFloat(requestData[AvgResSize]),
	}
}

func getConfigField(config *configs.ServiceConfig) string {
	res:= make(map[string]interface{})
	res["server"]=config.Service
	res["proto"]=config.Proto
	res["hatchRate"]=config.HatchRate
	res["users"]=config.NoWorkers
	res["slaves"]=config.NoConns
	stringRes, _ :=json.Marshal(res)
	return string(stringRes)
}
