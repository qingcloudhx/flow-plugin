//Created by zhbinary on 2019-03-12.
//Email: zhbinary@gmail.com
package session

import (
	"context"
	"github.com/256dpi/gomqtt/packet"
	"github.com/256dpi/gomqtt/transport"
)

type Manager struct {
	Ch       chan interface{}
	qos1Chan chan *packet.Publish
	context  context.Context
}

func NewManager() *Manager {
	return &Manager{}
}

func (this *Manager) Run() error {
	server, err := transport.Launch(this.Url)
	if err != nil {
		panic(err)
	}

	for {
		// accept next connection
		conn, err := server.Accept()
		if err != nil {
			return err
		}
		conn.Receive()
		this.handle(conn)
	}
}

func (this *Manager) handle(conn transport.Conn) {

}

// fan-in  fan-out
func (this *Manager) run() {
	for {
		select {
		case msg := <-this.Ch:
			// Send to persist
		case <-this.context.Done():
			return
		default:

		}
	}
}
