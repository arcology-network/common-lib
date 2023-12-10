package types

import (
	"encoding/json"
	"fmt"

	ethCommon "github.com/ethereum/go-ethereum/common"
)

type ExecutingLog struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (this *ExecutingLog) GetKey() string {
	return this.Key
}

func (this *ExecutingLog) GetValue() string {
	return this.Value
}

type ExecutingLogs struct {
	Txhash ethCommon.Hash `json:"txhash"`
	Logs   []ExecutingLog `json:"this"`
}

type ExecutingLogsMessage struct {
	Logs   ExecutingLogs
	Height uint64
	Round  uint64
	Msgid  uint64
}

func GetAssert(ret []byte) string {
	startIdx := 4 + 32 + 32
	pattern := []byte{8, 195, 121, 160}
	if ret != nil || len(ret) > startIdx {
		starts := ret[:4]
		if string(pattern) == string(starts) {
			remains := ret[startIdx:]
			end := 0
			for i := range remains {
				if remains[i] == 0 {
					end = i
					break
				}
			}
			return string(remains[:end])
		}
	}
	return ""
}

func NewExecutingLogs() *ExecutingLogs {
	return &ExecutingLogs{
		Logs: []ExecutingLog{},
	}
}

func (this *ExecutingLogs) Append(key, value string) {
	this.Logs = append(this.Logs, ExecutingLog{
		Key:   key,
		Value: value,
	})
}
func (this *ExecutingLogs) Appends(log []ExecutingLog) {
	this.Logs = append(this.Logs, log...)
}

func (this *ExecutingLogs) Marshal() (string, error) {
	data, err := json.Marshal(this)
	return fmt.Sprintf("%v", string(data)), err
}

func (this *ExecutingLogs) UnMarshal(data string) error {

	return json.Unmarshal([]byte(data), this)
}
