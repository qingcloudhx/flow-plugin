package persist

import (
	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
	"github.com/qingcloudhx/flow-plugin/trigger/mqtt-broker/broker/types"
	"time"
)

type BoltDB struct {
	*bolt.DB
	bucket []byte
}

const BucketName = "edgewise"

var ErrBucketNotExist = errors.New("Bucket not exist")
var ErrKeyNotExist = errors.New("Key not exist")
var ErrNoRecord = errors.New("No record")

func NewBoltDB(path string) (types.Persist, error) {
	db, err := bolt.Open(path, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}

	return &BoltDB{db, []byte(BucketName)}, nil
}

func (o *BoltDB) Get(key []byte) ([]byte, error) {
	return o.BucketGet(o.bucket, key)
}

func (o *BoltDB) Put(key, value []byte) error {
	return o.BucketPut(o.bucket, key, value)
}

func (o *BoltDB) Delete(key []byte) error {
	return o.BucketDelete(o.bucket, key)
}

func (o *BoltDB) BatchGet(offset []byte, size int) ([]*types.KV, error) {
	return o.BucketBatchGet(o.bucket, offset, size)
}

func (o *BoltDB) BatchPut(kv []*types.KV) error {
	return o.BucketBatchPut(o.bucket, kv)
}

func (o *BoltDB) BatchDelete(offset []byte, size int) error {
	return o.BucketBatchDelete(o.bucket, offset, size)
}

func (o *BoltDB) BucketGet(bucket, key []byte) ([]byte, error) {
	var value []byte
	err := o.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		if b == nil {
			return ErrBucketNotExist
		}
		v := b.Get(key)
		if len(v) == 0 {
			return ErrKeyNotExist
		}
		value = make([]byte, len(v))
		copy(value, v)
		return nil
	})

	return value, err
}

func (o *BoltDB) BucketPut(bucket, key, value []byte) error {
	err := o.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(bucket)
		if err != nil {
			return err
		}
		return b.Put(key, value)
	})

	return err
}

func (o *BoltDB) BucketDelete(bucket, key []byte) error {
	return o.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		if b == nil {
			return ErrBucketNotExist
		}
		return b.Delete(key)
	})
}

func (o *BoltDB) BucketBatchGet(bucket []byte, offset []byte, size int) ([]*types.KV, error) {
	var kvSlice []*types.KV
	err := o.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		if b == nil {
			return ErrBucketNotExist
		}
		i := int(0)
		c := b.Cursor()
		for ik, iv := c.Seek(offset); i < size && len(ik) != 0 && len(iv) != 0; ik, iv = c.Next() {
			key := make([]byte, len(ik))
			value := make([]byte, len(iv))
			copy(key, ik)
			copy(value, iv)
			kvSlice = append(kvSlice, &types.KV{Key: key, Value: value})
			i++
		}
		if i == 0 {
			return ErrNoRecord
		}
		return nil
	})

	return kvSlice, err
}

func (o *BoltDB) BucketBatchPut(bucket []byte, kv []*types.KV) error {
	return o.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(bucket)
		if err != nil {
			return err
		}
		for _, ele := range kv {
			err := b.Put(ele.Key, ele.Value)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (o *BoltDB) BucketBatchDelete(bucket, offset []byte, size int) error {
	return o.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		if b == nil {
			return ErrBucketNotExist
		}
		i := int(0)
		c := b.Cursor()
		key := make([][]byte, 0)
		for ik, iv := c.Seek(offset); i < size && len(ik) != 0 && len(iv) != 0; ik, iv = c.Next() {
			key = append(key, ik)
			i++
		}
		for _, k := range key {
			if err := b.Delete(k); err != nil {
				return err
			}
		}
		return nil
	})
}

func (o *BoltDB) Close() error {
	return o.DB.Close()
}
