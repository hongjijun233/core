package calcium

import (
	"sync"

	"gitlab.ricebook.net/platform/core/scheduler"
	"gitlab.ricebook.net/platform/core/scheduler/simple"
	"gitlab.ricebook.net/platform/core/store"
	"gitlab.ricebook.net/platform/core/store/etcd"
	"gitlab.ricebook.net/platform/core/types"
)

type Calcium struct {
	sync.Mutex
	store     store.Store
	config    types.Config
	scheduler scheduler.Scheduler
}

func New(config types.Config) (*Calcium, error) {
	store, err := etcdstore.NewKrypton(config)
	if err != nil {
		return nil, err
	}

	scheduler := &simplescheduler.Magnesium{}

	return &Calcium{store: store, config: config, scheduler: scheduler}, nil
}

func (c *Calcium) Run() {

}