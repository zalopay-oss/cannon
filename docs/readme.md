# Research technical stack

## 1. Demo 1

### 1.0 Delevirable

- Parse proto
- Use default data to benchmark
- Run multiple slaves
- Each slave has 1 stub (1 connection)
- Use boomer and Locust to benchmark

### 1.1 Parse proto and apply default data in order to make input

QuyenPT3

### 1.2 Run mutiple slaves

- Number of slaves is configed in `config.yaml`
- Each slave run in one seperate goroutines
- Each slace has one seperate `stub`. Each `stub` has one ClientConn to connect to server

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
