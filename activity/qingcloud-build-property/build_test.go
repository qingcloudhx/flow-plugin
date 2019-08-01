package qingcloud_build_property

import "testing"

/**
* @Author: hexing
* @Date: 19-8-1 下午5:40
 */
func TestBuild(t *testing.T) {
	s := make(map[string]interface{})
	s["xxxx"] = make(map[string]interface{})
	m := make(map[string]interface{})
	m["id"] = "1212"
	m["type"] = "int32"
	s["xxxx"] = m
	ss := &Settings{s}
	res := build(ss)
	t.Log(res)
}
