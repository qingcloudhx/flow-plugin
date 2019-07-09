package mqtts

import (
	"fmt"
	"github.com/muesli/cache2go"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

/**
* @Author: hexing
* @Date: 19-7-5 下午3:04
 */
func TestAdd(t *testing.T) {
	add("test", 5*time.Second, "xxxxxxxxx", func(item *cache2go.CacheItem) {
		fmt.Println("del", item.Data())
	})
	time.Sleep(6 * time.Second)
	value, err := get("test")
	assert.Nil(t, err)
	t.Log(value.Data())
	time.Sleep(6 * time.Second)

}
