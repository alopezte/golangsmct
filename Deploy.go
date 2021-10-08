package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/alopezte/golangsmct/api" // this would be your generated smart contract bindings

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const PRIVATE_KEY string = "dc5ead86b0fe79a8d1d04f0aa8482b9f8bb556336e80c5a57bcb2ba4e6e2b018"

func main() {
	client, err := ethclient.Dial("http://127.0.0.1:7545")
	if err != nil {
		panic(err)
	}

	privateKey, err := crypto.HexToECDSA(PRIVATE_KEY)
	if err != nil {
		panic(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		panic("invalid key")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		panic(err)
	}

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		panic(err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		panic(err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)      // in wei
	auth.GasLimit = uint64(3000000) // in units
	auth.GasPrice = big.NewInt(1000000)

	address, tx, instance, err := api.DeployApi(auth, client)
	if err != nil {
		panic(err)
	}

	fmt.Println(address.Hex())

	_, _ = instance, tx
}
