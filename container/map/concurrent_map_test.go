package concurrentmap

import (
	"bytes"
	"fmt"
	"math/rand"
	"reflect"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/arcology-network/common-lib/codec"
)

func TestCcmapBasic(t *testing.T) {
	ccmap := NewConcurrentMap()
	ccmap.Set("1", 1)
	ccmap.Set("2", 2)
	ccmap.Set("3", 3)
	ccmap.Set("4", 4)

	if v, ok := ccmap.Get("1"); !ok || v.(int) != 1 {
		t.Error("Error: Failed to get")
	}
	if v, ok := ccmap.Get("2"); !ok || v.(int) != 2 {
		t.Error("Error: Failed to get")
	}
	if v, ok := ccmap.Get("3"); !ok || v.(int) != 3 {
		t.Error("Error: Failed to get")
	}
	if v, ok := ccmap.Get("4"); !ok || v.(int) != 4 {
		t.Error("Error: Failed to get")
	}

	ccmap.Set("1", 4)
	ccmap.Set("2", 3)
	ccmap.Set("3", 2)
	ccmap.Set("4", 1)

	if v, ok := ccmap.Get("1"); !ok || v.(int) != 4 {
		t.Error("Error: Failed to get")
	}
	if v, ok := ccmap.Get("2"); !ok || v.(int) != 3 {
		t.Error("Error: Failed to get")
	}
	if v, ok := ccmap.Get("3"); !ok || v.(int) != 2 {
		t.Error("Error: Failed to get")
	}
	if v, ok := ccmap.Get("4"); !ok || v.(int) != 1 {
		t.Error("Error: Failed to get")
	}

	ccmap.Set("1", "first")
	ccmap.Set("2", "second")
	ccmap.Set("3", 3)
	ccmap.Set("4", 4)

	if v, ok := ccmap.Get("1"); !ok || v.(string) != "first" {
		t.Error("Error: Failed to get")
	}

	if v, ok := ccmap.Get("2"); !ok || v.(string) != "second" {
		t.Error("Error: Failed to get")
	}

	if v, ok := ccmap.Get("3"); !ok || v.(int) != 3 {
		t.Error("Error: Failed to get")
	}

	if v, ok := ccmap.Get("4"); !ok || v.(int) != 4 {
		t.Error("Error: Failed to get")
	}

	if ok := ccmap.Set("5", nil); ok != nil {
		t.Error("Error: Failed to set")
	}

	keys := ccmap.Keys()
	sort.SliceStable(keys, func(i, j int) bool {
		return bytes.Compare([]byte(keys[i][:]), []byte(keys[j][:])) < 0
	})

	if v, ok := ccmap.Get("4"); !ok || v.(int) != 4 {
		t.Error("Error: Failed to get")
	}

	if !reflect.DeepEqual(keys, []string{"1", "2", "3", "4"}) {
		t.Error("Error: Entries don't match")
	}

	if ccmap.Size() != 4 {
		t.Error("Error: Wrong entry count")
	}
}

func TestCcmapEmptyKeys(t *testing.T) {
	ccmap := NewConcurrentMap()
	ccmap.Set("1", 1)
	ccmap.Set("2", 2)
	ccmap.Set("3", 3)
	ccmap.Set("", 4)

	if v, ok := ccmap.Get("1"); !ok || v.(int) != 1 {
		t.Error("Error: Failed to get")
	}

	if v, ok := ccmap.Get("2"); !ok || v.(int) != 2 {
		t.Error("Error: Failed to get")
	}

	if v, ok := ccmap.Get("3"); !ok || v.(int) != 3 {
		t.Error("Error: Failed to get")
	}

	if ccmap.Size() != 3 {
		t.Error("Error: Total count should be 3 ")
	}

	if v, _ := ccmap.Get(""); v != nil {
		t.Error("Error: Failed to get")
	}
}

func TestCcmapBatchModeAllEntries(t *testing.T) {
	ccmap := NewConcurrentMap()
	flags := []bool{true, true, true, true}
	keys := []string{"1", "2", "3", "4"}
	values := make([]interface{}, len(keys))
	for i, v := range keys {
		values[i] = v
	}

	ccmap.BatchSet(keys, values, flags)
	outValues := ccmap.BatchGet(keys)

	if !reflect.DeepEqual(outValues, values) {
		t.Error("Error: Entries don't match")
	}
}

func TestCcmapBatchModeSomeEntries(t *testing.T) {
	ccmap := NewConcurrentMap()
	flags := []bool{true, true, false, true}
	keys := []string{"1", "2", "3", "4"}
	values := make([]interface{}, len(keys))
	for i, v := range keys {
		values[i] = v
	}

	ccmap.BatchSet(keys, values, flags)
	outValues := ccmap.BatchGet(keys)

	if !reflect.DeepEqual(outValues, []interface{}{"1", "2", nil, "4"}) {
		t.Error("Error: Entries don't match")
	}
}

func TestCCmapDump(t *testing.T) {
	ccmap := NewConcurrentMap()
	flags := []bool{true, true, true, true}
	keys := []string{"1", "2", "3", "4"}
	values := []interface{}{"1", "2", 3, "4"}

	ccmap.BatchSet(keys, values, flags)
	k, v := ccmap.Dump()
	if !reflect.DeepEqual(k, []string{"1", "2", "3", "4"}) {
		t.Error("Error: Entries don't match")
	}

	if !reflect.DeepEqual(v, []interface{}{"1", "2", 3, "4"}) {
		t.Error("Error: Entries don't match")
	}
}

func TestMinMax(t *testing.T) {
	ccmap := NewConcurrentMap()
	keys := []string{"1", "2", "3", "4"}
	values := []interface{}{1, 2, 3, 4}
	ccmap.BatchSet(keys, values)

	less := func(lhs interface{}, rhs interface{}) bool {
		return lhs.(int) < rhs.(int)
	}

	min := ccmap.Find(less)
	if min != 1 {
		t.Error("Error: Wrong min value")
	}

	greater := func(lhs interface{}, rhs interface{}) bool {
		return lhs.(int) > rhs.(int)
	}

	max := ccmap.Find(greater)
	if max != 4 {
		t.Error("Error: Wrong max value")
	}
}
func BenchmarkMinMax(b *testing.B) {
	ccmap := NewConcurrentMap()
	keys := make([]string, 1000000)
	values := make([]interface{}, len(keys))
	for i := 0; i < len(keys); i++ {
		keys[i] = fmt.Sprint(i)
		values[i] = i
	}

	ccmap.BatchSet(keys, values)

	t0 := time.Now()
	less := func(lhs interface{}, rhs interface{}) bool {
		return lhs.(int) < rhs.(int)
	}

	min := ccmap.Find(less)
	if min != 0 {
		b.Error("Error: Wrong min value")
	}

	greater := func(lhs interface{}, rhs interface{}) bool {
		return lhs.(int) > rhs.(int)
	}

	max := ccmap.Find(greater)
	if max != 1000000-1 {
		b.Error("Error: Wrong max value")
	}
	fmt.Println("Min + Max ", time.Since(t0))
}

func TestForeach(t *testing.T) {
	ccmap := NewConcurrentMap()
	keys := []string{"1", "2", "3", "4"}
	values := []interface{}{1, 2, 3, 4}
	ccmap.BatchSet(keys, values)

	adder := func(v interface{}) interface{} {
		return v.(int) + 10
	}
	ccmap.Foreach(adder)

	_, vs := ccmap.Dump()
	if !reflect.DeepEqual(vs, []interface{}{1 + 10, 2 + 10, 3 + 10, 4 + 10}) {
		t.Error("Error: Checksums don't match")
	}
}

func BenchmarkForeach(b *testing.B) {
	ccmap := NewConcurrentMap()
	keys := make([]string, 1000000)
	values := make([]interface{}, len(keys))
	for i := 0; i < len(keys); i++ {
		keys[i] = fmt.Sprint(i)
		values[i] = i
	}
	ccmap.BatchSet(keys, values)

	t0 := time.Now()
	adder := func(v interface{}) interface{} {
		return v.(int) + 10
	}
	ccmap.Foreach(adder)
	fmt.Println("Foreach + 10 ", time.Since(t0))
}

func TestChecksum(t *testing.T) {
	ccmap := NewConcurrentMap()
	flags := []bool{true, true, true, true}
	keys := []string{"1", "2", "3", "4"}
	values := []interface{}{codec.String("1"), codec.String("2"), codec.Int64(3), codec.String("4")}

	ccmap.BatchSet(keys, values, flags)
	if !reflect.DeepEqual(ccmap.Checksum(), ccmap.Checksum()) {
		t.Error("Error: Checksums don't match")
	}
}

func BenchmarkCcmapBatchSet(b *testing.B) {
	genString := func() string {
		var letters = []rune("abcdefghijklmnopqrstuvwxyz0123456789")
		rand.Seed(time.Now().UnixNano())
		b := make([]rune, 40)
		for i := range b {
			b[i] = letters[rand.Intn(len(letters))]
		}
		return string(b)
	}

	values := make([]interface{}, 0, 100000)
	paths := make([]string, 0, 100000)
	for i := 0; i < 100000; i++ {
		acct := genString()
		paths = append(paths, []string{
			"blcc://eth1.0/account/" + acct + "/",
			"blcc://eth1.0/account/" + acct + "/code",
			"blcc://eth1.0/account/" + acct + "/nonce",
			"blcc://eth1.0/account/" + acct + "/balance",
			"blcc://eth1.0/account/" + acct + "/defer/",
			"blcc://eth1.0/account/" + acct + "/storage/",
			"blcc://eth1.0/account/" + acct + "/storage/containers/",
			"blcc://eth1.0/account/" + acct + "/storage/native/",
			"blcc://eth1.0/account/" + acct + "/storage/containers/!/",
			"blcc://eth1.0/account/" + acct + "/storage/containers/KittyIndexToOwner/$ad90f8111111111111111111111111111111111111",
			"blcc://eth1.0/account/" + acct + "/storage/containers/KittyIndexToOwner/$ad90f8211111111111111111111111111111111111",
			"blcc://eth1.0/account/" + acct + "/storage/containers/KittyIndexToOwner/$ad90f8311111111111111111111111111111111111",
			"blcc://eth1.0/account/" + acct + "/storage/containers/KittyIndexToOwner/$ad90f8411111111111111111111111111111111111",
		}...)
	}

	for i := 0; i < len(paths); i++ {
		values = append(values, paths[i])
	}

	t0 := time.Now()
	var masterLock sync.RWMutex
	for i := 0; i < 1000000; i++ {
		masterLock.Lock()
		masterLock.Unlock()
	}
	fmt.Println("Lock() 1000000 "+fmt.Sprint(len(paths)), " in ", time.Since(t0))

	t0 = time.Now()
	ccmap := NewConcurrentMap()
	ccmap.BatchSet(paths, values)
	fmt.Println("BatchSet "+fmt.Sprint(len(paths)), " in ", time.Since(t0))
}
