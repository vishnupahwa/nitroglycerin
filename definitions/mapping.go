package definitions

import "eznft/scenario"

var NFT = map[string]scenario.Scenario{
	"quick": quick(),
	"slow":  slow(),
}
