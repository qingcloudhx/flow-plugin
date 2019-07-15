package qingcloud_mqtt

/**
* @Author: hexing
* @Date: 19-7-15 下午2:06
 */

const (
	mqtt_client_id = "clienId"
	mqtt_poctocol  = "poctocol"
	mqtt_cmd       = "cmd"
	mqtt_username  = "username"
	mqtt_password  = "password"
	mqtt_topic     = "topic"
	mqtt_qos       = "qos"
	mqtt_head      = "head"
	mqtt_body      = "body"
)
const (
	mqtt_cmd_connecting = "connecting"
	mqtt_cmd_connect    = "connected"
	mqtt_cmd_disconnect = "disconnected"
	mqtt_cmd_data       = "data"
)

//create head
func buildHead(cmd, clientId, username, password string) map[string]interface{} {
	head := make(map[string]interface{})
	head[mqtt_client_id] = clientId
	head[mqtt_cmd] = cmd
	head[mqtt_poctocol] = "mqtt"
	head[mqtt_username] = username
	head[mqtt_password] = password
	return head
}

//create package
func buildPackage(head map[string]interface{}, body []byte) *Output {
	out := &Output{
		Head: head,
		Body: body,
	}
	return out
}
