package encode

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
func add(key interface{}, time time.Duration, value interface{}, f func(key interface{})) {
	res := cache.Add(key, time, value)
	res.SetAboutToExpireCallback(f)
}
