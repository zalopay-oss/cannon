# Research technical stack

## 1. Demo 1

### 1.0 Deliverable

- Parse proto
- Use default data to benchmark
- Run multiple slaves
- Each slave has 1 stub (1 connection)
- Use boomer and Locust to benchmark

### 1.1 Parse proto and apply default data in order to make input

- Checkout this library: [protoreflect](https://github.com/jhump/protoreflect)
- First, we need to get the file descriptor of proto file. 
- From file descriptor, we can also get a service descriptor due to a service call provided by user.  
- After getting service descriptor, we can get method descriptor which contains information about gRPC service call, Input/Output type for that call.  

```go
md, err, fileDesc := parser.GetMethodDescFromProto("service.KeyValueStoreService.GetKV", "./zpkv.proto", []string{})
```  

In order to Invoke gRPC service call, we have define Input type for that call under `proto.Message` type. Therefore, I implement a method that generate the `dynamic.Message` input from method descriptor gotten from the above step:  

```go
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

func main() {
    inputs, _ := createPayloadsFromJSON(jsonStr, md)
	err = jsonpb.UnmarshalString(jsonStr, inputs[0])
    checkError(err)
    res, err := stub.InvokeRpc(ctx, md, inputs[0])
	checkError(err)
	logrus.Info(res)
}
```

### 1.2 Run mutiple slaves

- Number of slaves is configed in `config.yaml`
- Each slave run in one seperate goroutines
- Each slave has one seperate `stub`. Each `stub` has one ClientConn to connect to server

```go
    type Slave struct{
        stub grpcdynamic.Stub
        config *configs.ServiceConfig
    }

    type SlaveManager struct {
        slaves []Slave
        config *configs.ServiceConfig
    }
```

```go
func (manager *SlaveManager) createSlaves() ([]Slave,error){
    n:= manager.config.NoConns //number of slaves
    res := make([]Slave,0,n)
    //open connections
    conns, err := manager.openConnections()
    if err!=nil{
        return nil,err
    }

    for i:=0;i<n;i++ {
        // create new stub
        stub := grpcdynamic.NewStub(conns[i])
        slave := Slave{
            stub:stub,
            config:manager.config,
        }
        res = append(res, slave)
    }
    return res,nil
}

// create connections to server
func (manager *SlaveManager) openConnections() ([]*grpc.ClientConn, error) {
    address := fmt.Sprintf("%s:%d", manager.config.GRPCHost, manager.config.GRPCPort)
    n:= manager.config.NoConns // number of connections
    res := make([]*grpc.ClientConn,0,n)
    for i:=0;i<n;i++{
        conn, err := grpc.Dial(address, grpc.WithInsecure())
        if err != nil {
            logrus.WithFields(logrus.Fields{"Error": err}).Fatal("Did not connect server")
        }
        res = append(res, conn)
    }
    return res, nil
}

func (manager *SlaveManager)GetSlave(index int) *Slave{
    return &manager.slaves[index]
}
```

```go
    for i:=0;i<config.NoConns;i++{
        go runTask(managerSlave.GetSlave(i),config)
    }
```

### 1.3 Invoke server method

- Get method and data from proto (1.1)
- Use `grpcdynamic.Stub.InvokeRpc(context, methodDescriptor, input)` to invoke server method

```go
func (slave Slave) invoke(call string , proto string, data string) (proto.Message, error){
    md, err := parser.GetMethodDescFromProto(call, proto, []string{})
    if err != nil {
        logrus.Error("Error parse Proto: %+v", err.Error())
        return nil, err
    }
    ctx := context.Background()

    //get input from proto and default data
    inputs, err := utils.GetInputs(md,data)

    if err != nil {
        logrus.Error("Error creating client connection: %+v", err.Error())
        return nil, err
    }
    //Invoke method
    res, err := slave.stub.InvokeRpc(ctx, md, inputs[0])

    if err != nil {
        return nil, err
    }
    return res,nil
}
```

### 1.3 Apply boomber and locust to benchmark

```go
func runTask(slave *Slave, config *configs.ServiceConfig){
    task:= &boomer.Task{
        Name: config.Service,
        Weight: 1,
        Fn: slave.Invoke,
    }

    boomer.Events.Subscribe("boomer:hatch", func(workers int, hatchRate float64) {
        logrus.Info("The master asks me to spawn ", workers, " goroutines with a hatch rate of", int(hatchRate), "per second.")
    })

    boomer.Events.Subscribe("boomer:quit", func() {
        logrus.Info("Boomer is quitting now.")
    })
    boomer.Run(task)

    done <- true
}
```

```go

func (slave Slave) Invoke() {
    _, err := slave.invoke(slave.config.Service, slave.config.Proto, slave.config.Data)
    if err!=nil{
        logrus.Fatal(err)
    }
}

func (slave Slave) invoke(call string , proto string, data string) (proto.Message, error){
    md, err := parser.GetMethodDescFromProto(call, proto, []string{})
    if err != nil {
        logrus.Error("Error parse Proto: %+v", err.Error())
        return nil, err
    }
    ctx := context.Background()

    inputs, err := utils.GetInputs(md,data)

    if err != nil {
        logrus.Error("Error creating client connection: %+v", err.Error())
        return nil, err
    }
    start := boomer.Now()
    res, err := slave.stub.InvokeRpc(ctx, md, inputs[0])
    elapsed := boomer.Now() - start

    if err != nil {
        logrus.Error("Error InvokeRpc: %+v", err.Error())
        boomer.RecordFailure("tcp", call+" fail", elapsed, err.Error())
        return nil, err
    } else {
        boomer.RecordSuccess("tcp", call, elapsed, int64(len(res.String())))
    }
    return res,nil
}
```

## 2. Demo 2

### 2.0 Deliverable

- Generate random data with given proto
- Use random data to benchmark
- Start Locust in run time
- Run multiple slaves, each slave has multiple connections
- Store benchmark results in influxdb

### 2.1 Generate random data with given proto  

Up to now, the system supports generating data for following types:  
- string
- int32
- int64
- sint32
- sint64
- fixed32
- fixed64
- sfixed32
- sfixed64
- float
- bool
- bytes
- map
- repeated
- enum
- oneof

The data is just random. Feature works: generate sequential data.  

### 2.2 Start, stop, and get results from locust in run time

- **START LOCUST** by calling api  

```bash
curl --location --request POST 'http://localhost:7000/swarm' \
--header 'Connection: keep-alive' \
--header 'Accept: */*' \
--data-raw 'locust_count=1000&hatch_rate=100'
```

- **STOP LOCUST** by calling api  

```bash
curl --location --request GET 'http://localhost:7000/stop' \
--header 'Connection: keep-alive' \
--header 'Accept: */*' \
```

- **GET BENCHMARK RESULTS** by calling api

```bash
curl --location --request GET 'http://localhost:7000/stats/distribution/csv' \
```

```bash
curl --location --request GET 'http://localhost:7000/stats/requests/csv' \
```

### 2.3 Store results in influxdb

- Tags:
  - `id` : each benchmark test has a unique ID
- Fields:
  - `p90`
  - `p95`
  - `p99`
  - `rps`: Request/second
  - `requests`: Number of requests
  - `max_res_time`: Maximum response time
  - `min_res_time`: Minimum response time
  - `median_res_time`: Median response time
  - `avg_res_time`: Average response time
  - `avg_res_size`: Average content size
  - `configs`: JSON
    - `service`: service's name
    - `proto`: link to proto
    - `hatch_rate`
    - `users`: number of users
    - `slaves`: number of slaves
  