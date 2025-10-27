package contract

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type SafeCallData struct {
	Flag  *big.Int       `json:"flag"`  // uint256
	To    common.Address `json:"to"`    // address
	Value *big.Int       `json:"value"` // uint256
	Data  []byte         `json:"data"`  // bytes
	Hint  []byte         `json:"hint"`  // bytes
	Extra []byte         `json:"extra"` // bytes
}
