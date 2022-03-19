package mhasher

import (
	"fmt"
	"reflect"
	"testing"
)

func TestStringEngine(t *testing.T) {
	paths := make([]string, 4)
	paths[2] = "000000000000000000"
	paths[1] = "1111111111111111111110"
	paths[0] = "2222222222222222"
	paths[3] = "2222222222222222"

	se := Start()
	err := se.ToBuffer(paths)
	if err != nil {
		fmt.Printf("ToBuffer err: %v\n", err)
		return
	}

	retpaths, err := se.FromBuffer(paths)
	if err != nil {
		fmt.Printf("FromBuffer err: %v\n", err)
		return
	}

	for i := range retpaths {
		fmt.Printf("paths=%x\n", retpaths[i])
	}

	se.Clear()
	se.Stop()
}

func TestSortString(t *testing.T) {
	data := []string{"5678", "9101112", "1234"}

	result, err := SortStrings(data)
	if err != nil {
		t.Error("sort err:" + err.Error())
	}

	if !reflect.DeepEqual(result[0], data[2]) || !reflect.DeepEqual(result[1], data[0]) || !reflect.DeepEqual(result[2], data[1]) {
		t.Error("Wrong order !")
	}
	fmt.Printf("result=%v\n", result)
}
func TestUniqueSortStrings(t *testing.T) {
	data := []string{"5678", "9101112", "5678", "1234"}

	result, err := UniqueSortStrings(data)
	if err != nil {
		t.Error("sort err:" + err.Error())
	}

	if !reflect.DeepEqual(result[0], data[3]) || !reflect.DeepEqual(result[1], data[0]) || !reflect.DeepEqual(result[2], data[1]) {
		t.Error("Wrong order !")
	}
	fmt.Printf("result=%v\n", result)
}

func TestUniqueStrings(t *testing.T) {
	data := []string{"124", "5678", "1258", "5678"}

	result, err := UniqueStrings(data)
	if err != nil {
		t.Error("Unique err:" + err.Error())
	}
	if !reflect.DeepEqual(len(result), 3) || !reflect.DeepEqual(result[0], data[0]) || !reflect.DeepEqual(result[1], data[1]) || !reflect.DeepEqual(result[2], data[2]) {
		t.Error("Wrong Unique !")
	}
	fmt.Printf("result=%v\n", result)
}

func TestRemoveString(t *testing.T) {
	data := []string{"124", "5678", "9012"}
	toRemove := []string{"124", "5678"}
	result, err := RemoveString(data, toRemove)
	if err != nil {
		t.Error("Unique err:" + err.Error())
	}
	fmt.Printf("result=%v\n", result)
	if !reflect.DeepEqual(len(result), 1) || !reflect.DeepEqual(result[0], data[2]) {
		t.Error("Wrong Remove !")
	}

}
func TestSortBytes(t *testing.T) {
	data := [][]byte{[]byte{1, 2, 4}, []byte{5, 6, 7, 8}, []byte{9, 10, 12}}

	result, err := SortBytes(data)
	if err != nil {
		t.Error("sort err:" + err.Error())
	}
	if !reflect.DeepEqual(result[0], data[0]) || !reflect.DeepEqual(result[1], data[1]) || !reflect.DeepEqual(result[2], data[2]) {
		t.Error("Wrong order !")
	}
	fmt.Printf("result=%v\n", result)
}

func TestUniqueBytes(t *testing.T) {
	data := [][]byte{[]byte{1, 2, 4}, []byte{5, 6, 7, 8}, []byte{9, 10, 12}, []byte{5, 6, 7, 8}}

	result, err := UniqueBytes(data)
	if err != nil {
		t.Error("Unique err:" + err.Error())
	}
	if !reflect.DeepEqual(len(result), 3) || !reflect.DeepEqual(result[0], data[0]) || !reflect.DeepEqual(result[1], data[1]) || !reflect.DeepEqual(result[2], data[2]) {
		t.Error("Wrong Unique !")
	}
	fmt.Printf("result=%v\n", result)
}

func TestRemoveBytes(t *testing.T) {
	data := [][]byte{[]byte{1, 2, 4}, []byte{5, 6, 7, 8}, []byte{9, 10, 12}}
	toRemove := [][]byte{[]byte{1, 2, 4}, []byte{5, 6, 7, 8}}
	result, err := RemoveBytes(data, toRemove)
	if err != nil {
		t.Error("Unique err:" + err.Error())
	}
	fmt.Printf("result=%v\n", result)
	if !reflect.DeepEqual(len(result), 1) || !reflect.DeepEqual(result[0], data[2]) {
		t.Error("Wrong Remove !")
	}

}
