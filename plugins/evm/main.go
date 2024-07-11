/**
  @author: decision
  @date: 2024/7/10
  @note:
**/

package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gookit/config/v2"
	"github.com/sirupsen/logrus"
	"math/big"
)

var Instance = EVMChain{}

type EVMChain struct {
	client   *ethclient.Client
	instance *Dweb

	privateKey string
	id         *big.Int
}

func (e *EVMChain) Identity() (string, error) {
	result, err := e.instance.Identity(nil)
	if err != nil {
		return "", err
	}

	return result, nil
}

func (e *EVMChain) Bootstrap() (string, error) {
	result, err := e.instance.Boostrap(nil)
	if err != nil {
		return "", err
	}

	return result, nil
}

func (e *EVMChain) Initial(ident string, url string) error {
	if !e.check() {
		return fmt.Errorf("required configuration not exist")
	}
	if e.privateKey == "" {
		return fmt.Errorf("config item deploy.chain.private is empty")
	}

	privateKey, err := crypto.HexToECDSA(e.privateKey)
	if err != nil {
		return err
	}

	auth, err := e.addressAuth(privateKey)
	if err != nil {
		return err
	}

	tx, err := e.instance.Initial(auth, ident)
	if err != nil {
		return err
	}

	logrus.Infof("contract initial with tx: %s", tx.Hash().Hex())
	return nil
}

func (e *EVMChain) SetIdentity(ident string) error {
	if !e.check() {
		return fmt.Errorf("required configuration not exist")
	}
	if e.privateKey == "" {
		return fmt.Errorf("config item deploy.chain.private is empty")
	}

	privateKey, err := crypto.HexToECDSA(e.privateKey)
	if err != nil {
		return err
	}

	auth, err := e.addressAuth(privateKey)
	if err != nil {
		return err
	}

	tx, err := e.instance.Update(auth, ident)
	if err != nil {
		return err
	}

	logrus.Infof("contract update with tx: %s", tx.Hash().Hex())
	return nil
}

func (e *EVMChain) Join(url string) error {
	if !e.check() {
		return fmt.Errorf("required configuration not exist")
	}
	if e.privateKey == "" {
		return fmt.Errorf("config item deploy.chain.private is empty")
	}

	privateKey, err := crypto.HexToECDSA(e.privateKey)
	if err != nil {
		return err
	}

	auth, err := e.addressAuth(privateKey)
	if err != nil {
		return err
	}

	tx, err := e.instance.Join(auth, url)
	if err != nil {
		return err
	}

	logrus.Infof("contract bootstrap set with tx: %s", tx.Hash().Hex())
	return nil
}

func (e *EVMChain) Setup(address string) error {
	api := config.String("chain.url")
	e.privateKey = config.String("chain.private", "")
	id := config.Int64("chain.id", 0)
	e.id = big.NewInt(id)

	client, err := ethclient.Dial(api)
	if err != nil {
		return err
	}

	e.client = client
	hexAddress := common.HexToAddress(address)
	instance, err := NewDweb(hexAddress, client)
	if err != nil {
		return err
	}

	e.instance = instance
	return nil
}

func (e *EVMChain) addressAuth(key *ecdsa.PrivateKey) (*bind.TransactOpts,
	error) {
	publicKey := key.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("error casting public key to ECDSA")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := e.client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return nil, err
	}

	gasPrice, err := e.client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(key, e.id)
	if err != nil {
		return nil, err
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(300000)
	auth.GasPrice = gasPrice

	return auth, nil
}

func (e *EVMChain) check() bool {
	var cfgs = []string{
		"chain.private",
		"chain.id",
		"chain.url",
	}

	for _, item := range cfgs {
		if !config.Exists(item) {
			logrus.Errorf("config item %s not exist", item)
			return false
		}
	}

	return true
}
