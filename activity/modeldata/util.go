package modeldata

import "fmt"

/**
* @Author: hexing
* @Date: 19-7-5 下午3:44
 */
func buildUpTopic(id, thingId string) string {
	return fmt.Sprintf("/sys/%s/%s/thing/event/property/post", thingId, id)
}
