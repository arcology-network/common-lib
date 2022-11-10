package transactional

import (
	"context"
	"fmt"
)

type AddDataRequest struct {
	Data        interface{}
	RecoverFunc string
}

type TransactionalStore struct {
	tfdb     *TransactionalFileDB
	current  *Transaction
	previous *Transaction
}

func NewTransactionalStore() *TransactionalStore {
	return &TransactionalStore{}
}

func (ts *TransactionalStore) Config(params map[string]interface{}) {
	ts.tfdb = NewTransactionalFileDB(params["root"].(string))
}

func (ts *TransactionalStore) BeginTransaction(ctx context.Context, id *string, _ *int) (err error) {
	if ts.current != nil {
		panic("BeginTransaction called in another transaction.")
	}
	ts.current, err = ts.tfdb.BeginTransaction(*id)
	return
}

func (ts *TransactionalStore) AddData(ctx context.Context, request *AddDataRequest, _ *int) error {
	if ts.current == nil {
		panic("AddData called before BeginTransaction.")
	}
	//return ts.current.Add(request.Data, request.RecoverFunc)
	return nil
}

func (ts *TransactionalStore) EndTransaction(ctx context.Context, _ *int, _ *int) error {
	if ts.current == nil {
		panic("EndTransaction called before BeginTransaction.")
	}

	if ts.previous != nil {
		ts.previous.Clear()
	}

	err := ts.current.End()
	ts.previous = ts.current
	ts.current = nil
	return err
}

func (ts *TransactionalStore) Recover(ctx context.Context, id *string, _ *int) error {
	fmt.Printf("[TransactionalStore.Recover] id = %s\n", *id)
	return ts.tfdb.Recover(*id)
}
