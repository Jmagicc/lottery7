package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

// AESEncrypt 加密文本
func AESEncrypt(plaintext string, key string) (string, error) {
	// 将key转换为32位
	keyBytes := []byte(key)
	keyBytes = PKCS7Padding(keyBytes, 32)

	// 创建加密块
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	// 对明文进行填充
	content := []byte(plaintext)
	content = PKCS7Padding(content, aes.BlockSize)

	// 创建初始化向量
	iv := make([]byte, aes.BlockSize)

	// 创建CBC加密模式
	mode := cipher.NewCBCEncrypter(block, iv)

	// 加密
	encrypted := make([]byte, len(content))
	mode.CryptBlocks(encrypted, content)

	// 使用base64编码
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

// AESDecrypt 解密文本
func AESDecrypt(ciphertext string, key string) (string, error) {
	// 将key转换为32位
	keyBytes := []byte(key)
	keyBytes = PKCS7Padding(keyBytes, 32)

	// 创建解密块
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	// 解码base64
	encrypted, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	// 创建初始化向量
	iv := make([]byte, aes.BlockSize)

	// 创建CBC解密模式
	mode := cipher.NewCBCDecrypter(block, iv)

	// 解密
	decrypted := make([]byte, len(encrypted))
	mode.CryptBlocks(decrypted, encrypted)

	// 去除填充
	decrypted = PKCS7UnPadding(decrypted)

	return string(decrypted), nil
}

// PKCS7Padding 填充
func PKCS7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// PKCS7UnPadding 去除填充
func PKCS7UnPadding(data []byte) []byte {
	length := len(data)
	unPadding := int(data[length-1])
	return data[:(length - unPadding)]
}
