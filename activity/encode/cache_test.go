package encode

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

/**
* @Author: hexing
* @Date: 19-7-5 下午3:04
 */
func TestAdd(t *testing.T) {
	add("test", 5*time.Second, "xxxxxxxxx")
	time.Sleep(4 * time.Second)
	b := exists("test")
	assert.True(t, b)
	time.Sleep(2 * time.Second)
	b = exists("test")
	assert.False(t, b)
}
