/*
 *   Copyright (c) 2024 Arcology Network

 *   This program is free software: you can redistribute it and/or modify
 *   it under the terms of the GNU General Public License as published by
 *   the Free Software Foundation, either version 3 of the License, or
 *   (at your option) any later version.

 *   This program is distributed in the hope that it will be useful,
 *   but WITHOUT ANY WARRANTY; without even the implied warranty of
 *   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *   GNU General Public License for more details.

 *   You should have received a copy of the GNU General Public License
 *   along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

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
	tfdb         *TransactionalFileDB
	current      *Transaction
	previous     *Transaction
	optimization bool
}

func NewTransactionalStore() *TransactionalStore {
	return &TransactionalStore{}
}

func (ts *TransactionalStore) Config(params map[string]interface{}) {
	ts.tfdb = NewTransactionalFileDB(params["root"].(string))
	ts.optimization = params["optimization"].(bool)
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
	if ts.optimization {
		return nil
	} else {
		return ts.current.Add(request.Data, request.RecoverFunc)
	}
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
