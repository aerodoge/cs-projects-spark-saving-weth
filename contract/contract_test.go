package contract

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
)

func TestBuildApprove(t *testing.T) {
	spender := common.HexToAddress("0xfe6eb3b609a7c8352a241f7f3a21cea4e9209b8f")
	amount := big.NewInt(10000000000000000)
	data, err := BuildApproveData(spender, amount)
	if err != nil {
		panic(err)
	}
	target := "0x095ea7b3000000000000000000000000fe6eb3b609a7c8352a241f7f3a21cea4e9209b8f000000000000000000000000000000000000000000000000002386f26fc10000"
	assert.Equal(t, hexutil.Encode(data), target)
}

func TestBuildSparkDeposit(t *testing.T) {
	assets := big.NewInt(10000000000000000)
	receiver := common.HexToAddress("0x6B69A26E05Ba29Dcf0C3E354BB0823Cc56Aeabfc")
	referral := uint16(128)
	data, err := BuildSparkDepositData(assets, receiver, referral)
	if err != nil {
		panic(err)
	}
	target := "0x9b8d6d38000000000000000000000000000000000000000000000000002386f26fc100000000000000000000000000006b69a26e05ba29dcf0c3e354bb0823cc56aeabfc0000000000000000000000000000000000000000000000000000000000000080"
	assert.Equal(t, hexutil.Encode(data), target)
}

func TestBuildSparkRedeem(t *testing.T) {
	shares := big.NewInt(9999376887387850)
	receiver := common.HexToAddress("0x6B69A26E05Ba29Dcf0C3E354BB0823Cc56Aeabfc")
	owner := common.HexToAddress("0x6B69A26E05Ba29Dcf0C3E354BB0823Cc56Aeabfc")
	data, err := BuildSparkRedeemData(shares, receiver, owner)
	if err != nil {
		panic(err)
	}
	target := "0xba087652000000000000000000000000000000000000000000000000002386615b5916ca0000000000000000000000006b69a26e05ba29dcf0c3e354bb0823cc56aeabfc0000000000000000000000006b69a26e05ba29dcf0c3e354bb0823cc56aeabfc"
	assert.Equal(t, hexutil.Encode(data), target)
}
