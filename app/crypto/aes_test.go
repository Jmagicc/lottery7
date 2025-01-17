package crypto

import (
	"fmt"
	"testing"
)

func TestAESEncryptDecrypt(t *testing.T) {
	tests := []struct {
		name      string
		plaintext string
		key       string
	}{
		{
			name:      "基本测试",
			plaintext: "Hello, World!",
			key:       "mysecretkey12345",
		},
		{
			name:      "中文测试",
			plaintext: "你好，世界！",
			key:       "这是一个密钥12345",
		},
		{
			name:      "空字符串测试",
			plaintext: "",
			key:       "testkey123456789",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 加密
			encrypted, err := AESEncrypt(tt.plaintext, tt.key)
			if err != nil {
				t.Errorf("AESEncrypt error: %v", err)
				return
			}
			fmt.Println(encrypted)

			// 解密
			decrypted, err := AESDecrypt(encrypted, tt.key)
			if err != nil {
				t.Errorf("AESDecrypt error: %v", err)
				return
			}

			// 验证结果
			if decrypted != tt.plaintext {
				t.Errorf("TestAESEncryptDecrypt failed: got %v, want %v", decrypted, tt.plaintext)
			}
		})
	}
}
