package cachepolicy

type BackupPolicy[K comparable, V any] struct {
	datastore interface {
		Cache(interface{}) interface{}
	}
	interval uint32
}

func NewBackupPolicy[K comparable, V any](datastore interface{ Cache(interface{}) interface{} }, interval uint32) *BackupPolicy[K, V] {
	return &BackupPolicy[K, V]{
		datastore: datastore,
		interval:  interval,
	}
}

func (this *BackupPolicy[K, V]) FullBackup() {
	// keys, values := this.datastore.Cache(nil).(*expmap.ConcurrentMap[string, any]).KVs()
	// codec.Strings(keys).Encode()

	// encoder := this.datastore.Encoder(nil)
	// byteset := make([][]byte, len(keys))
	// for i := 0; i < len(keys); i++ {
	// 	byteset[i] = encoder(keys[i], values[i])
	// }
	// codec.Strings(keys).Encode()
}
