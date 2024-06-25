package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

func pkcs7UnPadding(text []byte) []byte {
	unPadding := int(text[len(text)-1])
	if len(text) < unPadding {
		return text
	}
	return text[:(len(text) - unPadding)]
}

func pkcs7Padding(text []byte, blockSize int) []byte {
	padding := blockSize - len(text)%blockSize
	var paddingText []byte
	if padding == 0 {
		paddingText = bytes.Repeat([]byte{byte(blockSize)}, blockSize)
	} else {
		paddingText = bytes.Repeat([]byte{byte(padding)}, padding)
	}
	return append(text, paddingText...)
}

func AesCBCDecryptBase64(data string, key []byte) ([]byte, error) {
	d, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}
	return AesCBCDecrypt(d, key)
}

func AesCBCEncryptBase64(data, key []byte) (string, error) {
	result, err := AesCBCEncrypt(data, key)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(result), nil
}

func AesCBCEncrypt(data, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	// 填充
	paddingText := pkcs7Padding(data, block.BlockSize())

	blockMode := cipher.NewCBCEncrypter(block, key[:block.BlockSize()])
	// 加密
	result := make([]byte, len(paddingText))
	blockMode.CryptBlocks(result, paddingText)
	// 返回密文
	return result, nil
}

func AesCBCDecrypt(data, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockMode := cipher.NewCBCDecrypter(block, key[:block.BlockSize()])

	result := make([]byte, len(data))

	blockMode.CryptBlocks(result, data)
	// 去除填充
	result = pkcs7UnPadding(result)
	return result, nil
}
