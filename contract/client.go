package contract

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	sparkAddress = common.HexToAddress("0xfE6eb3b609a7C8352A241f7F3A21CEA4e9209B8f")
	wethAddress  = common.HexToAddress("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2")
)

type Delegate struct {
	client     *ethclient.Client
	address    common.Address
	privateKey *ecdsa.PrivateKey
	safe       common.Address
	coboSafe   common.Address
}

func NewDelegate(rpcUrl, delegatePrivateKey, safeAddr, coboSafeAddr string) *Delegate {
	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		log.Fatal(err)
	}
	privateKey, err := crypto.HexToECDSA(delegatePrivateKey)
	if err != nil {
		log.Fatal(err)
	}
	address := crypto.PubkeyToAddress(privateKey.PublicKey)
	return &Delegate{
		client:     client,
		address:    address,
		privateKey: privateKey,
		safe:       common.HexToAddress(safeAddr),
		coboSafe:   common.HexToAddress(coboSafeAddr),
	}
}

func (d *Delegate) ApproveWETHToSpark(amount *big.Int) {
	//assets, _ := new(big.Int).SetString("100000000000000000000", 10) // 100ETH
	// approve
	approveData, err := BuildApproveData(sparkAddress, amount)
	if err != nil {
		log.Printf("Execute WETH approve Error: %v\n", err)
		return
	}
	txHash, err := d.executeSafeTransaction(wethAddress, big.NewInt(0), approveData)
	if err != nil {
		log.Printf("Execute Spark Deposit Error: %v\n", err)
		return
	}
	err = d.waitForConfirmation(txHash)
	if err != nil {
		log.Printf("Execute Spark Deposit Error: %v\n", err)
	}
}

func (d *Delegate) Deposit(amount *big.Int) {
	// deposit
	depositData, err := BuildSparkDepositData(amount, d.safe, uint16(128))
	if err != nil {
		log.Printf("Build Spark Deposit Error: %v\n", err)
		return
	}
	txHash, err := d.executeSafeTransaction(sparkAddress, big.NewInt(0), depositData)
	if err != nil {
		log.Printf("Execute Spark Deposit Error: %v\n", err)
		return
	}
	err = d.waitForConfirmation(txHash)
	if err != nil {
		log.Printf("Execute Spark Deposit Error: %v\n", err)
	}
}

func (d *Delegate) Redeem(amount *big.Int) {
	redeemData, err := BuildSparkRedeemData(amount, d.safe, d.safe)
	if err != nil {
		log.Printf("Build Spark Redeem Error: %v\n", err)
		return
	}

	txHash, err := d.executeSafeTransaction(sparkAddress, big.NewInt(0), redeemData)
	if err != nil {
		log.Printf("Execute Spark Redeem Error: %v\n", err)
		return
	}

	err = d.waitForConfirmation(txHash)
	if err != nil {
		log.Printf("Execute Spark Redeem Error: %v\n", err)
		return
	}
}

// 通过Safe执行交易x
func (d *Delegate) executeSafeTransaction(to common.Address, value *big.Int, data []byte) (*common.Hash, error) {
	// 构建Safe交易
	GasLimit := 6000000
	// 获取nonce
	nonce, err := d.client.PendingNonceAt(context.Background(), d.address)
	if err != nil {
		return nil, fmt.Errorf("获取nonce失败: %v", err)
	}

	// 获取gas价格
	gasPrice, err := d.client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, fmt.Errorf("获取gas价格失败: %v", err)
	}

	// 构建Safe执行交易的calldata
	safeExecData := d.buildSafeExecTransactionData(to, value, data)
	fmt.Printf("Gas price: %v, gas limit: %v\n", gasPrice, GasLimit)
	// 创建交易
	tx := types.NewTransaction(
		nonce,
		d.coboSafe, // 发送到Cobo Safe合约
		big.NewInt(0),
		uint64(GasLimit),
		gasPrice,
		safeExecData,
	)

	chainID, err := d.client.NetworkID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("获取chain ID失败: %v", err)
	}
	fmt.Println("ChainID:", chainID)
	// 签名交易
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), d.privateKey)
	if err != nil {
		return nil, fmt.Errorf("签名交易失败: %v", err)
	}

	// 发送交易
	err = d.client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return nil, fmt.Errorf("发送交易失败: %v", err)
	}

	hash := signedTx.Hash()
	return &hash, nil
}

func (d *Delegate) buildSafeExecTransactionData(to common.Address, value *big.Int, data []byte) []byte {
	contractABI, err := abi.JSON(strings.NewReader(safeABI))
	if err != nil {
		return nil
	}
	callData := SafeCallData{
		Flag:  big.NewInt(0),
		To:    to,
		Value: value,
		Data:  data,
		Hint:  []byte{},
		Extra: []byte{},
	}

	// 使用ABI编码
	encodedData, err := contractABI.Pack("execTransaction", callData)
	if err != nil {
		log.Fatal("build safe exec transaction err: ", err)
		return nil
	}
	return encodedData
}

func (d *Delegate) waitForConfirmation(txHash *common.Hash) error {
	fmt.Printf("等待交易确认: %s\n", txHash.Hex())

	for i := 0; i < 60; i++ { // 最多等待5分钟
		receipt, err := d.client.TransactionReceipt(context.Background(), *txHash)
		if err == nil {
			if receipt.Status == types.ReceiptStatusSuccessful {
				fmt.Printf("交易确认成功! Gas使用: %d\n", receipt.GasUsed)
				return nil
			} else {
				return fmt.Errorf("交易执行失败\n")
			}
		}
		time.Sleep(5 * time.Second)
	}

	return fmt.Errorf("交易确认超时")
}
