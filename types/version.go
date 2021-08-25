package types

// Version has the format |Actuator Prefix (1 byte)|Execution Unit Prefix (1 byte)|Workload (6 bytes)|.
// At most 256 executors, 256 execution units, 2^48-1 versions are supported.
type Version uint64

const (
	// InvalidVersion means no version info available.
	InvalidVersion = 0
	// MinVersion is the initial version for new element.
	MinVersion = 1
	// MaxVersion is the largest version number supported.
	// TODO: what should we do if version number reached MaxVersion?
	MaxVersion = 0xffffffffffff
	// VersionMask & version => 'Workload' part.
	VersionMask = 0xffffffffffff
	// PrefixMask & version => 'Prefix' part.
	PrefixMask = 0xffff000000000000
)
