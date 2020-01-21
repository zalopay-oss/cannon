package utils

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/zalopay-oss/benchmark/configs"
	"github.com/zalopay-oss/benchmark/model"
)

func checkResponse(resp *http.Response) {
	var response model.LocustCommandResponse
	body, _ := ioutil.ReadAll(resp.Body)

	json.Unmarshal(body, &response)

	if !response.Success {
		Log(logrus.ErrorLevel, nil, "Call fail with msg: "+response.Message)
	}
}

// StartLocust call locust API to start the test
func StartLocust(config *configs.CannonConfig) error {
	logrus.Info("Starting the test...")
	data := "locust_count=" + strconv.Itoa(config.NoWorkers) + "&hatch_rate=" + strconv.Itoa(config.HatchRate)
	des := config.LocustWebTarget + "/swarm"
	req, err := http.NewRequest("POST", des, strings.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.88 Safari/537.36")
	req.Header.Set("Accept-Language", "gzip, deflate")
	req.Header.Set("Accept-Encoding", "*/*")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		Log(logrus.FatalLevel, err, "Error start locust")
	}
	defer resp.Body.Close()

	checkResponse(resp)
	return nil
}

func CloseLocust(config *configs.CannonConfig) error {
	logrus.Info("Stopping the test...")
	des := config.LocustWebTarget + "/stop"
	req, err := http.NewRequest("GET", des, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.88 Safari/537.36")
	req.Header.Set("Accept-Encoding", "*/*")
	req.Header.Set("Accept-Language", "gzip, deflate")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	checkResponse(resp)
	return nil
}

func GetDistributedFile(config *configs.CannonConfig) (map[string]string, error) {
	des := config.LocustWebTarget + "/stats/distribution/csv"
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

func GetRequestsFile(config *configs.CannonConfig) (map[string]string, error) {
	des := config.LocustWebTarget + "/stats/requests/csv"
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

func GetLocustStatus(config *configs.CannonConfig) (*model.LocustStatus, error) {
	des := config.LocustWebTarget + "/stats/requests"
	req, err := http.NewRequest("GET", des, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		Log(logrus.FatalLevel, err, "Cannot connect "+config.LocustWebTarget)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var status *model.LocustStatus

	json.Unmarshal(body, &status)
	return status, nil
}
