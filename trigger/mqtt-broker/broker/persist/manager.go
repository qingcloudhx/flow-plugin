//Created by zhbinary on 2019-03-13.
//Email: zhbinary@gmail.com
package persist

import (
	"context"
	"github.com/qingcloudhx/flow-plugin/trigger/mqtt-broker/broker/types"
	"sync"
)

type Manager struct {
	ch      chan interface{}
	factory types.PluginFactory
	ctx     context.Context
	cancel  context.CancelFunc
	wg      *sync.WaitGroup
}

func (this *Manager) start() {
	this.ch = make(chan interface{}, 100)
	defer close(this.ch)

	this.ctx, this.cancel = context.WithCancel(context.Background())

	plugin := this.factory.Create()
	err := plugin.Start(this.ctx, this.ch)
	if err != nil {
		this.wg.Add(1)
	}

	//nums := runtime.NumCPU()
	//for i := 0; i < nums; i++ {
	//	plugin := this.factory.Create()
	//	go func() {
	//		err := plugin.Start(this.ctx, this.ch)
	//		if err != nil {
	//			this.wg.Add(1)
	//		}
	//	}()
	//}
}

func (this *Manager) close() {
	this.cancel()
	this.wg.Wait()
}
