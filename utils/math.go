/**
  @author: decision
  @date: 2024/6/24
  @note:
**/

package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

func MerkleRoot(data []string) string {
	count := len(data)
	log := Log2(count)
	var root string

	if (2 << (log - 1)) == count {
		log -= 1
	}

	size := (1 << log)
	contents := make([][]byte, size, size)
	for idx, s := range data {
		contents[idx] = []byte(s)
	}

	for {
		if len(contents) <= 1 {
			root = hex.EncodeToString(contents[0])
			break
		}

		left := contents[0]
		right := contents[1]
		contents = contents[2:]

		sha := sha256.New()
		sha.Write(left)
		sha.Write(right)
		contents = append(contents, sha.Sum(nil))
	}

	return root
}

func Log2(num int) int {
	count := 0
	for {
		if num <= 0 {
			return count
		}
		num >>= 1
		count++
	}
}
