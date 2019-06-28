//Created by zhbinary on 2019-03-19.
//Email: zhbinary@gmail.com
package types

import (
	"github.com/256dpi/gomqtt/packet"
	"github.com/pkg/errors"
)

type Event string

const (
	EventSessionActive   Event = "Session active"
	EventSessionInactive Event = "Session inactive"
	EventSessionDequeue  Event = "Session dequeue"
)

var (
	ConnectError = errors.New("Connect error")
)

type Session interface {
	// NextID should return the next id for outgoing packets.
	NextID() packet.ID

	Close()
}
