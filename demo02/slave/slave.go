package slave

import (
	"github.com/sirupsen/logrus"
	"github.com/tranndc/benchmark/configs"
)

type Slave struct{
	Pool *StubPool
	config *configs.ServiceConfig
	noWorkers int
}

type SlaveManager struct {
	slaves []Slave
	config *configs.ServiceConfig
}


func NewSlaveManager(config *configs.ServiceConfig) (*SlaveManager,error){
	manager := &SlaveManager{
		config: config,
	}
	slaves, err := manager.createSlaves()
	if err!=nil{
		return nil,err
	}
	manager.slaves = slaves
	return manager,nil
}

func (manager *SlaveManager)GetSlave(index int) *Slave{
	return &manager.slaves[index]
}

func (manager *SlaveManager) createSlaves() ([]Slave,error){
	n := manager.config.NoConns
	logrus.Warn("No:",n)
	res := make([]Slave,0,n)
	for i:=0;i<n;i++ {
		var noWorkers int
		if i<n-1 || n==1{
			noWorkers = manager.config.NoWorkers /n
		} else {
			normalWorker := manager.config.NoWorkers /n
			noWorkers = manager.config.NoWorkers - (n-1)*normalWorker
		}
		logrus.Warn("No:",noWorkers)
		pool,err := NewStubPool(noWorkers,manager.config)
		if err != nil {
			logrus.WithFields(logrus.Fields{"Error": err}).Fatal("Did not create pool")
			return nil,err
		}
		slave := Slave{
			Pool:    pool,
			config:  manager.config,
			noWorkers: noWorkers,
		}
		res = append(res, slave)
	}
	return res,nil
}



//func (manager *SlaveManager) openConnections() ([]*grpc.ClientConn, error) {
//	address := fmt.Sprintf("%s:%d", manager.config.GRPCHost, manager.config.GRPCPort)
//	n:= manager.config.NoConns
//	res := make([]*grpc.ClientConn,0,n)
//	for i:=0;i<n;i++{
//		conn, err := grpc.Dial(address, grpc.WithInsecure())
//		if err != nil {
//			logrus.WithFields(logrus.Fields{"Error": err}).Fatal("Did not connect server")
//		}
//		res = append(res, conn)
//	}
//	return res, nil
//}