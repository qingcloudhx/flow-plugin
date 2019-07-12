package timerrand

import (
	"github.com/qingcloudhx/core/data/coerce"
	"math/rand"
)

/**
* @Author: hexing
* @Date: 19-7-12 下午2:36
 */

func build(settings *HandlerSettings) *Output {
	out := &Output{
		DeviceId: settings.DeviceId,
		ThingId:  settings.ThingId,
	}
	for _, v := range settings.Device {
		if m, err := coerce.ToObject(v); err != nil {

		} else {
			s := ThingData{
				Id:   m["id"].(string),
				Name: m["name"].(string),
				Type: m["type"].(string),
			}
			switch m["type"] {
			case "float":
				s.Value = rand.Float64() + float64(rand.Intn(10))
			case "int32":
				s.Value = rand.Intn(100)
			case "string":
				s.Value = m["value"]
			default:

			}
			out.Device = append(out.Device, s)
		}
	}

	return out
}
