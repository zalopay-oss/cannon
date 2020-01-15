package utils

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/tranndc/benchmark/configs"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func StartLocust(config *configs.ServiceConfig) error{
	logrus.Info("START LOCUST")
	data := "locust_count="+strconv.Itoa(config.NoWorkers)+"&hatch_rate="+strconv.Itoa(config.HatchRate)
	des := config.Locust+"/swarm"
	req, err := http.NewRequest("POST",des,strings.NewReader(data))
	if err!=nil{
		return err
	}
	req.Header.Set("Connection","keep-alive")
	req.Header.Set("X-Requested-With","XMLHttpRequest")
	req.Header.Set("Content-Type","application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Accept","*/*")
	req.Header.Set("User-Agent","Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.88 Safari/537.36")
	req.Header.Set("Accept-Language","gzip, deflate")
	req.Header.Set("Accept-Encoding","*/*")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	logrus.Info(string(body))
	return nil
}

func CloseLocust(config *configs.ServiceConfig) error{
	logrus.Info("CLOSE LOCUST")
	des := config.Locust+"/stop"
	req, err := http.NewRequest("GET",des,nil)
	if err!=nil{
		return err
	}
	req.Header.Set("Connection","keep-alive")
	req.Header.Set("Accept","*/*")
	req.Header.Set("X-Requested-With","XMLHttpRequest")
	req.Header.Set("User-Agent","Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.88 Safari/537.36")
	req.Header.Set("Accept-Encoding","*/*")
	req.Header.Set("Accept-Language","gzip, deflate")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	logrus.Info(string(body))
	return nil
}



func GetDistributedFile(config *configs.ServiceConfig) (map[string]string, error) {
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
	return parseMessage(body),nil
}


func GetRequestsFile(config *configs.ServiceConfig) (map[string]string, error) {
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
	return parseMessage(body), nil
}
