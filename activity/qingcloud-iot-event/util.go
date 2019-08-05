package qingcloud_iot_event

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/satori/go.uuid"
)

/**
* @Author: hexing
* @Date: 19-7-4 上午10:30
 */
const (
	device_id = "orgi"
	thing_id  = "thid"
)

// return token payload
func parseToken(tokenString string) (string, string, error) {
	//key, err := jwt.ParseRSAPublicKeyFromPEM([]byte("1212"))
	//if err != nil {
	//	return "1212", "", err
	//}
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return "", nil
	})
	payload := token.Claims.(jwt.MapClaims)
	if err := payload.Valid(); err != nil {
		return "", "", err
	}
	id, ok := payload[device_id].(string)
	if !ok {
		return "", "", errors.New("device id type error")
	}
	thingId, ok := payload[thing_id].(string)
	if !ok {
		return "", "", errors.New("device id type error")
	}
	return id, thingId, nil
}
func buildUpTopic(id, thingId, eventId string) string {
	return fmt.Sprintf("/sys/%s/%s/thing/event/%s/post", thingId, id, eventId)
}
func buildMessage(data interface{}) []byte {
	message := make(map[string]interface{})
	message["id"] = uuid.NewV4().String()
	message["version"] = "v1.0.0"
	message["params"] = data
	res, _ := json.Marshal(message)
	return res
}
