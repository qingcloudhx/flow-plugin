package encode

/**
* @Author: hexing
* @Date: 19-7-4 上午11:43
 */
type EventData struct {
	Id    int64       `json:"id"`
	Name  string      `json:"name"`
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}
type ThingEventData struct {
	Id   string       `json:"id"`
	Time int64        `json:"time"`
	Data []*EventData `json:"data"`
}
type ThingEventMsg struct {
	Id      string          `json:"id"`
	Version string          `json:"version"`
	Params  *ThingEventData `json:"params"`
}

type ThingMsg struct {
	Id      string                `json:id`
	Version string                `json:version`
	Params  map[string]*ThingData `json:params`
}

type ThingData struct {
	Id    string      `json:"id"`
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
	Time  int64       `json:"time"`
}

type DeviceUpStatusMsg struct {
	Status     string `json:"status"`
	UserId     string `json:"user_id"`
	ThingId    string `json:"thing_id"`
	PropertyId string `json:"property_id"`
	DeviceId   string `json:"device_id"`
	Time       int64  `json:"time"`
}
type DeviceInfo struct {
	ThingId  string `md:"thingId"`
	DeviceId string `md:"deviceId"`
}
