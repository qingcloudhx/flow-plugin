//Created by zhbinary on 2019-04-12.
//Email: zhbinary@gmail.com
package types

type KV struct {
	Key   []byte
	Value []byte
}

type Persist interface {
	Get(key []byte) ([]byte, error)
	Put(key, value []byte) error
	Delete(key []byte) error

	BatchGet(offset []byte, size int) ([]*KV, error)
	BatchPut(kv []*KV) error
	BatchDelete(offset []byte, size int) error

	BucketGet(bucket, key []byte) ([]byte, error)
	BucketPut(bucket, key, value []byte) error
	BucketDelete(bucket, key []byte) error

	BucketBatchGet(bucket []byte, offset []byte, size int) ([]*KV, error)
	BucketBatchPut(bucket []byte, kv []*KV) error
	BucketBatchDelete(bucket, offset []byte, size int) error

	Close() error
}
