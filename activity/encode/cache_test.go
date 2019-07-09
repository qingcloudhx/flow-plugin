package encode

import (
	"fmt"
	"testing"
	"time"
)

/**
* @Author: hexing
* @Date: 19-7-5 下午3:04
 */
func TestAdd(t *testing.T) {
	add("test", 5*time.Second, "xxxxxxxxx", func(key interface{}) {
		fmt.Println("del", key)
	})
	time.Sleep(4 * time.Second)
}
