package persist

import (
	"github.com/qingcloudhx/flow-plugin/trigger/mqtt-broker/broker/types"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"testing"
)

func randInt(min int, max int) int {
	return rand.Intn(max-min) + min
}

func TestBoltDB(t *testing.T) {
	key := []byte("testKey")
	value := []byte("testValue")
	dbPath := filepath.Join(os.TempDir(), "bolt_test.db")
	db, err := NewBoltDB(dbPath)
	defer func() {
		err := db.Close()
		if err != nil {
			t.Fatal(err)
		}
		_ = os.RemoveAll(dbPath)
	}()
	if err != nil {
		t.Fatal(err)
	}

	v, err := db.Get(key)
	if v != nil || err == nil {
		t.Fatal(err)
	}
	err = db.Put(key, value)
	if err != nil {
		t.Fatal(err)
	}
	v, err = db.Get(key)
	if v == nil || err != nil || !reflect.DeepEqual(v, value) {
		t.Fatal(err)
	}
	err = db.Delete(key)
	if err != nil {
		t.Fatal(err)
	}
	v, err = db.Get(key)
	if v != nil || err == nil {
		t.Fatal(err)
	}
}

func TestBoltDBBucket(t *testing.T) {
	key := []byte("testKey")
	value := []byte("testValue")
	bucket := []byte("testBucket")
	dbPath := filepath.Join(os.TempDir(), "bolt_test.db")
	db, err := NewBoltDB(dbPath)
	defer func() {
		_ = db.Close()
		_ = os.RemoveAll(dbPath)
	}()
	if err != nil {
		t.Fatal(err)
	}

	v, err := db.BucketGet(bucket, key)
	if v != nil || err == nil {
		t.Fatal(err)
	}
	err = db.BucketPut(bucket, key, value)
	if err != nil {
		t.Fatal(err)
	}
	v, err = db.BucketGet(bucket, key)
	if v == nil || err != nil || !reflect.DeepEqual(v, value) {
		t.Fatal(err)
	}
	v, err = db.Get(key)
	if v != nil || err == nil {
		t.Fatal(err)
	}
	err = db.BucketDelete(bucket, key)
	if err != nil {
		t.Fatal(err)
	}
	v, err = db.BucketGet(bucket, key)
	if v != nil || err == nil {
		t.Fatal(err)
	}
}

func TestBoltDBBatch(t *testing.T) {
	cnt := 10
	kvSlice := make([]*types.KV, cnt)
	for i := 0; i < cnt; i++ {
		kvSlice[i] = &types.KV{Key: []byte(strconv.Itoa(i)), Value: []byte(strconv.Itoa(i))}
	}
	dbPath := filepath.Join(os.TempDir(), "bolt_test.db")
	db, err := NewBoltDB(dbPath)
	defer func() {
		err := db.Close()
		if err != nil {
			t.Fatal(err)
		}
		_ = os.RemoveAll(dbPath)
	}()
	if err != nil {
		t.Fatal(err)
	}

	v, err := db.BatchGet([]byte(strconv.Itoa(0)), cnt)
	if len(v) != 0 || err == nil {
		t.Fatal(err)
	}
	err = db.BatchPut(kvSlice)
	if err != nil {
		t.Fatal(err)
	}
	start := randInt(0, cnt/2)
	num := randInt(0, cnt-start+1)
	v, err = db.BatchGet([]byte(strconv.Itoa(start)), num)
	if v == nil || err != nil || len(v) != num {
		t.Fatal(err)
	}
	for i := start; i < start+num; i++ {
		if string(kvSlice[i].Key[:]) != string(v[i-start].Key[:]) || string(kvSlice[i].Value[:]) != string(v[i-start].Value[:]) {
			t.Fatal(string(kvSlice[i].Key[:]), string(v[i-start].Key[:]))
		}
	}

	err = db.BatchDelete([]byte(strconv.Itoa(start)), num)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < start; i++ {
		v, err := db.Get([]byte(strconv.Itoa(i)))
		if err != nil {
			t.Fatal(err)
		}
		if string(kvSlice[i].Value[:]) != string(v[:]) {
			t.Fatal(string(kvSlice[i].Value[:]), string(v[:]))
		}
	}
	for i := start; i < start+num; i++ {
		v, err := db.Get([]byte(strconv.Itoa(i)))
		if len(v) != 0 || err == nil {
			t.Fatal(err)
		}
	}
	for i := start + num; i < cnt; i++ {
		v, err := db.Get([]byte(strconv.Itoa(i)))
		if err != nil {
			t.Fatal(err)
		}
		if string(kvSlice[i].Value[:]) != string(v[:]) {
			t.Fatal(string(kvSlice[i].Value[:]), string(v[:]))
		}
	}
}

func TestBoltDBBucketBatch(t *testing.T) {
	cnt := 10
	kvSlice := make([]*types.KV, cnt)
	bucket := []byte("testBucket")
	for i := 0; i < cnt; i++ {
		kvSlice[i] = &types.KV{Key: []byte(strconv.Itoa(i)), Value: []byte(strconv.Itoa(i))}
	}
	dbPath := filepath.Join(os.TempDir(), "bolt_test.db")
	db, err := NewBoltDB(dbPath)
	defer func() {
		err := db.Close()
		if err != nil {
			t.Fatal(err)
		}
		_ = os.RemoveAll(dbPath)
	}()
	if err != nil {
		t.Fatal(err)
	}

	v, err := db.BucketBatchGet(bucket, []byte(strconv.Itoa(0)), cnt)
	if len(v) != 0 || err == nil {
		t.Fatal(err)
	}
	err = db.BucketBatchPut(bucket, kvSlice)
	if err != nil {
		t.Fatal(err)
	}
	start := randInt(0, cnt/2)
	num := randInt(0, cnt-start+1)
	v, err = db.BucketBatchGet(bucket, []byte(strconv.Itoa(start)), num)
	if v == nil || err != nil || len(v) != num {
		t.Fatal(err)
	}
	for i := start; i < start+num; i++ {
		if string(kvSlice[i].Key[:]) != string(v[i-start].Key[:]) || string(kvSlice[i].Value[:]) != string(v[i-start].Value[:]) {
			t.Fatal(string(kvSlice[i].Key[:]), string(v[i-start].Key[:]))
		}
	}

	err = db.BucketBatchDelete(bucket, []byte(strconv.Itoa(start)), num)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < start; i++ {
		v, err := db.BucketGet(bucket, []byte(strconv.Itoa(i)))
		if err != nil {
			t.Fatal(err)
		}
		if string(kvSlice[i].Value[:]) != string(v[:]) {
			t.Fatal(string(kvSlice[i].Value[:]), string(v[:]))
		}
	}
	for i := start; i < start+num; i++ {
		v, err := db.BucketGet(bucket, []byte(strconv.Itoa(i)))
		if len(v) != 0 || err == nil {
			t.Fatal(err)
		}
	}
	for i := start + num; i < cnt; i++ {
		v, err := db.BucketGet(bucket, []byte(strconv.Itoa(i)))
		if err != nil {
			t.Fatal(err)
		}
		if string(kvSlice[i].Value[:]) != string(v[:]) {
			t.Fatal(string(kvSlice[i].Value[:]), string(v[:]))
		}
	}
}
