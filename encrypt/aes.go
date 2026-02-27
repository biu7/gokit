package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
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
	paddingText := make([]byte, padding)
	for i := range paddingText {
		paddingText[i] = byte(padding)
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

// AesCBCEncrypt encrypts data using AES-CBC with a random IV.
// The returned ciphertext has the IV prepended (first blockSize bytes).
func AesCBCEncrypt(data, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	paddingText := pkcs7Padding(data, block.BlockSize())

	ciphertext := make([]byte, block.BlockSize()+len(paddingText))
	iv := ciphertext[:block.BlockSize()]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, fmt.Errorf("failed to generate IV: %w", err)
	}

	blockMode := cipher.NewCBCEncrypter(block, iv)
	blockMode.CryptBlocks(ciphertext[block.BlockSize():], paddingText)
	return ciphertext, nil
}

// AesCBCDecrypt decrypts data encrypted by AesCBCEncrypt.
// It expects the IV to be prepended to the ciphertext (first blockSize bytes).
func AesCBCDecrypt(data, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(data) < block.BlockSize() {
		return nil, fmt.Errorf("ciphertext too short")
	}

	iv := data[:block.BlockSize()]
	data = data[block.BlockSize():]

	if len(data)%block.BlockSize() != 0 {
		return nil, fmt.Errorf("ciphertext is not a multiple of the block size")
	}

	blockMode := cipher.NewCBCDecrypter(block, iv)
	result := make([]byte, len(data))
	blockMode.CryptBlocks(result, data)
	result = pkcs7UnPadding(result)
	return result, nil
}
