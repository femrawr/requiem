package higher

import "shared"

var cryptoKey []byte

func InitKey(key1 string, key2 string) []byte {
	key1Bytes := []byte(key1)
	key2Bytes := []byte(key2)

	maxLen := len(key1Bytes)
	if len(key2Bytes) > maxLen {
		maxLen = len(key2Bytes)
	}

	key1Buffer := make([]byte, maxLen)
	key2Buffer := make([]byte, maxLen)

	copy(key1Buffer, key1Bytes)
	copy(key2Buffer, key2Bytes)

	summed := make([]int, maxLen)
	for i := 0; i < maxLen; i++ {
		summed[i] = int(key1Buffer[i]) + int(key2Buffer[i])
	}

	xored := make([]int, maxLen)
	for i := 0; i < maxLen; i++ {
		xored[i] = int(key1Buffer[i]) ^ int(key2Buffer[i])
	}

	xoredAll := make([]int, maxLen)
	for i := 0; i < maxLen; i++ {
		val := summed[i]
		for j := 0; j < maxLen; j++ {
			val ^= xored[j]
		}

		xoredAll[i] = val
	}

	for _, v := range xoredAll {
		for v > 255 {
			cryptoKey = append(cryptoKey, byte(v&0xFF))
			v >>= 8
		}

		cryptoKey = append(cryptoKey, byte(v))
	}

	for _, v := range key2Bytes {
		cryptoKey = append(cryptoKey, v)
	}

	return cryptoKey
}

func EncryptConfig(config string) string {
	if config == "" {
		return ""
	}

	if cryptoKey == nil {
		return ""
	}

	enc, err := shared.EncryptData(config, cryptoKey, false)
	if err != nil {
		return ""
	}

	return enc
}

func DecryptConfig(config string) string {
	if config == "" {
		return ""
	}

	if cryptoKey == nil {
		return ""
	}

	dec, err := shared.DecryptData(config, cryptoKey, false)
	if err != nil {
		return ""
	}

	return dec
}
