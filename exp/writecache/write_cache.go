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

/*
 *   Copyright (c) 2023 Arcology Network
 *   All rights reserved.

 *   Licensed under the Apache License, Version 2.0 (the "License");
 *   you may not use this file except in compliance with the License.
 *   You may obtain a copy of the License at

 *   http://www.apache.org/licenses/LICENSE-2.0

 *   Unless required by applicable law or agreed to in writing, software
 *   distributed under the License is distributed on an "AS IS" BASIS,
 *   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *   See the License for the specific language governing permissions and
 *   limitations under the License.
 */

package cache

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"

	common "github.com/arcology-network/common-lib/common"
	mapi "github.com/arcology-network/common-lib/exp/map"
	mempool "github.com/arcology-network/common-lib/exp/mempool"
	slice "github.com/arcology-network/common-lib/exp/slice"
	stgcomm "github.com/arcology-network/storage-committer"
	committercommon "github.com/arcology-network/storage-committer/common"
	platform "github.com/arcology-network/storage-committer/platform"

	"github.com/arcology-network/storage-committer/commutative"
	importer "github.com/arcology-network/storage-committer/importer"
	"github.com/arcology-network/storage-committer/interfaces"
	intf "github.com/arcology-network/storage-committer/interfaces"
	"github.com/arcology-network/storage-committer/noncommutative"
	univalue "github.com/arcology-network/storage-committer/univalue"
)

// WriteCache is a read-only data store used for caching.
type WriteCache[T0 comparable, T1 any] struct {
	store    intf.ReadOnlyDataStore
	kvDict   map[T0]T1 // Local KV lookup
	platform intf.Platform
	buffer   []T1 // Transition + access record buffer
	uniPool  *mempool.Mempool[*univalue.Univalue]
}

// NewWriteCache creates a new instance of WriteCache; the store can be another instance of WriteCache,
// resulting in a cascading-like structure.
func NewWriteCache[T0 comparable, T1 any](store intf.ReadOnlyDataStore, perPage int, numPages int, args ...interface{}) *WriteCache[T0, T1] {
	// t0 := time.Now()
	var writeCache WriteCache[T0, T1]
	writeCache.store = store
	writeCache.kvDict = make(map[T0]*univalue.Univalue)
	writeCache.platform = platform.NewPlatform()
	writeCache.buffer = make([]*univalue.Univalue, 0, perPage*numPages)

	writeCache.uniPool = mempool.NewMempool[*univalue.Univalue](perPage, numPages, func() T1 {
		return new(univalue.Univalue)
	}, (&univalue.Univalue{}).Reset)
	// fmt.Println("NewWriteCache ------------- ", time.Since(t0))
	return &writeCache
}

// CreateNewAccount creates a new account in the write cache.
// It returns the transitions and an error, if any.
func (this *WriteCache[T0, T1]) CreateNewAccount(tx uint32, acct string) ([]*univalue.Univalue, error) {
	paths, typeids := platform.NewPlatform().GetBuiltins(acct)

	transitions := []*univalue.Univalue{}
	for i, path := range paths {
		var v interface{}
		switch typeids[i] {
		case commutative.PATH: // Path
			v = commutative.NewPath()

		case uint8(reflect.Kind(noncommutative.STRING)): // delta big int
			v = noncommutative.NewString("")

		case uint8(reflect.Kind(commutative.UINT256)): // delta big int
			v = commutative.NewUnboundedU256()

		case uint8(reflect.Kind(commutative.UINT64)):
			v = commutative.NewUnboundedUint64()

		case uint8(reflect.Kind(noncommutative.INT64)):
			v = new(noncommutative.Int64)

		case uint8(reflect.Kind(noncommutative.BYTES)):
			v = noncommutative.NewBytes([]byte{})
		}

		// fmt.Println(path)
		if !this.IfExists(path) {
			transitions = append(transitions, univalue.NewUnivalue(tx, path, 0, 1, 0, v, nil))

			if _, err := this.Write(tx, path, v); err != nil { // root path
				return nil, err
			}

			if !this.IfExists(path) {
				_, err := this.Write(tx, path, v)
				return transitions, err // root path
			}
		}
	}
	return transitions, nil
}

func (this *WriteCache[T0, T1]) SetReadOnlyDataStore(store intf.ReadOnlyDataStore) *WriteCache[T0, T1] {
	this.store = store
	return this
}

func (this *WriteCache[T0, T1]) ReadOnlyDataStore() intf.ReadOnlyDataStore { return this.store }
func (this *WriteCache[T0, T1]) Cache() map[T0]T1                          { return this.kvDict }
func (this *WriteCache[T0, T1]) MinSize() int                              { return this.uniPool.MinSize() }
func (this *WriteCache[T0, T1]) NewUnivalue() T1                           { return this.uniPool.New() }

// If the access has been recorded
func (this *WriteCache[T0, T1]) GetOrNew(tx uint32, key T0, Type any) (*univalue.Univalue, bool) {
	unival, inCache := this.kvDict[key]
	if unival == nil { // Not in the kvDict, check the datastore
		var typedv interface{}
		if store := this.ReadOnlyDataStore(); store != nil {
			typedv = common.FilterFirst(store.Retrive(key, T))
		}

		unival = this.NewUnivalue().Init(tx, key, 0, 0, 0, typedv, this)
		this.kvDict[key] = unival // Adding to kvDict
	}
	return unival, inCache // From cache
}

func (this *WriteCache[T0, T1]) Read(tx uint32, key T0, T any) (interface{}, interface{}, uint64) {
	univalue, _ := this.GetOrNew(tx, key, T)
	return univalue.Get(tx, key, nil), univalue, 0
}

func (this *WriteCache[T0, T1]) write(tx uint32, key T0, value interface{}) error {
	parentPath := common.GetParentPath(key)
	if this.IfExists(parentPath) || tx == committercommon.SYSTEM { // The parent path exists or to inject the path directly
		univalue, inCache := this.GetOrNew(tx, key, value) // Get a univalue wrapper
		err := univalue.Set(tx, key, value, inCache, this)

		if err == nil {
			if strings.HasSuffix(parentPath, "/container/") || (!this.platform.IsSysPath(parentPath) && tx != committercommon.SYSTEM) { // Don't keep track of the system children
				parentMeta, inCache := this.GetOrNew(tx, parentPath, new(commutative.Path))
				err = parentMeta.Set(tx, path, univalue.Value(), inCache, this)
			}
		}
		return err
	}
	return errors.New("Error: The parent path doesn't exist: " + parentPath)
}

func (this *WriteCache[T0, T1]) Write(tx uint32, key T0, value interface{}) (int64, error) {
	fee := int64(0) //Fee{}.Writer(key, value, this.writeCache)
	if value == nil || (value != nil && value.(interfaces.Type).TypeID() != uint8(reflect.Invalid)) {
		return fee, this.write(tx, key, value)
	}
	return fee, errors.New("Error: Unknown data type !")
}

// Get data from the DB direcly, still under conflict protection
func (this *WriteCache[T0, T1]) ReadCommitted(tx uint32, key T0, T any) (interface{}, uint64) {
	if v, _, Fee := this.Read(tx, key, this); v != nil { // For conflict detection
		return v, Fee
	}

	v, _ := this.ReadOnlyDataStore().Retrive(key, T)
	if v == nil {
		return v, 0 //Fee{}.Reader(univalue.NewUnivalue(tx, key, 1, 0, 0, v, nil))
	}
	return v, 0 //Fee{}.Reader(univalue.NewUnivalue(tx, key, 1, 0, 0, v.(interfaces.Type), nil))
}

// Get the raw value directly, skip the access counting at the univalue level
func (this *WriteCache[T0, T1]) InCache(key T0) (interface{}, bool) {
	univ, ok := this.kvDict[key]
	return univ, ok
}

// Get the raw value directly, skip the access counting at the univalue level
func (this *WriteCache[T0, T1]) Find(key T0, T any) (interface{}, interface{}) {
	if univ, ok := this.kvDict[key]; ok {
		return univ.Value(), univ
	}

	v, _ := this.ReadOnlyDataStore().Retrive(key, T)
	univ := univalue.NewUnivalue(committercommon.SYSTEM, key, 0, 0, 0, v, nil)
	return univ.Value(), univ
}

func (this *WriteCache[T0, T1]) Retrive(key T0, T any) (interface{}, error) {
	typedv, _ := this.Find(key, T)
	if typedv == nil || typedv.(intf.Type).IsDeltaApplied() {
		return typedv, nil
	}

	rawv, _, _ := typedv.(intf.Type).Get()
	return typedv.(intf.Type).New(rawv, nil, nil, typedv.(intf.Type).Min(), typedv.(intf.Type).Max()), nil // Return in a new univalue
}

func (this *WriteCache[T0, T1]) IfExists(key T0) bool {
	if committercommon.ETH10_ACCOUNT_PREFIX_LENGTH == len(key) {
		return true
	}

	if v := this.kvDict[key]; v != nil {
		return v.Value() != nil // If value == nil means either it's been deleted or never existed.
	}

	if this.store == nil {
		return false
	}
	return this.store.IfExists(key) //this.RetriveShallow(path, nil) != nil
}

// The function is used to add the transitions to the writecache, which usually comes from
// the child writecaches. It usually happens with the sub processeses are completed.
func (this *WriteCache[T0, T1]) AddTransitions(transitions []*univalue.Univalue) {
	if len(transitions) == 0 {
		return
	}

	// Filter out the key creations transitions as they will be treated differently.
	newPathCreations := slice.MoveIf(&transitions, func(_ int, v *univalue.Univalue) bool {
		return common.IsPath(*v.GetPath()) && !v.Preexist()
	})

	// Not necessary to sort the path creations at the moment,
	// but it is good for the future if multiple level containers are available
	newPathCreations = univalue.Univalues(importer.Sorter(newPathCreations))
	slice.Foreach(newPathCreations, func(_ int, v **univalue.Univalue) {
		(*v).CopyTo(this) // Write back to the parent writecache
	})

	// Remove the changes to the existing path meta, as they will be updated automatically when inserting sub elements.
	transitions = slice.RemoveIf(&transitions, func(_ int, v *univalue.Univalue) bool {
		return common.IsPath(*v.GetPath())
	})

	// Write back to the parent writecache
	slice.Foreach(transitions, func(_ int, v **univalue.Univalue) {
		(*v).CopyTo(this)
	})
}

// Reset the writecache to the initial state for the next round of processing.
func (*WriteCache[T0, T1]) Reset(this *WriteCache[T0, T1]) {
	if clear(this.buffer); cap(this.buffer) > 3*this.uniPool.MinSize() {
		this.buffer = make([]*univalue.Univalue, 0, this.uniPool.MinSize())
	}
	this.uniPool.Reset()
	clear(this.kvDict)
}

func (this *WriteCache[T0, T1]) Equal(other *WriteCache[T0, T1]) bool {
	thisBuffer := mapi.Values(this.kvDict)
	sort.SliceStable(thisBuffer, func(i, j int) bool {
		return *thisBuffer[i].GetPath() < *thisBuffer[j].GetPath()
	})

	otherBuffer := mapi.Values(other.kvDict)
	sort.SliceStable(otherBuffer, func(i, j int) bool {
		return *otherBuffer[i].GetPath() < *otherBuffer[j].GetPath()
	})

	cacheFlag := reflect.DeepEqual(thisBuffer, otherBuffer)
	return cacheFlag
}

func (this *WriteCache[T0, T1]) Export(preprocessors ...func([]*univalue.Univalue) []*univalue.Univalue) []T1 {
	this.buffer = mapi.Values(this.kvDict) //this.buffer[:0]

	for _, processor := range preprocessors {
		this.buffer = common.IfThenDo1st(processor != nil, func() []T1 {
			return processor(this.buffer)
		}, this.buffer)
	}

	slice.RemoveIf(&this.buffer, func(_ int, v *univalue.Univalue) bool { return v.Reads() == 0 && v.IsReadOnly() }) // Remove peeks
	return this.buffer
}

func (this *WriteCache[T0, T1]) ExportAll(preprocessors ...func([]*univalue.Univalue) []*univalue.Univalue) ([]*univalue.Univalue, []*univalue.Univalue) {
	all := this.Export(importer.Sorter)
	// univalue.Univalues(all).Print()

	accesses := univalue.Univalues(slice.Clone(all)).To(importer.ITAccess{})
	transitions := univalue.Univalues(slice.Clone(all)).To(importer.ITTransition{})
	return accesses, transitions
}

func (this *WriteCache[T0, T1]) Print() {
	values := mapi.Values(this.kvDict)
	sort.SliceStable(values, func(i, j int) bool {
		return *values[i].GetPath() < *values[j].GetPath()
	})

	for i, elem := range values {
		fmt.Println("Level : ", i)
		elem.Print()
	}
}

func (this *WriteCache[T0, T1]) KVs() ([]string, []intf.Type) {
	transitions := univalue.Univalues(slice.Clone(this.Export(importer.Sorter))).To(importer.ITTransition{})

	values := make([]intf.Type, len(transitions))
	keys := slice.ParallelAppend(transitions, 4, func(i int, v *univalue.Univalue) string {
		values[i] = v.Value().(intf.Type)
		return *v.GetPath()
	})
	return keys, values
}

// This function is used to write the cache to the data source directly to bypass all the intermediate steps,
// including the conflict detection.
//
// It's mainly used for TESTING purpose.
func (this *WriteCache[T0, T1]) FlushToDataSource(store interfaces.Datastore) interfaces.Datastore {
	committer := stgcomm.NewStorageCommitter(store)
	acctTrans := univalue.Univalues(slice.Clone(this.Export(importer.Sorter))).To(importer.IPTransition{})

	txs := slice.Transform(acctTrans, func(_ int, v *univalue.Univalue) uint32 {
		return v.GetTx()
	})

	committer.Import(acctTrans)
	committer.Sort()
	committer.Precommit(txs)
	committer.Commit()
	this.Reset(this)

	return store
}
