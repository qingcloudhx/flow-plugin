package modelmqtt

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
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
func buildUpTopic(id, thingId string) string {
	return fmt.Sprintf("/sys/%s/%s/thing/event/ObjectDection/post", thingId, id)
}
