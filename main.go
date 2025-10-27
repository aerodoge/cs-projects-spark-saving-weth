package main

import (
	"cs-projects-spark-saving-weth/contract"
	"math/big"
)

func main() {
	rpcURL := "http://127.0.0.1:8545"
	privateKey := ""
	safe := "0x969b37A287bBFb4080E6cb293fB6E21995fd1f83"
	coboSafe := "0xad11d63611B133c6659F1C9127c9DD6d4e1aFf0D"
	delegate := contract.NewDelegate(rpcURL, privateKey, safe, coboSafe)
	//amount, _ := new(big.Int).SetString("100000000000000000000", 10)
	//delegate.ApproveWETHToSpark(amount)
	//delegate.Deposit(amount)
	shares, _ := new(big.Int).SetString("999909250609038983997", 10)
	delegate.Redeem(shares)
}
