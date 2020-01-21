package slave

import (
	"context"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/jhump/protoreflect/desc"
	"github.com/myzhan/boomer"
	"github.com/sirupsen/logrus"
	"github.com/zalopay-oss/benchmark/generator"
	"github.com/zalopay-oss/benchmark/generator/parser"
	"github.com/zalopay-oss/benchmark/utils"
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
		utils.PrintBanner("STOP SLAVE")
	})

	wg.Wait()
}

func (slave *Slave) RunTask(waitRun *sync.WaitGroup) {
	slaveBoomer = boomer.NewBoomer(slave.config.LocustHost, slave.config.LocustPort)
	var err error

	if slave.config.Proto == "" {
		utils.Log(logrus.FatalLevel, nil, "You must set proto file")
		os.Exit(1)
	}

	md, err, fd = parser.GetMethodDescFromProto(slave.config.Method, slave.config.Proto, []string{})

	if err != nil {
		utils.Log(logrus.FatalLevel, nil, "Invalid method: "+err.Error()+". Use -m to set methodName")
		os.Exit(1)
	}

	task := &boomer.Task{
		Name:   slave.config.Method,
		Weight: 1,
		Fn:     slave.Invoke,
	}

	if err = boomer.Events.Subscribe("boomer:hatch", func(workers int, hatchRate float64) {
		err := slave.CreateStubPool(workers)
		if err != nil {
			utils.Log(logrus.FatalLevel, err, "Cannot init pool")
		}
		atomic.AddInt32(&startTest, 1)
		logrus.Info("The master asks me to spawn ", workers, " goroutines with a hatch rate of ", int(hatchRate), " per second.")
	}); err != nil {
		utils.Log(logrus.FatalLevel, err, "Subcribe locust fail")
	}

	utils.PrintBanner("START SLAVE attack " + slave.config.GRPCHost + ":" + strconv.Itoa(slave.config.GRPCPort))
	slaveBoomer.Run(task)
	waitRun.Done()
	waitForQuit()

}

func (slave *Slave) Invoke() {
	_, _ = slave.invoke()
}

func (slave *Slave) invoke() (proto.Message, error) {
	for startTest == 0 {
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
		logrus.Error("Error call service: %v", err.Error())
		slaveBoomer.RecordFailure("tcp", call, elapsed.Nanoseconds()/int64(time.Millisecond), err.Error())
		return nil, err
	} else {
		slaveBoomer.RecordSuccess("tcp", call, elapsed.Nanoseconds()/int64(time.Millisecond), int64(len(res.String())))
	}
	return res, nil
}
