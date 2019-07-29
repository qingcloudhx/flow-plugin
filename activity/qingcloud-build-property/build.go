package qingcloud_build_property

import (
	"github.com/qingcloudhx/core/data/coerce"
	"math/rand"
	"time"
)

/**
* @Author: hexing
* @Date: 19-7-12 下午2:36
 */

func build(settings *Settings) *Output {
	out := &Output{Device: make(map[string]interface{})}
	for k, v := range settings.Device {
		vm, _ := coerce.ToObject(v)
		out.Device[k] = make(map[string]interface{})
		msg := make(map[string]interface{})
		if m, err := coerce.ToObject(vm); err != nil {
		} else {
			msg["id"] = m["id"]
			msg["type"] = m["type"]
			msg["time"] = time.Now().Unix() * 1000
			switch m["type"] {
			case "float":
				msg["value"] = rand.Float64() + float64(rand.Intn(10))
			case "int32":
				msg["value"] = rand.Intn(100)
			case "string":
				msg["value"] = m["value"]
			default:

			}
			out.Device[k] = msg
		}
	}

	return out
}
