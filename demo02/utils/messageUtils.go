package utils

import (
	"bytes"
	"encoding/csv"
	"io"
	"log"
)

func parseMessage(body []byte) map[string]string {
	reader := csv.NewReader(bytes.NewReader(body))
	var data = make([][]string,0,3)
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		data = append(data, line)
	}
	return parseKV(data[0],data[1])
}

func parseKV(key []string, value []string)map[string]string{
	size:= len(key)
	res := make(map[string]string)
	for i:=0;i<size;i++{
		res[key[i]] = value[i]
	}
	return res
}

