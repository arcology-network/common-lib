package transactional

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"testing"
	"time"
)

func TestTransactionalFileDB(t *testing.T) {
	RegisterRecoverFunc("rf1", func(obj interface{}, bs []byte) error {
		var str string
		if obj != nil {
			str = obj.(string)
		} else {
			err := gob.NewDecoder(bytes.NewBuffer(bs)).Decode(&str)
			if err != nil {
				return err
			}
		}

		fmt.Printf("apply data: %v\n", str)
		return nil
	})
	RegisterRecoverFunc("rf2", func(obj interface{}, bs []byte) error {
		var array []byte
		if obj != nil {
			array = obj.([]byte)
		} else {
			err := gob.NewDecoder(bytes.NewBuffer(bs)).Decode(&array)
			if err != nil {
				return err
			}
		}

		fmt.Printf("apply data: %v\n", array)
		return nil
	})
	RegisterRecoverFunc("rf3", func(obj interface{}, bs []byte) error {
		var array []byte
		if obj != nil {
			array = obj.([]byte)
		} else {
			err := gob.NewDecoder(bytes.NewBuffer(bs)).Decode(&array)
			if err != nil {
				return err
			}
		}

		fmt.Printf("data len: %d\n", len(array))
		return nil
	})

	begin := time.Now()
	tfdb := NewTransactionalFileDB("./tfdb/")
	tx, err := tfdb.BeginTransaction("1")
	if err != nil {
		t.Error(err)
		return
	}

	if err = tx.Add("test string", "rf1"); err != nil {
		t.Error(err)
		return
	}
	if err = tx.Add([]byte("test byte array"), "rf2"); err != nil {
		t.Error(err)
		return
	}
	if err = tx.Add(make([]byte, 100000000), "rf3"); err != nil {
		t.Error(err)
		return
	}
	if err = tx.End(); err != nil {
		t.Error(err)
		return
	}

	if err = tfdb.Recover("1"); err != nil {
		t.Error(err)
	}
	t.Logf("elapsed time: %v\n", time.Since(begin))
}
