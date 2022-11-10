package types

import (
	"fmt"
	"reflect"
	"testing"

	ethCommon "github.com/arcology-network/3rd-party/eth/common"
)

func TestAssert(t *testing.T) {

	ret := []byte{8, 195, 121, 160, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 123, 11, 23, 33, 0, 0, 0}

	str := GetAssert(ret)

	fmt.Printf("str=%x\n", str)
}

func TestJson(t *testing.T) {
	logs := NewExecutingLogs()
	logs.Append("start", "s")
	logs.Append("doing", "s")
	logs.Append("end", "s")
	logs.Txhash = ethCommon.BytesToHash([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})

	str, err := logs.Marshal()
	if err != nil {
		t.Error("Marshal err !", err)
	}
	fmt.Printf("str=%v\n", str)

	logsn := &ExecutingLogs{}
	err = logsn.UnMarshal(str)
	if err != nil {
		t.Error("UnMarshal err !", err)
	}

	if !reflect.DeepEqual(logs, logsn) {
		t.Error("UnMarshal err !", logs, logsn)
	}

}
