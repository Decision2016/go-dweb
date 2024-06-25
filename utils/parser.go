package utils

import (
	"context"
	"fmt"
	"github.com/gookit/config/v2"
	"github.io/decision2016/go-dweb/interfaces"
	"plugin"
	"strings"
)

func ParseOnChain(ident string) (*interfaces.IChain, error) {
	identArray := strings.Split(ident, "/")
	if len(identArray) != 4 || identArray[1] != "chain" {
		return nil, fmt.Errorf("on-chain identity not correct")
	}

	var symbol plugin.Symbol = nil
	var err error

	switch identArray[2] {
	case "evm":
		symbol, err = LoadSymbol("evm")
	case "norn":
		symbol, err = LoadSymbol("norn")
	case "custom":
		pluginName := config.String("plugins.chain", "")
		if pluginName == "" {
			return nil, fmt.Errorf("plugin binary file not found")
		}
		symbol, err = LoadSymbol(pluginName)
	}

	if err != nil {
		return nil, err
	}
	chain := symbol.(interfaces.IChain)
	address := identArray[3]
	err = chain.Setup(address)
	if err != nil {
		return nil, fmt.Errorf("setup chain interface failed")
	}
	return &chain, nil
}

func ParseFileStorage(ctx context.Context, ident string) (string,
	*interfaces.IFileStorage,
	error) {
	//	/storage/ipfs/QSsdafaw
	identArray := strings.Split(ident, "/")

	if len(identArray) != 4 || identArray[1] != "storage" {
		return "", nil, fmt.Errorf("on-chain identity not correct")
	}

	var symbol plugin.Symbol
	var err error

	symbol, err = LoadSymbol(identArray[2])
	if err != nil {
		return "", nil, err
	}

	fs := symbol.(interfaces.IFileStorage)
	err = fs.Initial(ctx)
	if err != nil {
		return "", nil, err
	}

	return identArray[3], &fs, nil
}
