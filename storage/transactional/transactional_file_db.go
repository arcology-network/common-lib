package transactional

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"sync"
)

type RecoverFunc func(obj interface{}, bs []byte) error

var RecoverFuncRegistry = make(map[string]RecoverFunc)

func RegisterRecoverFunc(name string, rf RecoverFunc) {
	RecoverFuncRegistry[name] = rf
}

type TransactionalFileDB struct {
	root string
	db   *SimpleFileDB
}

func NewTransactionalFileDB(root string) *TransactionalFileDB {
	return &TransactionalFileDB{
		root: root,
		db:   NewSimpleFileDB(root),
	}
}

func (tfdb *TransactionalFileDB) BeginTransaction(id string) (*Transaction, error) {
	return NewTransaction(id, tfdb.db)
}

func (tfdb *TransactionalFileDB) Recover(id string) error {
	fmt.Printf("[TransactionalFileDB.Recover] id = %s\n", id)
	bs, err := tfdb.db.Get(id)
	if err != nil {
		fmt.Printf("[TransactionalFileDB.Recover] transaction not found, err = %v\n", err)
		return nil
	}

	var rfs map[string]string
	err = gob.NewDecoder(bytes.NewBuffer(bs)).Decode(&rfs)
	if err != nil {
		fmt.Printf("[TransactionalFileDB.Recover] Decode transaction file err: %v\n", err)
		return err
	}
	tx := Transaction{
		id:  id,
		db:  tfdb.db,
		rfs: rfs,
	}
	return tx.commit()
}

type Transaction struct {
	id   string
	db   *SimpleFileDB
	rfs  map[string]string
	buf  map[string]interface{}
	wg   sync.WaitGroup
	lock sync.Mutex
}

func NewTransaction(id string, db *SimpleFileDB) (*Transaction, error) {
	if _, err := db.Get(id); err != nil {
		return &Transaction{
			id:  id,
			db:  db,
			rfs: make(map[string]string),
			buf: make(map[string]interface{}),
		}, nil
	}
	return nil, fmt.Errorf("Transaction already exists: %v", id)
}

func (t *Transaction) Add(obj interface{}, rf string) error {
	if _, ok := RecoverFuncRegistry[rf]; !ok {
		return fmt.Errorf("Recover function not found: %v", rf)
	}

	t.wg.Add(1)

	go func() {
		var buf bytes.Buffer
		err := gob.NewEncoder(&buf).Encode(obj)
		if err != nil {
			return
		}

		value := buf.Bytes()
		key := fmt.Sprintf("%x", sha256.Sum256(value))
		err = t.db.Set(key, buf.Bytes())
		if err != nil {
			return
		}

		t.lock.Lock()
		t.rfs[key] = rf
		t.buf[key] = obj
		t.lock.Unlock()

		t.wg.Done()
	}()
	return nil
}

func (t *Transaction) End() error {
	t.wg.Wait()

	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(t.rfs)
	if err != nil {
		return err
	}

	return t.db.Set(t.id, buf.Bytes())
}

func (t *Transaction) commit() error {
	fmt.Printf("[Transaction.commit] rfs = %v\n", t.rfs)
	for key, rf := range t.rfs {
		if len(t.buf) != 0 {
			RecoverFuncRegistry[rf](t.buf[key], nil)
			continue
		}

		if bs, err := t.db.Get(key); err != nil {
			fmt.Printf("[Transaction.commit] transaction file not found: %v, err: %v\n", key, err)
			return err
		} else {
			if err = RecoverFuncRegistry[rf](nil, bs); err != nil {
				return err
			}
		}
	}

	return t.Clear()
}

func (t *Transaction) Clear() error {
	t.db.Delete(t.id)
	for key := range t.rfs {
		t.db.Delete(key)
	}
	return nil
}
