package config

// Patch is the cluster specific configuration.
type Patch struct {
	Workers []Worker `json:"workers,omitempty"`
}

// DefaultPatch provides a default patch by best effort.
func DefaultPatch() Patch {
	return Patch{
		Workers: []Worker{},
	}
}
