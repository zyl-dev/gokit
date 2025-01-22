package conv

import (
	"encoding/hex"
	"fmt"
	"strconv"
)

// Hex2Bin 十六进制转二进制
func Hex2Bin(hexString string) (string, error) {
	ui, err := strconv.ParseUint(hexString, 16, 64)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%016b", ui), nil
}

func Hex2BinCoinBaseTags(hexString string) (string, error) {
	hexByte, err := hex.DecodeString(hexString)
	return string(hexByte), err
}
