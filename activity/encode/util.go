package encode

import (
	"encoding/json"
	"fmt"
	"time"
)

/**
* @Author: hexing
* @Date: 19-7-5 下午3:44
 */
const (
	IOT_DEVICE_STATUS_END  = "iote-global-onoffline-end"
	IOT_DEVICE_STATUS_EDGE = "iote-global-onoffline-edge"
	DEVICE_STATUS_ONLINE   = "online"  // 在线
	DEVICE_STATUS_OFFLINE  = "offline" // 离线
)

func buildHeartBeat(id, thingId, status string) []byte {
	msg := &DeviceUpStatusMsg{
		DeviceId:   id,
		ThingId:    thingId,
		PropertyId: IOT_DEVICE_STATUS_END,
		Time:       time.Now().Unix(),
		Status:     status,
	}
	data, _ := json.Marshal(msg)
	return data
}
func buildUpTopic(id, thingId, eventId string) string {
	return fmt.Sprintf("/sys/%s/%s/thing/event/%s/post", thingId, id, eventId)
}

func buildHeartbeatTopic(id, thingId string) string {
	return fmt.Sprintf("/as/mqtt/status/%s/%s", thingId, id)
}
