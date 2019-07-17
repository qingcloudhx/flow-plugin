package encode

import (
	"encoding/json"
	"fmt"
	"github.com/yunify/qingstor-sdk-go/config"
	//qsErrors "github.com/yunify/qingstor-sdk-go/request/errors"
	"github.com/qingcloudhx/core/support/log"
	qs "github.com/yunify/qingstor-sdk-go/service"
	"os"
	"strconv"
	"time"
)

/**
* @Author: hexing
* @Date: 19-7-5 下午3:44
 */
const (
	IOT_DEVICE_STATUS_END  = "iote-global-onoffline-end"
	IOT_DEVICE_STATUS_EDGE = "iote-global-onoffline-edge"
	DEVICE_STATUS_ONLINE   = "online"  // 在线
	DEVICE_STATUS_OFFLINE  = "offline" // 离线
)

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
func buildUpTopic(id, thingId, eventId string) string {
	return fmt.Sprintf("/sys/%s/%s/thing/event/%s/post", thingId, id, eventId)
}

func buildHeartbeatTopic(id, thingId string) string {
	return fmt.Sprintf("/as/mqtt/status/%s/%s", thingId, id)
}

var bucket *qs.Bucket

func init() {
	configuration, _ := config.New("LCZPYVBDRWGFQYNTLLHF", "nptZQsUpqc7UIdjqgjrL1GXSbapKTkNGsAPjPr1z")
	qsService, _ := qs.Init(configuration)
	bucket, _ = qsService.Bucket("facetest", "pek3a")
	//putBucketOutput, _ := bucket.Put()
}
func getPictureUrl(path string, logger log.Logger) string {
	// Open file
	//if file, err := os.Open(path); err != nil {
	//	panic(err)
	//}
	file, err := os.Open(path)
	if err != nil {
		logger.Errorf("open fail err:%s", err.Error())
		return ""
	}
	defer func() {
		file.Close()
		if err := os.Remove(path); err != nil {
			logger.Error(err)
		} else {
			logger.Infof("remove %s success", path)
		}
	}()

	// Put object
	name := "ai-" + strconv.FormatUint(uint64(time.Now().Unix()), 10) + ".jpg"
	oOutput, err := bucket.PutObject(name, &qs.PutObjectInput{Body: file})
	// 所有 >= 400 的 HTTP 返回码都被视作错误
	if err != nil {
		// Example: QingStor Error: StatusCode 403, Code "permission_denied"...
		logger.Errorf("PutObject fail err:%s", err.Error())
		return ""
	} else {
		// Print the HTTP status code.
		// Example: 201
		logger.Infof("up success code:%d", qs.IntValue(oOutput.StatusCode))
		return fmt.Sprintf("https://facetest.pek3b.qingstor.com/%s", name)
	}
}
