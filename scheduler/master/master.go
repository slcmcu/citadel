package master

import (
	"errors"
	"sync"
	"time"

	"citadelapp.io/citadel"
	"citadelapp.io/citadel/utils"
	"github.com/Sirupsen/logrus"
)

var (
	ErrNoValidOffers = errors.New("no valid offers for tasks")
)

// Master is the master node in a cluster
type Master struct {
	sync.Mutex

	log     *logrus.Logger
	timeout time.Duration
}

func New(logger *logrus.Logger, timeout time.Duration) (*Master, error) {
	m := &Master{
		timeout: timeout,
		log:     logger,
	}
	return m, nil
}

func (m *Master) Schedule(task *citadel.Task) error {
	m.Lock()
	transactionId := utils.GenerateUUID(32)

	defer func() {
		m.log.WithField("transaction_id", transactionId).Debug("ending scheduled transaction")
		m.Unlock()
	}()
	m.log.WithField("transaction_id", transactionId).Debug("starting scheduled transaction")

	return nil
}
