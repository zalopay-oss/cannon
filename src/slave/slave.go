package slave

import (
	"github.com/sirupsen/logrus"
	"github.com/zalopay-oss/benchmark/configs"
)

type Slave struct {
	Pool   *StubPool
	config *configs.CannonConfig
}

func CreateSlave(config *configs.CannonConfig) (*Slave, error) {
	slave := &Slave{
		config: config,
	}
	return slave, nil
}

func (slave *Slave) CreateStubPool(noConns int) error {
	pool, err := NewStubPool(noConns, slave.config)
	if err != nil {
		logrus.WithFields(logrus.Fields{"Error": err}).Fatal("Did not create pool")
		return err
	}
	slave.Pool = pool
	return nil
}
