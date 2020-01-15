package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/jsonpb"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"strings"
	"text/template"
)

func executeData(data string) ([]byte, error) {
	t := template.Must(template.New("").Parse(data))
	var m interface{}
	var tpl bytes.Buffer

	if err := json.Unmarshal([]byte(data), &m); err != nil {
		return []byte{}, nil
	}
	if err := t.Execute(&tpl, m); err != nil {
		return []byte{}, nil
	}
	return tpl.Bytes(),nil
}


func messageFromMap(input *dynamic.Message, data *map[string]interface{}) error {
	strData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = jsonpb.UnmarshalString(string(strData), input)
	if err != nil {
		return err
	}

	return nil
}


func createPayloadsFromJSON(data string, mtd *desc.MethodDescriptor) ([]*dynamic.Message, error) {
	md := mtd.GetInputType()
	var inputs []*dynamic.Message

	if len(data) > 0 {
		if strings.IndexRune(data, '[') == 0 {
			dataArray := make([]map[string]interface{}, 5)
			err := json.Unmarshal([]byte(data), &dataArray)
			if err != nil {
				return nil, fmt.Errorf("Error unmarshalling payload. Data: '%v' Error: %v", data, err.Error())
			}

			elems := len(dataArray)
			if elems > 0 {
				inputs = make([]*dynamic.Message, elems)
			}

			for i, elem := range dataArray {
				elemMsg := dynamic.NewMessage(md)
				err := messageFromMap(elemMsg, &elem)
				if err != nil {
					return nil, fmt.Errorf("Error creating message: %v", err.Error())
				}

				inputs[i] = elemMsg
			}
		} else {
			inputs = make([]*dynamic.Message, 1)
			inputs[0] = dynamic.NewMessage(md)
			err := jsonpb.UnmarshalString(data, inputs[0])
			if err != nil {
				return nil, fmt.Errorf("Error creating message from data. Data: '%v' Error: %v", data, err.Error())
			}
		}
	}

	return inputs, nil
}


func GetInputs(md *desc.MethodDescriptor, data string) ([]*dynamic.Message, error){
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	jsonStr := strings.TrimRight(data, "\n")
	inputData, err := executeData(jsonStr)
	if err != nil {
		logrus.Error(err)
	}

	return createPayloadsFromJSON(string(inputData), md)
}
