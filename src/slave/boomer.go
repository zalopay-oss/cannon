package slave

import (
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/jhump/protoreflect/desc"
	"github.com/myzhan/boomer"
	"github.com/sirupsen/logrus"
	"github.com/zalopay-oss/benchmark/generator"
	"github.com/zalopay-oss/benchmark/generator/parser"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

var startTest int32 = 0

var slaveBoomer *boomer.Boomer
var md *desc.MethodDescriptor
var fd *desc.FileDescriptor

func waitForQuit() {
	wg := sync.WaitGroup{}
	wg.Add(1)

	quitByMe := false
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		<-c
		quitByMe = true
		slaveBoomer.Quit()
		wg.Done()
	}()

	_ = boomer.Events.Subscribe("boomer:quit", func() {
		if !quitByMe {
			wg.Done()
		}
		logrus.Info("STOP SLAVE")
	})

	wg.Wait()
}

func (slave *Slave) RunTask() {
	slaveBoomer = boomer.NewBoomer(slave.config.LocustHost,slave.config.LocustPort)
	var err error
	md, err, fd = parser.GetMethodDescFromProto(slave.config.Method,slave.config.Proto,[]string{})

	if err!=nil{
		panic(err)
		os.Exit(1)
	}

	task:= &boomer.Task{
		Name:	slave.config.Method,
		Weight: 1,
		Fn: slave.Invoke,
	}

	boomer.Events.Subscribe("boomer:hatch", func(workers int, hatchRate float64) {
		err := slave.CreateStubPool(workers)
		if err!=nil{
			panic(err)
			os.Exit(1)
		}
		atomic.AddInt32(&startTest,1)
		logrus.Info("The master asks me to spawn ", workers, " goroutines with a hatch rate of ", int(hatchRate), " per second.")
	})

	logrus.Info("START SLAVE")
	slaveBoomer.Run(task)
	waitForQuit()

}

func (slave *Slave) Invoke() {
	_, err := slave.invoke()
	if err!=nil{
		logrus.Fatal(err)
	}
}

func (slave *Slave) invoke() (proto.Message, error){
	for startTest==0{
	}
	ctx := context.Background()
	call := slave.config.Method
	inputs, md, err := generator.GetInput(call, md, fd)

	if err != nil {
		logrus.Error("Error creating client connection: %v", err.Error())
		return nil, err
	}
	start := time.Now()
	stub, err := slave.Pool.Get()
	if err != nil {
		logrus.Error("Error getting stub: %v", err.Error())
		return nil, err
	}
	res, err := stub.InvokeRpc(ctx, md, inputs[0])
	elapsed := time.Since(start)

	if err != nil {
		logrus.Error("Error InvokeRpc: %v", err.Error())
		slaveBoomer.RecordFailure("tcp", call+" fail", elapsed.Nanoseconds()/int64(time.Millisecond), err.Error())
		return nil, err
	} else {
		slaveBoomer.RecordSuccess("tcp", call, elapsed.Nanoseconds()/int64(time.Millisecond), int64(len(res.String())))
	}
	return res,nil
}
