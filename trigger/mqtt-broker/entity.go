package mqtt

/**
* @Author: hexing
* @Date: 19-6-28 下午5:31
 */

type DeviceUpStatusMsg struct {
	Status string `json:"status"`
	//UserId     string `json:"user_id"`
	ThingId    string `json:"thing_id"`
	PropertyId string `json:"property_id"`
	DeviceId   string `json:"device_id"`
	Time       int64  `json:"time"`
}
