//Created by zhbinary on 2019-04-12.
//Email: zhbinary@gmail.com
package model

import (
	"github.com/256dpi/gomqtt/packet"
	"github.com/qingcloudhx/flow-plugin/trigger/mqtt-broker/broker/types"
)

type Message struct {
	Pkt packet.Generic
	Ack types.AckCallback
}
