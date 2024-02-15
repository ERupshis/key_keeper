package ska

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
)

type AESKeyLength int

const (
	Key16 = AESKeyLength(16)
	Key24 = AESKeyLength(24)
	Key32 = AESKeyLength(32)
)

type SKA struct {
	keyAES []byte
}

func NewSKA(userKey string, AESKeyLength AESKeyLength) *SKA {
	return &SKA{
		keyAES: generateKey(userKey, AESKeyLength),
	}
}

func (s *SKA) SetAESKey(userKey string, AESKeyLength AESKeyLength) {
	s.keyAES = generateKey(userKey, AESKeyLength)
}

func (s *SKA) Encrypt(rawText []byte) ([]byte, error) {
	errMsg := "encrypt bytes: %w"
	block, err := aes.NewCipher(s.keyAES)
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}

	paddedData := padData(rawText, aes.BlockSize)
	paddedCiphertext := make([]byte, aes.BlockSize+len(paddedData))
	iv := paddedCiphertext[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(paddedCiphertext[aes.BlockSize:], paddedData)

	encodedData := make([]byte, base64.StdEncoding.EncodedLen(len(paddedCiphertext)))
	base64.StdEncoding.Encode(encodedData, paddedCiphertext)
	return encodedData, nil
}

func (s *SKA) Decrypt(ciphertext []byte) ([]byte, error) {
	errMsg := "decode: %w"

	decodedCiphertext := make([]byte, base64.StdEncoding.DecodedLen(len(ciphertext)))
	n, err := base64.StdEncoding.Decode(decodedCiphertext, ciphertext)
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}

	iv := decodedCiphertext[:aes.BlockSize]
	decodedCiphertext = decodedCiphertext[aes.BlockSize:n]

	block, err := aes.NewCipher(s.keyAES)
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(decodedCiphertext, decodedCiphertext)

	return unPadData(decodedCiphertext), nil
}

func generateKey(input string, AESKeyLength AESKeyLength) []byte {
	hasher := sha256.New()
	hasher.Write([]byte(input))
	hashedKey := hasher.Sum(nil)

	if len(hashedKey) == int(AESKeyLength) {
		return hashedKey
	}

	for len(hashedKey) < int(AESKeyLength) {
		hasher.Reset()
		hasher.Write(hashedKey)
		hashedKey = append(hashedKey, hasher.Sum(nil)...)
	}

	return hashedKey[:AESKeyLength]
}

func padData(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	if padding == blockSize {
		return data
	}

	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

func unPadData(data []byte) []byte {
	if len(data) == 0 {
		return data
	}

	padding := int(data[len(data)-1])
	if len(data) < padding {
		return data
	}
	return data[:len(data)-padding]
}
