//Created by zhbinary on 2019-04-12.
//Email: zhbinary@gmail.com
package types

import "context"

type Plugin interface {
	New() Plugin
	Start(context context.Context, ch chan interface{}) error
}

type PluginFactory interface {
	Create() Plugin
}
