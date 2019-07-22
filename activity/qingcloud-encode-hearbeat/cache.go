package qingcloud_encode_hearbeat

/**
* @Author: hexing
* @Date: 19-7-5 下午2:56
 */
import (
	"github.com/muesli/cache2go"
	"time"
)

var cache *cache2go.CacheTable

func init() {
	cache = cache2go.Cache("iot-device-hearbeat")
}
func add(key interface{}, time time.Duration, value interface{}, f func(*cache2go.CacheItem)) {
	cache.Add(key, time, value)
	cache.SetAboutToDeleteItemCallback(f)
}

func exists(key interface{}) bool {
	return cache.Exists(key)
}
