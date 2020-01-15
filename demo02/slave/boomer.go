package slave

import (
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/myzhan/boomer"
	"github.com/sirupsen/logrus"
	"github.com/tranndc/benchmark/configs"
	"github.com/tranndc/benchmark/generator"
)



func RunTask(slave *Slave, config *configs.ServiceConfig, done chan bool){
	var count = make(chan bool)

	task:= &boomer.Task{
		Name:	config.Service,
		Weight: 1,
		Fn: slave.Invoke,
	}

	boomer.Events.Subscribe("boomer:hatch", func(workers int, hatchRate float64) {
		logrus.Info("The master asks me to spawn ", workers, " goroutines with a hatch rate of", int(hatchRate), "per second.")
	})

	boomer.Events.Subscribe("boomer:quit", func() {
		logrus.Info("Boomer is quitting now.")
	})

	logrus.Info("START SLAVE")
	boomer.Run(task)

	done <- true
	count <- true
}

func (slave Slave) Invoke() {
	_, err := slave.invoke(slave.config.Service, slave.config.Proto)
	if err!=nil{
		logrus.Fatal(err)
	}
}


func (slave Slave) invoke(call string , proto string) (proto.Message, error){
	ctx := context.Background()
	inputs, md, err := generator.GetInput(proto, call)

	if err != nil {
		logrus.Error("Error creating client connection: %+v", err.Error())
		return nil, err
	}
	start := boomer.Now()

	stub, err := slave.Pool.Get()
	if err != nil {
		logrus.Error("Error getting stub: %+v", err.Error())
		return nil, err
	}
	res, err := stub.InvokeRpc(ctx, md, inputs[0])
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
