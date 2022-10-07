module github.com/HPISTechnologies/common-lib

go 1.17

require (
	github.com/HPISTechnologies/3rd-party v1.3.1-0.20220302005842-3524e305a016
	github.com/HPISTechnologies/evm v1.10.4-0.20220123034347-eb8d747ab2b2
	github.com/elliotchance/orderedmap v1.4.0
	github.com/google/uuid v1.1.5
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519
)

require (
	github.com/btcsuite/btcd v0.21.0-beta // indirect
	golang.org/x/sys v0.0.0-20210615035016-665e8c7367d1 // indirect
)

// replace github.com/HPISTechnologies/evm => ../evm/

// replace github.com/HPISTechnologies/3rd-party => ../3rd-party/
