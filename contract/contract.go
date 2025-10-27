package contract

import (
	_ "embed"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

//go:embed abis/weth.abi.json
var wethABI string

//go:embed abis/spark.abi.json
var sparkABI string

//go:embed abis/safe.abi.json
var safeABI string

func BuildApproveData(spender common.Address, amount *big.Int) ([]byte, error) {
	data, err := abi.JSON(strings.NewReader(wethABI))
	if err != nil {
		return nil, err
	}

	encodedData, err := data.Pack("approve", spender, amount)
	if err != nil {
		return nil, err
	}

	return encodedData, nil
}

func BuildSparkDepositData(assets *big.Int, receiver common.Address, referral uint16) ([]byte, error) {
	data, err := abi.JSON(strings.NewReader(sparkABI))
	if err != nil {
		return nil, err
	}

	encodedData, err := data.Pack("deposit", assets, receiver, referral)
	if err != nil {
		return nil, err
	}

	return encodedData, nil
}

func BuildSparkRedeemData(shares *big.Int, receiver, owner common.Address) ([]byte, error) {
	data, err := abi.JSON(strings.NewReader(sparkABI))
	if err != nil {
		return nil, err
	}

	encodedData, err := data.Pack("redeem", shares, receiver, owner)
	if err != nil {
		return nil, err
	}

	return encodedData, nil
}

func calculateMethodID(signature string) []byte {
	hash := crypto.Keccak256Hash([]byte(signature))
	return hash[:4]
}
