package request

// Vault represents the information necessary to deal with Vault.
type Vault struct {
	Address string `json:"address"`
	Token   string `json:"token"`
}

// DefaultVault provides a default Vault configuration by best effort.
func DefaultVault() Vault {
	return Vault{
		Address: "",
		Token:   "",
	}
}
