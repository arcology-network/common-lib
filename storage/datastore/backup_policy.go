package datastore

import "github.com/arcology-network/common-lib/codec"

type BackupPolicy struct {
	datastore *DataStore
	interval  uint32
}

func NewBackupPolicy(datastore *DataStore, interval uint32) *BackupPolicy {
	return &BackupPolicy{
		datastore: datastore,
		interval:  interval,
	}
}

func (this *BackupPolicy) FullBackup() {
	keys, values := this.datastore.Cache().KVs()
	codec.Strings(keys).Encode()

	encoder := this.datastore.Encoder()
	byteset := make([][]byte, len(keys))
	for i := 0; i < len(keys); i++ {
		byteset[i] = encoder(keys[i], values[i])
	}
	codec.Strings(keys).Encode()
}
