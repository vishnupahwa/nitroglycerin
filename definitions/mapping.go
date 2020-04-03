package definitions

import "eznft/scenario"

// NFT is the mapping of scenario names to scenario definitions
var NFT = map[string]scenario.Scenario{
	"quick": quick(),
	"slow":  slow(),
}
