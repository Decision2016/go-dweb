/**
  @author: decision
  @date: 2024/6/5
  @note:
**/

package main

import (
	"fmt"
	"path/filepath"
)

//var example = "/chain/evm/0x7313AA3d928479Deb7debaC9E9d38286496D542e"

func main() {
	//fmt.Println(len(strings.Split(example, "/")))
	fmt.Println(filepath.Join("./plugins/", "infura.so"))
	fmt.Println(filepath.Join("./plugins", "infura.so"))
}
