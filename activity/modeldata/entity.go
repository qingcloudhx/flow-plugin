package modeldata

/**
* @Author: hexing
* @Date: 19-7-12 下午2:03
 */
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
