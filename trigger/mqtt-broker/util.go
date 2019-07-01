package mqtt_broker

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/dgrijalva/jwt-go"
)

/**
* @Author: hexing
* @Date: 19-6-28 下午3:43
 */

const (
	device_id = "orgi"
	thing_id  = "thid"
)

var (
	inc uint16
)

type ConInfo struct {
	// The clients client id.
	ClientID string `json:"client_id"`
	// The authentication username.
	Username string `json:"username"`
	// The authentication password.
	Password string `json:"password"`
}

func uuid() uint16 {
	inc++
	return inc
}
func buildUpTopic(id, thingId string) string {
	return fmt.Sprintf("/sys/%s/%s/thing/event/property/post", thingId, id)
}
func buildDownTopic(id, thingId string) string {
	return fmt.Sprintf("/sys/%s/%s/thing/event/property/post_reply", thingId, id)
}

//up status
func buildStatusToipc(id, thingId string) string {
	return fmt.Sprintf("/as/mqtt/status/%s/%s", thingId, id)
}
func parseUser(data []byte) string {
	info := &ConInfo{}
	if err := json.Unmarshal(data, info); err != nil {
		return ""
	}
	return info.Password
}

// return token payload
func parseToken(pubKey, tokenString string) (string, string, error) {
	key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(pubKey))
	if err != nil {
		return "", "", err
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		return "", "", err
	}
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
