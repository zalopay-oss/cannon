package slave

import (
	"errors"
	"fmt"
	"github.com/jhump/protoreflect/dynamic/grpcdynamic"
	"github.com/sirupsen/logrus"
	"github.com/tranndc/benchmark/configs"
	"google.golang.org/grpc"
	"sync"
)

type StubPool struct{
	stubs []grpcdynamic.Stub
	index uint64
	mtx sync.Mutex
}

func NewStubPool(size int,config *configs.ServiceConfig) (*StubPool, error){
	res := &StubPool{index:0}
	res.stubs = make([]grpcdynamic.Stub,0,size)
	conns,err := openConnection(size, config)
	if err != nil {
		logrus.WithFields(logrus.Fields{"Error": err}).Fatal("Did not connect server")
		return nil,err
	}
	for i:=0;i<size;i++{
		stub := grpcdynamic.NewStub(conns[i])
		res.stubs = append(res.stubs,stub)
	}
	return res,nil
}

func (pool *StubPool)Get() (grpcdynamic.Stub,error){
	if len(pool.stubs)==0{
		return grpcdynamic.Stub{},errors.New("Empty pool")
	}

	pool.mtx.Lock()
	i:= pool.index + 1
	pool.index = uint64(int(i) % len(pool.stubs))
	i= pool.index
	pool.mtx.Unlock()

	i = uint64(int(i) % len(pool.stubs))
	res := pool.stubs[i]
	return res,nil
}

func openConnection(size int, config *configs.ServiceConfig) ([]*grpc.ClientConn, error) {

	address := fmt.Sprintf("%s:%d", config.GRPCHost, config.GRPCPort)
	res := make([]*grpc.ClientConn,0,size)
	for i:=0;i<size;i++{
		conn, err := grpc.Dial(address, grpc.WithInsecure())
		if err != nil {
			logrus.WithFields(logrus.Fields{"Error": err}).Fatal("Did not connect server")
		}
		res = append(res, conn)
	}
	return res, nil
}