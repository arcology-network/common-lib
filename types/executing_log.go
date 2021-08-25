package types

import (
	"encoding/json"
	"fmt"

	ethCommon "github.com/HPISTechnologies/3rd-party/eth/common"
)

type ExecutingLog struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (log *ExecutingLog) GetKey() string {
	return log.Key
}

func (log *ExecutingLog) GetValue() string {
	return log.Value
}

type ExecutingLogs struct {
	Txhash ethCommon.Hash `json:"txhash"`
	Logs   []ExecutingLog `json:"logs"`
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

func (logs *ExecutingLogs) Append(key, value string) {
	logs.Logs = append(logs.Logs, ExecutingLog{
		Key:   key,
		Value: value,
	})
}
func (logs *ExecutingLogs) Appends(log []ExecutingLog) {
	logs.Logs = append(logs.Logs, log...)
}

func (logs *ExecutingLogs) Marshal() (string, error) {
	data, err := json.Marshal(logs)
	return fmt.Sprintf("%v", string(data)), err
}

func (logs *ExecutingLogs) UnMarshal(data string) error {

	return json.Unmarshal([]byte(data), logs)
}
