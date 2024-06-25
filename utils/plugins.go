/**
  @author: decision
  @date: 2024/6/25
  @note:
**/

package utils

import (
	"fmt"
	"github.com/gookit/config/v2"
	"path/filepath"
	"plugin"
)

func LoadSymbol(name string) (plugin.Symbol, error) {
	workDir := config.String("plugins.base", "")
	if workDir == "" {
		return nil, fmt.Errorf("plugin config not exists")
	}

	pluginPath := filepath.Join(workDir, fmt.Sprintf(name+".so"))
	p, err := plugin.Open(pluginPath)
	if err != nil {
		return nil, fmt.Errorf("load plugin file failed")
	}

	symbol, err := p.Lookup("Instance")
	if err != nil {
		return nil, err
	}

	return symbol, nil
}
