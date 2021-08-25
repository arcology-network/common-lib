package types

type DeferCall struct {
	DeferID         string
	ContractAddress Address
	Signature       string
}

type TxAccessRecords struct {
	Hash     string
	ID       uint32
	Accesses []interface{}
}

type EuResult struct {
	H           string
	ID          uint32
	Transitions []interface{}
	DC          *DeferCall
	Status      uint64
	GasUsed     uint64
}
