//Created by zhbinary on 2019-03-12.
//Email: zhbinary@gmail.com
package persist

/*import (
	"context"
	"fmt"
	"github.com/256dpi/gomqtt/packet"
	"github.com/boltdb/bolt"
	"log"
	"github.com/qingcloudhx/flow-plugin/trigger/mqtt-broker/broker/types"
)

type BoltDbPlugin struct {
	ch      chan interface{}
	context context.Context
}

func (this *BoltDbPlugin) New(context context.Context, ch chan interface{}) types.Plugin {
	this.context = context
	this.ch = ch
	return nil
}

func (this *BoltDbPlugin) Start() error {
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("upstream_queue"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	for {
		select {
		case msg := <-this.ch:
			db.Update(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte("upstream_queue"))
				pub := msg.(*packet.Publish)
				//b.Put([]byte(pub.ID),pub.Message.Payload)
				return nil
			})
		case <-this.context.Done():
			break
		}
	}
}

func PubQos0Msg() {

}

func PubQos1Msg() {

}*/
