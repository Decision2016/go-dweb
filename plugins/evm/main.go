/**
  @author: decision
  @date: 2024/7/10
  @note:
**/

package main

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gookit/config/v2"
)

type EVMChain struct {
	client   *ethclient.Client
	instance *Dweb
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
	//TODO implement me
	panic("implement me")
}

func (e *EVMChain) SetIdentity(ident string) error {
	//TODO implement me
	panic("implement me")
}

func (e *EVMChain) Join(url string) error {
	//TODO implement me
	panic("implement me")
}

func (e *EVMChain) Setup(address string) error {
	api := config.String("chain.api")
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
