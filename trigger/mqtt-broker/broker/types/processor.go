//Created by zhbinary on 2019-03-20.
//Email: zhbinary@gmail.com
package types

import "github.com/256dpi/gomqtt/packet"

type Processor interface {
	Process(generic packet.Generic)
}
