package mqtts

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
	cache = cache2go.Cache("iot")
}
func add(key interface{}, time time.Duration, value interface{}, f func(item *cache2go.CacheItem)) {
	cache.Add(key, time, value)
	cache.SetAboutToDeleteItemCallback(f)
}

func get(key interface{}) (*cache2go.CacheItem, error) {
	return cache.Value(key)
}
