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
	out := &Output{}
	for _, v := range settings.Device {
		vm, _ := coerce.ToObject(v)
		out.Device[vm["name"].(string)] = make(map[string]interface{})
		msg := make(map[string]interface{})
		if m, err := coerce.ToObject(v); err != nil {
		} else {
			msg["id"] = m["id"]
			msg["type"] = m["type"]
			msg["time"] = time.Now().Unix() * 1000
			switch m["type"] {
			case "float":
				m["value"] = rand.Float64() + float64(rand.Intn(10))
			case "int32":
				m["value"] = rand.Intn(100)
			case "string":
				m["value"] = m["value"]
			default:

			}
		}
	}

	return out
}
