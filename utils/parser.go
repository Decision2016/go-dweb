package utils

import (
	"context"
	"fmt"
	"github.com/gookit/config/v2"
	"github.com/sirupsen/logrus"
	"github.io/decision2016/go-dweb/interfaces"
	"path/filepath"
	"plugin"
	"strings"
)

func ParseOnChain(ident string) (*interfaces.IChain, error) {
	var identity Ident
	err := identity.FromString(ident)
	if err != nil {
		return nil, err
	}

	var symbol plugin.Symbol = nil

	switch identity.SubType {
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
	address := identity.Address
	err = chain.Setup(address)
	if err != nil {
		return nil, fmt.Errorf("setup chain interface failed")
	}
	return &chain, nil
}

func URLPathToChainIdent(url string) (string, error) {
	identArray := strings.Split(url, "/")
	logrus.Debugf("identity array length: %d", len(identArray))
	if len(identArray) < 3 {
		return "", fmt.Errorf("url path not valid")
	}

	ident := fmt.Sprintf("/chain/%s/%s", identArray[1], identArray[2])
	return ident, nil
}

func ExtractFilePath(url string) (string, error) {
	parts := strings.Split(url, "/")

	if len(parts) < 4 {
		return "", fmt.Errorf("url path too short")
	}

	remove := filepath.Join("/", parts[1], parts[2])
	result := strings.Replace(url, remove, "", 1)
	return result, nil
}

func ParseFileStorage(ctx context.Context, ident string) (string, *interfaces.IFileStorage, error) {
	var identity Ident
	err := identity.FromString(ident)
	if err != nil {
		return "", nil, err
	}

	var symbol plugin.Symbol

	symbol, err = LoadSymbol(identity.SubType)
	if err != nil {
		return "", nil, err
	}

	fs := symbol.(interfaces.IFileStorage)
	err = fs.Initial(ctx)
	if err != nil {
		return "", nil, err
	}

	return identity.Address, &fs, nil
}
