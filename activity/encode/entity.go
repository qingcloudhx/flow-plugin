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
